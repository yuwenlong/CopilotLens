package client

import (
	"encoding/csv"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"copilotlens/domain/dto"
)

func LoadUsernameMap(dataDir string) map[string]string {
	m := make(map[string]string)
	f, err := os.Open(filepath.Join(dataDir, "username.csv"))
	if err != nil {
		return m
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil {
		return m
	}

	for i, row := range records {
		if i == 0 || len(row) < 2 {
			continue
		}
		m[row[0]] = row[1]
	}
	return m
}

func LoadCopilotCSV(dataDir, month string) []dto.CopilotRecord {
	filePath := filepath.Join(dataDir, "copilot", month+".csv")
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	rows, err := reader.ReadAll()
	if err != nil {
		return nil
	}

	var records []dto.CopilotRecord
	for i, row := range rows {
		if i == 0 || len(row) < 9 {
			continue
		}
		qty, _ := strconv.ParseFloat(row[5], 64)
		cost, _ := strconv.ParseFloat(row[8], 64)
		records = append(records, dto.CopilotRecord{
			Date:        row[0],
			Username:    row[1],
			Model:       row[4],
			AICQuantity: qty,
			AICCost:     cost,
		})
	}
	return records
}

func Round2(f float64) float64 {
	return math.Round(f*100) / 100
}
