package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
)

type Args struct {
	file string
	q    int
}

func parseArgs() (Args, error) {
	file := flag.String("file", "", "path to the input file")
	q := flag.Int("q", 3, "q-grams")
	flag.Parse()

	if *file == "" {
		return Args{}, fmt.Errorf("file must be specified")
	}

	return Args{
		file: *file,
		q:    *q,
	}, nil
}

func main() {
	args, err := parseArgs()
	if err != nil {
		slog.Error("error parsing arguments", "error", err)
		os.Exit(1)
	}

	index := NewQGramIndex(args.q)

	slog.Info("building index",
		"file", args.file,
		"q", args.q,
	)
	err = index.BuildFormFile(args.file)
	if err != nil {
		slog.Error("error building index", "error", err)
		os.Exit(1)
	}

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Query: ")
		scanner.Scan()
		query := scanner.Text()

		query = Normalize(query)
		delta := len(query) / (args.q + 1)

		postings, err := index.FindMatches(query, delta)
		if err != nil {
			slog.Error("error finding matches", "error", err)
		}
		slog.Info("found matches", "count", len(postings))
		// Print only the first 5 results.
		for _, posting := range postings[:min(5, len(postings))] {
			infos := index.Infos[posting.ID]
			fmt.Printf(
				"\n\033[1m%s\033[0m (score=%d, ped=%d, qid=%s):\n%s\n",
				infos.Name,
				infos.Score,
				posting.PED,
				infos.Infos[0],
				infos.Infos[1],
			)
		}
	}
}
