package client

import (
	"encoding/csv"
	"math"
	"os"
	"path/filepath"

	"copilotlens/internal/github"
)

var GitHubClient *github.Client

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

func LoadUsers(dataDir string) []string {
	return LoadUsersFromClient(GitHubClient)
}

func LoadUsersFromClient(c *github.Client) []string {
	if c != nil {
		members, err := c.GetOrgMembers()
		if err == nil && len(members) > 0 {
			return members
		}
	}
	return nil
}

func Round2(f float64) float64 {
	return math.Round(f*100) / 100
}
