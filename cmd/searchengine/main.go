package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"

	"github.com/finnagethen/searchengine/internal/embeddingindex"
	"github.com/finnagethen/searchengine/internal/qgramindex"
	"github.com/finnagethen/searchengine/internal/utils"
)

// BUG: "hunger games" query does not find "The Hunger Games"

type Args struct {
	filePath       string
	embeddingsPath string
	q              int
	topK           int
}

func parseArgs() (Args, error) {
	filePath := flag.String("file", "", "path to the input file")
	embeddingsPath := flag.String("embeddings", "./data/embeddings.bin", "path to the embeddings file")
	q := flag.Int("q", 3, "q-grams")
	topK := flag.Int("k", 5, "number of top similar documents to return")
	flag.Parse()

	if *filePath == "" {
		return Args{}, fmt.Errorf("file must be specified")
	}

	return Args{
		filePath:       *filePath,
		embeddingsPath: *embeddingsPath,
		q:              *q,
		topK:           *topK,
	}, nil
}

func main() {
	go func() {
		log.Print(http.ListenAndServe(":1234", nil))
	}()

	// Redirect logs to JSON file.
	//file, err := os.Create("logs.json")
	//if err != nil {
	//	panic(err)
	//}
	//defer file.Close()
	//
	//logger := slog.New(
	//	slog.NewJSONHandler(file, &slog.HandlerOptions{}),
	//)
	//
	//slog.SetDefault(logger)

	// Parse command line arguments.
	args, err := parseArgs()
	if err != nil {
		slog.Error("error parsing arguments", "error", err)
		os.Exit(1)
	}

	// Initialize the QGramIndex.
	qIndex := qgramindex.NewQGramIndex(args.q)

	slog.Info("building qIndex",
		"file", args.filePath,
		"q", args.q,
	)
	err = qIndex.BuildFormFile(args.filePath)
	if err != nil {
		slog.Error("error building qIndex", "error", err)
		os.Exit(1)
	}

	// Initialize the EmbeddingIndex.
	eIndex := embeddingindex.NewEmbeddingIndex()

	slog.Info("loading embeddings",
		"file", args.embeddingsPath,
	)
	err = eIndex.LoadEmbeddings(args.embeddingsPath)
	if err != nil {
		slog.Error("error loading embeddings", "error", err)
		os.Exit(1)
	}

	// Query loop.
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Query: ")
		scanner.Scan()
		query := scanner.Text()

		query = utils.Normalize(query)
		delta := len(query) / (args.q + 1)

		postings, err := qIndex.FindMatches(query, delta)
		if err != nil {
			slog.Error("error finding matches", "error", err)
		}
		slog.Info("found matches", "count", len(postings))

		// Print only the first 5 results.
		var names []string
		var plots []string

		numResults := min(5, len(postings))
		if numResults == 0 {
			fmt.Println("No matching movies found.")
			continue
		}

		for i, posting := range postings[:numResults] {
			info := qIndex.Infos[posting.ID]
			name, score, infos := info.Name, info.Score, info.Infos
			plot := infos[1]
			year := infos[0]

			fmt.Printf("  %d. %s (%s | %d votes)\n", i+1, name, year, score)
			names = append(names, name)
			plots = append(plots, plot)

			//fmt.Printf(
			//	"\n\033[1m%s\033[0m (score=%d, ped=%d, qid=%s):\n%s\n",
			//	infos.Name,
			//	infos.Score,
			//	posting.PED,
			//	infos.Infos[0],
			//	infos.Infos[1],
			//)
		}

		fmt.Print("\nSelect a movie: ")
		if !scanner.Scan() {
			break
		}
		selectionStr := scanner.Text()

		selection, err := strconv.Atoi(selectionStr)
		if err != nil || selection < 1 || selection > numResults {
			fmt.Println("Invalid selection.")
			continue
		}

		index := selection - 1
		name := names[index]

		fmt.Printf("\nTop %d most similar movies to '%s' (and the movie itself):\n", args.topK, name)

		// Search for top K + 1 (the first result is always the movie itself).
		documentID := int(postings[index].ID)
		indices, err := eIndex.TopKNeighbors(documentID, args.topK+1)
		if err != nil {
			slog.Error("error finding top k neighbors", "error", err)
			continue
		}

		for i, idx := range indices {
			info := qIndex.Infos[idx]
			nName, nInfos := info.Name, info.Infos
			nYear := nInfos[0]
			nPlot := nInfos[1]

			if len(nPlot) > 1000 {
				nPlot = nPlot[:1000] + "..."
			}

			header := fmt.Sprintf("%d. %s %s", i, nName, nYear)
			separator := strings.Repeat("-", len(header))
			fmt.Printf("  %s\n  %s\n  %s\n\n", header, separator, nPlot)
		}
	}
}
