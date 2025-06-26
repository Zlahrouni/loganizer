package analyzer

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type FileNotFoundOrUnreadableError struct {
	FilePath string
	Err      error
}

func (e *FileNotFoundOrUnreadableError) Error() string {
	return fmt.Sprintf("file %s not found or unreadable: %v", e.FilePath, e.Err)
}

func (e *FileNotFoundOrUnreadableError) Unwrap() error {
	return e.Err
}

type ParsingError struct {
	FilePath string
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("parsing error for file %s", e.FilePath)
}

type AnalysisResult struct {
	FilePath string
	Entries  int
	Duration time.Duration
}

func AnalyzeLogs(filePaths []string) ([]AnalysisResult, []error) {
	var wg sync.WaitGroup
	results := make(chan AnalysisResult, len(filePaths))
	errs := make(chan error, len(filePaths))

	for _, filePath := range filePaths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			res, err := analyzeFile(path)
			if err != nil {
				errs <- err
				return
			}
			results <- res
		}(filePath)
	}

	wg.Wait()
	close(results)
	close(errs)

	var resultSlice []AnalysisResult
	for res := range results {
		resultSlice = append(resultSlice, res)
	}

	var errSlice []error
	for e := range errs {
		errSlice = append(errSlice, e)
	}

	return resultSlice, errSlice
}

func analyzeFile(filePath string) (AnalysisResult, error) {
	delay := time.Duration(50+rand.Intn(150)) * time.Millisecond
	time.Sleep(delay)

	file, err := os.Open(filePath)
	if err != nil {
		return AnalysisResult{}, &FileNotFoundOrUnreadableError{
			FilePath: filePath,
			Err:      err,
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if rand.Float32() < 0.1 {
		return AnalysisResult{}, &ParsingError{FilePath: filePath}
	}

	return AnalysisResult{
		FilePath: filePath,
		Entries:  lineCount,
		Duration: delay,
	}, nil
}
