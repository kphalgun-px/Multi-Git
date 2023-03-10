package helpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
)

func DebugPrint(output ...any) {
	fmt.Println()
	fmt.Println("------------------")
	fmt.Println(output...)
	fmt.Println("------------------")
}

func ConfigureGit() (err error) {
	// err = exec.Command("git", "config", "--global", "user.name", "kphalgun").Run()
	// if err != nil {
	// 	return
	// }
	// err = exec.Command("git", "config", "--global", "user.email", "kphalgun@purestorage.com").Run()

	return
}

func CreateDir(baseDir string, name string, initGit bool) (err error) {
	dirName := path.Join(baseDir, name)

	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return
	}

	if !initGit {
		return
	}

	currDir, err := os.Getwd()
	if err != nil {
		return
	}
	defer os.Chdir(currDir)

	os.Chdir(dirName)
	err = exec.Command("git", "init").Run()

	return
}

func AddFiles(baseDir string, dirName string, commit bool, filenames ...string) (err error) {
	dir := path.Join(baseDir, dirName)

	for _, f := range filenames {
		data := []byte("data for" + f)
		err = os.WriteFile(path.Join(dir, f), data, 0777)
		if err != nil {
			return
		}
	}

	if !commit {
		return
	}

	currDir, err := os.Getwd()
	if err != nil {
		return
	}
	defer os.Chdir(currDir)

	os.Chdir(dir)
	err = exec.Command("git", "add", "-A").Run()
	if err != nil {
		return
	}

	err = exec.Command("git", "commit", "-m", "added some files...").Run()

	return
}

func RunMultiGit(command string, ignoreErrors bool, mgRoot string, mgRepos string) (output string, err error) {
	var out []byte

	if runtime.GOOS == "windows" {
		out, err = exec.Command("where.exe", "multi-git").CombinedOutput()
	} else {
		out, err = exec.Command("which", "multi-git").CombinedOutput()
	}

	if err != nil {
		return
	}

	if len(out) == 0 {
		err = errors.New("multi-git is not in the PATH")
		return
	}

	components := []string{"--command", command}
	if ignoreErrors {
		components = append(components, "--ignore-errors")
	}
	cmd := exec.Command("multi-git", components...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "MG_ROOT="+mgRoot, "MG_REPOS="+mgRepos)
	out, err = cmd.CombinedOutput()
	output = string(out)
	return
}
