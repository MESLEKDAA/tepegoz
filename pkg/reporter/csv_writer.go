package reporter

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"time"

	"github.com/MESLEKDAA/tepegoz/pkg/model"
)

func WriteToCsv(filename string, event model.LogEvent) error {

	reportDir := "reports"
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return err
	}

	fullPath := filepath.Join(reportDir, filename)

	fileExists := false

	if _, err := os.Stat(fullPath); err == nil {
		fileExists = true
	}

	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	defer writer.Flush()

	if !fileExists {

		header := []string{"Timestamp", "Level", "RuleID", "RuleName", "SourceFile", "LogContent"}

		if err := writer.Write(header); err != nil {
			return err
		}

	}

	record := []string{
		event.TimeStamp.Format(time.RFC3339),
		event.Level,
		event.RuleID,
		event.RuleName,
		event.SourceFile,
		event.OriginalLine,
	}

	return writer.Write(record)
}
