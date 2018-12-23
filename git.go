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
