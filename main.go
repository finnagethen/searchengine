package main

import (
	"fmt"
	"os"
)

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

	for {
		var query string
		fmt.Print("Query: ")
		_, err := fmt.Scanln(&query)
		if err != nil {
			return
		}

		query = normalize(query)
		delta := len(query) / index.Q

		postings, err := index.FindMatches(query, delta)
		if err != nil {
			panic(err)
		}

		// Print only the first 5 results.
		for _, posting := range postings[:min(5, len(postings))] {
			infos := index.GetInfo(posting.ID)
			fmt.Printf(
				"\n\033[1m%s\033[0m (score=%d, ped=%d, qid=%s, via '%s'):\n%s\n",
				infos.Name,
				infos.Score,
				posting.PED,
				infos.Infos[0],
				infos.Infos[1],
			)
		}
	}
}
