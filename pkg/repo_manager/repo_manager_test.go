package repo_manager_test

import (
	"fmt"
	. "multi-git/pkg/helpers"
	. "multi-git/pkg/repo_manager"
	"os"
	"path"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const baseDir = "/tmp/test-multi-git"

var repoList = []string{}

var _ = Describe("Repo manager tests", func() {
	var err error

	removeAll := func() {
		err = os.RemoveAll(baseDir)
		Expect(err).Should(BeNil())
	}

	BeforeEach(func() {
		err = ConfigureGit()
		Expect(err).Should(BeNil())
		removeAll()

		err = CreateDir(baseDir, "dir-1", true)
		Expect(err).Should(BeNil())
		repoList = []string{"dir-1"}
	})
	AfterEach(removeAll)

	It("Should fail with invalid base dir", func() {
		_, err := NewRepoManager("/no-such-dir", repoList, true)
		Expect(err).ShouldNot(BeNil())
	})

	It("Should fail with empty repo list", func() {
		_, err := NewRepoManager(baseDir, []string{}, true)
		Expect(err).ShouldNot(BeNil())
	})

	It("Should get repo list successfully", func() {
		rm, err := NewRepoManager(baseDir, repoList, true)
		Expect(err).Should(BeNil())

		repos := rm.GetRepos()
		Expect(repos).Should(HaveLen(1))
		Expect(repos[0] == path.Join(baseDir, repoList[0])).Should(BeTrue())
	})

	It("Should get repo list successfully with non-git directories", func() {
		repoList = append(repoList, "dir-2")
		CreateDir(baseDir, repoList[1], true)
		CreateDir(baseDir, "not-a-git-repo", false)

		rm, err := NewRepoManager(baseDir, repoList, true)
		Expect(err).Should(BeNil())

		repos := rm.GetRepos()
		Expect(repos).Should(HaveLen(2))
		Expect(repos[0] == path.Join(baseDir, repoList[0])).Should(BeTrue())
		Expect(repos[1] == path.Join(baseDir, repoList[1])).Should(BeTrue())
	})

	It("Should get repo list successfully with non-git directories", func() {
		repoList = append(repoList, "dir-2")
		CreateDir(baseDir, repoList[1], true)
		CreateDir(baseDir, "not-a-git-repo", false)

		rm, err := NewRepoManager(baseDir, repoList, true)
		Expect(err).Should(BeNil())

		repos := rm.GetRepos()
		Expect(repos).Should(HaveLen(2))
		Expect(repos[0] == path.Join(baseDir, repoList[0])).Should(BeTrue())
		Expect(repos[1] == path.Join(baseDir, repoList[1])).Should(BeTrue())
	})

	It("Should create branches successfully", func() {
		repoList = append(repoList, "dir-2")
		CreateDir(baseDir, repoList[1], true)

		rm, err := NewRepoManager(baseDir, repoList, true)
		Expect(err).Should(BeNil())

		output, err := rm.Exec("checkout -b test-branch")
		Expect(err).Should(BeNil())

		for _, out := range output {
			Expect(out).Should(Equal("Switched to a new branch 'test-branch'\n"))
		}
	})

	It("Should commit files successfully", func() {
		rm, err := NewRepoManager(baseDir, repoList, true)
		Expect(err).Should(BeNil())

		output, err := rm.Exec("checkout -b test-branch")
		Expect(err).Should(BeNil())

		for _, out := range output {
			Expect(out).Should(Equal("Switched to a new branch 'test-branch'\n"))
		}

		err = AddFiles(baseDir, repoList[0], true, "file_1.txt", "file_2.txt")
		Expect(err).Should(BeNil())

		wd, _ := os.Getwd()
		defer os.Chdir(wd)

		dir := path.Join(baseDir, repoList[0])
		err = os.Chdir(dir)
		Expect(err).Should(BeNil())

		output, err = rm.Exec("log --oneline")
		fmt.Println(output)
		Expect(err).Should(BeNil())

		ok := strings.HasSuffix(output[dir], "added some files...\n")
		Expect(ok).Should(BeTrue())
	})
})
