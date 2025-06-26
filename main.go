package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/zlahrouni/loganizer/internal/analyzer"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	logFiles := []string{
		"test_logs/access.log",
		"test_logs/corrupted.log",
		"test_logs/empty.log",
		"test_logs/errors.log",
	}

	results, errs := analyzer.AnalyzeLogs(logFiles)

	fmt.Println("\nLog Analysis Results:")
	for _, res := range results {
		fmt.Printf("✓ %s: %d entries analyzed in %v\n",
			res.FilePath, res.Entries, res.Duration)
	}

	fmt.Println("\nErrors:")
	for _, e := range errs {
		var notFoundErr *analyzer.FileNotFoundOrUnreadableError
		var parsingErr *analyzer.ParsingError

		switch {
		case errors.As(e, &notFoundErr):
			fmt.Printf("⚠️ %v\n", notFoundErr)
		case errors.As(e, &parsingErr):
			fmt.Printf("⚠️ %v\n", parsingErr)
		default:
			fmt.Printf("⚠️ Unexpected error: %v\n", e)
		}
	}
}
