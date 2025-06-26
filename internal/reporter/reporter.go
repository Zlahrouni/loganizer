package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zlahrouni/loganizer/internal/analyzer"
	"github.com/zlahrouni/loganizer/internal/config"
)

func ExportResultsToJsonFile(filePath string, results []analyzer.FileNotFoundOrUnreadableError) error {
	timestampedPath := addTimestamp(filePath)
	createDir(timestampedPath)

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(timestampedPath, data, 0644)
}

func ExportReport(filePath string, configs []config.LogConfig, results []analyzer.AnalysisResult, errs []error, statusFilter string) error {
	timestampedPath := addTimestamp(filePath)
	createDir(timestampedPath)

	report := buildReport(configs, results, errs)

	if statusFilter != "" {
		report = filterReport(report, statusFilter)
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(timestampedPath, data, 0644)
}

func AddLog(configFile, id, path, logType string) error {
	configs, err := config.LoadConfig(configFile)
	if err != nil {
		return err
	}

	for _, cfg := range configs {
		if cfg.ID == id {
			return fmt.Errorf("ID %s existe déjà", id)
		}
	}

	newConfig := config.LogConfig{
		ID:   id,
		Path: path,
		Type: logType,
	}
	configs = append(configs, newConfig)

	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}

func addTimestamp(filePath string) string {
	ext := filepath.Ext(filePath)
	name := strings.TrimSuffix(filePath, ext)
	timestamp := time.Now().Format("060102")
	return fmt.Sprintf("%s_%s%s", name, timestamp, ext)
}

func createDir(filePath string) {
	os.MkdirAll(filepath.Dir(filePath), 0755)
}

func buildReport(configs []config.LogConfig, results []analyzer.AnalysisResult, errs []error) []map[string]interface{} {
	var report []map[string]interface{}

	resultMap := make(map[string]analyzer.AnalysisResult)
	for _, r := range results {
		resultMap[r.FilePath] = r
	}

	errorMap := make(map[string]string)
	for _, err := range errs {
		if fnfErr, ok := err.(*analyzer.FileNotFoundOrUnreadableError); ok {
			errorMap[fnfErr.FilePath] = err.Error()
		}
		if parseErr, ok := err.(*analyzer.ParsingError); ok {
			errorMap[parseErr.FilePath] = err.Error()
		}
	}

	for _, cfg := range configs {
		entry := map[string]interface{}{
			"log_id":    cfg.ID,
			"file_path": cfg.Path,
		}

		if result, ok := resultMap[cfg.Path]; ok {
			entry["status"] = "OK"
			entry["entries"] = result.Entries
		} else if errMsg, ok := errorMap[cfg.Path]; ok {
			entry["status"] = "FAILED"
			entry["error"] = errMsg
		} else {
			entry["status"] = "NOT_PROCESSED"
		}

		report = append(report, entry)
	}

	return report
}

func filterReport(report []map[string]interface{}, statusFilter string) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, entry := range report {
		if entry["status"] == statusFilter {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
