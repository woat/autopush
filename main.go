package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
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

func main() {
	fmt.Println("Starting watcher...")

	// this shit does not watch recursively gg
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

	if len(os.Args) < 2 {
		err = watcher.Add("./")
	} else {
		err = watcher.Add(os.Args[1])
	}

	if err != nil {
		panic(err)
	}

	<-done
}
