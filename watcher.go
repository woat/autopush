package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
)

func watch(path string) {
	fmt.Println("Starting watcher...")

	// TODO: make this recursive
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

	if path == "" {
		err = watcher.Add("./")
	} else {
		// TODO: validate path
		err = watcher.Add(path)
	}

	if err != nil {
		panic(err)
	}

	<-done
}
