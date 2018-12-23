package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
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

func commit(msg string) {
	git := exec.Command("git", "commit", "-m", msg)
	gitOut, err := git.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(gitOut))
}

func push() {
	git := exec.Command("git", "push", "-f")
	gitOut, err := git.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(gitOut))
}

func main() {
	fmt.Println("Starting watcher...")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Chmod != fsnotify.Chmod {
					msg := fmt.Sprintf("%s: %s\n", event.Op, event.Name)
					add()
					commit(msg)
					push()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				panic(err)
			}
		}
	}()

	err = watcher.Add("./")
	if err != nil {
		panic(err)
	}
	<-done
}
