package main

import (
	"fmt"
	"os"

	"github.com/weiser/lox/mainhelpers"
)

//TODO start 9.4 while loops on 143
func main() {
	// [0] is the program name
	args := os.Args[1:]
	switch len(args) {
	case 1:
		mainhelpers.RunFile(args[0])
	case 0:
		fmt.Println("Starting mainHelper...")
		mainhelpers.RunPrompt()
	default:
		println("Usage: lox [path to script]")
		os.Exit(1)
	}

}
