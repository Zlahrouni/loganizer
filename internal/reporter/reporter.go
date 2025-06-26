package reporter

import (
	"encoding/json"
	"github.com/zlahrouni/loganizer/internal/analyzer"
	"os"
)

func ExportResultsToJsonFile(filePath string, results []analyzer.FileNotFoundOrUnreadableError) error {
	data, err := json.MarshalIndent(results, "", "")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}
	return nil

}
