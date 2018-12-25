package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"strings"
)

type watcher struct {
	*fsnotify.Watcher
	paths map[string]bool
}

func (w *watcher) handleEvents() {
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Chmod != fsnotify.Chmod {
				msg := fmt.Sprintf("%s: %s\n", event.Op, event.Name)
				outf(msg)
				// TODO: fix commit + prevent pollution
				add()
				err := commit(msg)
				if err != nil {
					break
				}
				push()
				reset()

				// If dir -> Remove from watcher
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					for p, _ := range w.paths {
						if strings.Contains(p, event.Name) {
							msg, err := w.removePath(p)
							if err != nil {
								panic(err)
							}
							outf(msg)
						}
					}
				}

				// If dir -> Add to watcher
				if event.Op&fsnotify.Create == fsnotify.Create {
					file, err := os.Stat(event.Name)
					if err != nil {
						panic(err)
					}
					if file.IsDir() {
						msg, err := w.addPath(event.Name)
						if err != nil {
							panic(err)
						}
						outf(msg)
					}
				}
			}
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			panic(err)
		}
	}
}

func (w *watcher) addPath(p string) (string, error) {
	w.Add(p)
	if _, ok := w.paths[p]; !ok {
		w.paths[p] = true
	} else {
		return "", fmt.Errorf("Attempted to add pre-existing path: %s\n", p)
	}
	msg := fmt.Sprintf("Added %s to watcher\n", p)
	return msg, nil
}

func (w *watcher) removePath(p string) (string, error) {
	if _, ok := w.paths[p]; ok {
		delete(w.paths, p)
	} else {
		return "", fmt.Errorf("Attempted to remove non-existant path: %s\n", p)
	}
	w.Remove(p)

	msg := fmt.Sprintf("Removed %s from watcher\n", p)
	return msg, nil
}

func newWatcher() *watcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	paths := make(map[string]bool)

	watcher := &watcher{w, paths}
	return watcher
}

func lsSubdirs(root string) []string {
	var dirs []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// TODO: see git.go
			// TODO: reduce number of calls to exec
			// TODO: may have to change approach to filter-last, see #w.handleEvents
			if !strings.Contains(path, ".git") && !shouldIgnore(path) {
				dirs = append(dirs, path)
			}
		}
		return err
	})

	if err != nil {
		panic(err)
	}

	return dirs
}

func watch(path string) {
	watcher := newWatcher()
	defer watcher.Close()

	done := make(chan bool)

	go watcher.handleEvents()
	outf("Starting watcher on path: %s\n", path)
	paths := lsSubdirs(path)

	for _, p := range paths {
		outf("Watching: %s\n", p)

		watcher.Add(p)
		watcher.paths[p] = true
	}

	<-done
}
