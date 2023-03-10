package repo_manager

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type RepoManager struct {
	repos        []string
	ignoreErrors bool
}

func NewRepoManager(baseDir string, repoNames []string, ignoreErrors bool) (repoManager *RepoManager, err error) {
	_, err = os.Stat(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("base dir: '%s' does not exist", baseDir)
		}
		return
	}

	if baseDir[len(baseDir)-1] != '/' {
		baseDir += "/"
	}

	if len(repoNames) == 0 {
		err = errors.New("repo list cannot be empty")
		return
	}

	repoManager = &RepoManager{
		ignoreErrors: ignoreErrors,
	}

	for _, r := range repoNames {
		if r == "" {
			err = errors.New("repo name cannot be empty")
			return
		}
		path := baseDir + r
		repoManager.repos = append(repoManager.repos, path)
	}

	return
}

func (m *RepoManager) GetRepos() []string {
	return m.repos
}

func (m *RepoManager) Exec(cmd string) (output map[string]string, err error) {
	output = map[string]string{}
	var components []string
	var multiWord []string

	for _, component := range strings.Split(cmd, " ") {
		if strings.HasPrefix(component, "\"") {
			multiWord = append(multiWord, component[1:])
			continue
		}

		if len(multiWord) > 0 {
			if !strings.HasSuffix(component, "\"") {
				multiWord = append(multiWord, component)
				continue
			}
			multiWord = append(multiWord, component[:len(component)-1])
			component = strings.Join(multiWord, " ")
			multiWord = []string{}
		}

		components = append(components, component)
	}

	wd, _ := os.Getwd()
	defer os.Chdir(wd)

	var out []byte

	for _, r := range m.repos {
		os.Chdir(r)
		out, err = exec.Command("git", components...).CombinedOutput()
		output[r] = string(out)
		if err != nil && !m.ignoreErrors {
			return
		}
	}

	return
}
