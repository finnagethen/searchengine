package main

import "os"

func parseArgs() []string {
	if len(os.Args) != 2 {
		panic("Usage: ./searchengine <path-to-file>")
	}

	return os.Args[1:]
}

func main() {
	args := parseArgs()

	index := NewQGramIndex(3)
	err := index.BuildFormFile(args[0])
	if err != nil {
		panic(err)
	}
}
