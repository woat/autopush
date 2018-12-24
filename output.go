package main

import "fmt"

const prompt = "\033[1;34m[autopush]\033[0m "

func printfout(f string, args ...interface{}) {
	f = prompt + f
	fmt.Printf(f, args...)
}
