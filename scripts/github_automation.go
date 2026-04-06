//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Config represents the input JSON structure
type Config struct {
	Owner  string  `json:"owner"`
	Repo   string  `json:"repo"`
	Labels []Label `json:"labels"`
	Issues []Issue `json:"issues"`
}

type Label struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type Issue struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Labels    []string `json:"labels"`
	Assignees []string `json:"assignees"`
}

func main() {
	configPath := flag.String("config", "scripts/audit-issues-config.json", "path to config JSON")
	token := flag.String("token", "", "GitHub Personal Access Token")
	dryRun := flag.Bool("dry-run", false, "Preview changes without execution")
	flag.Parse()

	// 1. Resolve token
	ghToken := *token
	if ghToken == "" {
		ghToken = os.Getenv("GITHUB_TOKEN")
	}
	if ghToken == "" {
		// Try reading from .env file
		envData, err := os.ReadFile(".env")
		if err == nil {
			lines := strings.Split(string(envData), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "GITHUB_TOKEN=") {
					ghToken = strings.TrimSpace(strings.TrimPrefix(line, "GITHUB_TOKEN="))
					break
				}
			}
		}
	}

	if ghToken == "" {
		fmt.Println("❌ Error: GITHUB_TOKEN not provided and not found in environment or .env file")
		os.Exit(1)
	}

	// 2. Load config
	data, err := os.ReadFile(*configPath)
	if err != nil {
		fmt.Printf("❌ Error reading config: %v\n", err)
		os.Exit(1)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("❌ Error parsing config: %v\n", err)
		os.Exit(1)
	}

	client := &http.Client{}

	fmt.Printf("🚀 Starting GitHub Automation for %s/%s\n", cfg.Owner, cfg.Repo)
	if *dryRun {
		fmt.Println("⚠️  DRY-RUN MODE: No changes will be made")
	}
	fmt.Println(strings.Repeat("-", 50))

	// 3. Create Labels
	fmt.Printf("Step 1: Creating/Verifying Labels (%d)...\n", len(cfg.Labels))
	for _, l := range cfg.Labels {
		if *dryRun {
			fmt.Printf("  [DRY-RUN] Would create label: %s (#%s)\n", l.Name, l.Color)
			continue
		}

		success, msg := createLabel(client, cfg.Owner, cfg.Repo, ghToken, l)
		fmt.Printf("  %s\n", msg)
		_ = success
	}

	// 4. Create Issues
	fmt.Printf("\nStep 2: Creating Issues (%d)...\n", len(cfg.Issues))
	for i, issue := range cfg.Issues {
		if *dryRun {
			fmt.Printf("  [DRY-RUN] Would create issue #%d: %s\n", i+1, issue.Title)
			continue
		}

		success, msg := createIssue(client, cfg.Owner, cfg.Repo, ghToken, issue)
		if success {
			fmt.Printf("  ✅ Created Issue #%d: %s\n", i+1, issue.Title)
		} else {
			fmt.Printf("  ❌ Failed Issue #%d: %s (%s)\n", i+1, issue.Title, msg)
		}
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("✨ Task complete!")
}

func createLabel(client *http.Client, owner, repo, token string, l Label) (bool, string) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/labels", owner, repo)

	body, _ := json.Marshal(l)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("Error request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		return true, fmt.Sprintf("Created: %s", l.Name)
	}
	if resp.StatusCode == 422 {
		return true, fmt.Sprintf("Skipped (exists): %s", l.Name)
	}

	respBody, _ := io.ReadAll(resp.Body)
	return false, fmt.Sprintf("Failed: %s (Status: %d, Msg: %s)", l.Name, resp.StatusCode, string(respBody))
}

func createIssue(client *http.Client, owner, repo, token string, iss Issue) (bool, string) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)

	body, _ := json.Marshal(iss)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("Error request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		return true, "Success"
	}

	respBody, _ := io.ReadAll(resp.Body)
	return false, string(respBody)
}
