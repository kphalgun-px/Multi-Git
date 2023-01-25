package e2e_tests_test

import (
	"fmt"
	"os"
	"strings"

	. "github.com/kphalgun/multi-git/pkg/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const baseDir = "/tmp/test-multi-git"

var err error
var repoList string

var removeAll = func() {
	err = os.RemoveAll(baseDir)
	Expect(err).Should(BeNil())
}

var _ = Describe("multi-git e2e tests", func() {
	fmt.Println("e2e_tests starting")

	BeforeEach(func() {
		err = ConfigureGit()
		Expect(err).Should(BeNil())
		removeAll()
		err = CreateDir(baseDir, "", false)
		Expect(err).Should(BeNil())
	})

	Context("Tests for empty/undefined environment failure cases", func() {
		It("Should fail with invalid base dir", func() {
			output, err := RunMultiGit("status", false, "/no-such-dir", repoList)
			Expect(err).ShouldNot(BeNil())
			suffix := "base dir: '/no-such-dir/' does not exist\n"
			Expect(output).Should(HaveSuffix(suffix))
		})

		It("Should fail with empty repo list", func() {
			output, err := RunMultiGit("status", false, baseDir, repoList)
			Expect(err).ShouldNot(BeNil())
			Expect(output).Should(ContainSubstring("repo list cannot be empty"))
		})
	})

	Context("Tests for success cases", func() {
		It("Should do git init successfully", func() {
			err = CreateDir(baseDir, "dir-1", false)
			Expect(err).Should(BeNil())
			err = CreateDir(baseDir, "dir-2", false)
			Expect(err).Should(BeNil())
			repoList = "dir-1,dir-2"

			output, err := RunMultiGit("init", false, baseDir, repoList)
			Expect(err).Should(BeNil())
			fmt.Println(output)
			count := strings.Count(output, "Initialized empty Git repository")
			Expect(count).Should(Equal(2))
		})

		It("Should do git status successfully for git directories", func() {
			err = CreateDir(baseDir, "dir-1", true)
			Expect(err).Should(BeNil())
			err = CreateDir(baseDir, "dir-2", true)
			Expect(err).Should(BeNil())
			repoList = "dir-1,dir-2"

			output, err := RunMultiGit("status", false, baseDir, repoList)
			Expect(err).Should(BeNil())
			count := strings.Count(output, "nothing to commit")
			Expect(count).Should(Equal(2))
		})

		It("Should create branches successfully", func() {
			err = CreateDir(baseDir, "dir-1", true)
			Expect(err).Should(BeNil())
			err = CreateDir(baseDir, "dir-2", true)
			Expect(err).Should(BeNil())
			repoList = "dir-1,dir-2"

			output, err := RunMultiGit("checkout -b test-branch", false, baseDir, repoList)
			Expect(err).Should(BeNil())

			count := strings.Count(output, "Switched to a new branch 'test-branch'")
			Expect(count).Should(Equal(2))
		})
	})

	Context("Tests for non-git directories", func() {
		It("Should fail git status", func() {
			err = CreateDir(baseDir, "dir-1", false)
			Expect(err).Should(BeNil())
			err = CreateDir(baseDir, "dir-2", false)
			Expect(err).Should(BeNil())
			repoList = "dir-1,dir-2"

			output, err := RunMultiGit("status", false, baseDir, repoList)
			Expect(err).Should(BeNil())
			Expect(output).Should(ContainSubstring("fatal: not a git repository"))
		})
	})

	Context("Tests for ignoreErrors flag", func() {
		Context("First directory is invalid", func() {
			When("ignoreErrors is true", func() {
				It("git status should succeed for the second directory", func() {
					err = CreateDir(baseDir, "dir-1", false)
					Expect(err).Should(BeNil())
					err = CreateDir(baseDir, "dir-2", true)
					Expect(err).Should(BeNil())
					repoList = "dir-1,dir-2"

					output, err := RunMultiGit("status", true, baseDir, repoList)
					Expect(err).Should(BeNil())

					Expect(output).Should(ContainSubstring("[dir-1]: git status\nfatal: not a git repository"))
					Expect(output).Should(ContainSubstring("[dir-2]: git status\nOn branch master"))
				})
			})

			When("ignoreErrors is false", func() {
				It("Should fail on first directory and bail out", func() {
					err = CreateDir(baseDir, "dir-1", false)
					Expect(err).Should(BeNil())
					err = CreateDir(baseDir, "dir-2", true)
					Expect(err).Should(BeNil())
					repoList = "dir-1,dir-2"

					output, err := RunMultiGit("status", false, baseDir, repoList)
					Expect(err).Should(BeNil())

					Expect(output).Should(ContainSubstring("[dir-1]: git status\nfatal: not a git repository"))
					Expect(output).ShouldNot(ContainSubstring("[dir-2]"))
				})
			})
		})
	})
})

var _ = AfterSuite(removeAll)
