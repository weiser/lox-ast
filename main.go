package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// [0] is the program name
	args := os.Args[1:]
	switch len(args) {
	case 1:
		runFile(args[0])
	case 0:
		fmt.Println("Starting interpreter...")
		runPrompt()
	default:
		println("Usage: lox [path to script]")
		os.Exit(1)
	}

}

func runFile(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("couldn't read file: %v", err)
		os.Exit(1)
	}

	run(string(data))
}

func run(data string) {
	//TODO: maybe refactor runPrompt and runFile to pass a scanner into run?
	scanner := bufio.NewScanner(strings.NewReader(data))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")

	for scanner.Scan() {
		line := scanner.Text()
		if line == "q" || line == "" {
			fmt.Println("Exiting. goodbye")
			break
		}
		run(line)
		fmt.Print("> ")
	}
}
