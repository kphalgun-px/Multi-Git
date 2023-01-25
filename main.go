package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	command := flag.String("command", "", "The git command")
	ignoreErros := flag.Bool(
		"ignore-errors",
		false,
		"Keep running after error if true")
	flag.Parse()
	root := os.Getenv("MG_ROOT")
	if root[len(root)-1] != '/' {
		root += "/"
	}
	repo_names := strings.Split(os.Getenv("MG_REPOS"), ",")
	var repos []string
	for _, r := range repo_names {
		path := root + r
		_, err := os.Stat(path + "/.git")
		if err != nil {
			log.Fatal(err)
		}
		repos = append(repos, path)
	}
	var git_components []string
	git_components = append(git_components, strings.Split(*command, " ")...)
	command_string := "git " + *command
	for _, r := range repos {
		os.Chdir(r)
		fmt.Printf("[%s] %s\n", r, command_string)
		out, err := exec.Command("git", git_components...).CombinedOutput()
		fmt.Println(string(out))
		if err != nil && !*ignoreErros {
			os.Exit(1)
		}
	}
	fmt.Println("Done.")
}
