package main

import (
	"fmt"
	"os/exec"
)

func add() {
	git := exec.Command("git", "add", ".")
	gitOut, err := git.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(gitOut))
}

// this fucks up when we remove a dir
func commit(msg string) {
	git := exec.Command("git", "commit", "-m", msg)
	gitOut, err := git.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(gitOut))
}

func push() {
	git := exec.Command("git", "push", "-f", "origin", "head:autopush")
	gitOut, err := git.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(gitOut))
}

// TODO: handle exit codes properly (0, 1, 128) -> (true, false, fatal)
// TODO: take a list of paths and return all paths that are true
func shouldIgnore(path string) bool {
	git := exec.Command("git", "check-ignore", path)
	_, err := git.Output() // run() ?
	if err != nil {
		return false
	}
	return true
}
