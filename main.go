package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

type Args struct {
	file string
	q    int
}

func parseArgs() Args {
	file := flag.String("file", "", "path to the input file")
	q := flag.Int("q", 3, "q-grams")
	flag.Parse()

	return Args{
		file: *file,
		q:    *q,
	}
}

func main() {
	args := parseArgs()
	if args.file == "" {
		log.Fatal("file must be specified")
	}

	index := NewQGramIndex(args.q)

	log.Printf("Building index for %s with q=%d.", args.file, args.q)
	err := index.BuildFormFile(args.file)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done building index.")

	for {
		var query string
		fmt.Print("Query: ")
		_, err := fmt.Scanln(&query)
		if err != nil {
			log.Println(err)
			return
		}

		query = normalize(query)
		delta := len(query) / (args.q + 1)

		start := time.Now()
		postings, err := index.FindMatches(query, delta)
		if err != nil {
			log.Fatal(err)
		}
		elapsed := time.Since(start)
		log.Printf("Found %d matches in %s", len(postings), elapsed)

		// Print only the first 5 results.
		for _, posting := range postings[:min(5, len(postings))] {
			infos := index.Infos[posting.ID]
			fmt.Printf(
				"\n\033[1m%s\033[0m (score=%d, ped=%d)\n",
				infos.Name,
				infos.Score,
				posting.PED,
			)
		}
	}
}
