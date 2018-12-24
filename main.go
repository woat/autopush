package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		watch("./")
	} else {
		watch(os.Args[1])
	}
}
