//go:build ignore

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type moduleSpec struct {
	id                string
	title             string
	testPkgs          []string
	requiredFiles     []string
	requiredRepoFiles []string
	readmeDir         string
	nextStep          string
}

var modules = []moduleSpec{
	{
		id:    "OPSL.1",
		title: "Foundation and Configuration",
		testPkgs: []string{
			"./internal/config/...",
		},
		requiredFiles: []string{
			"cmd/server/main.go",
			"internal/config/config.go",
			"internal/config/environment.go",
			".env.example",
			"Dockerfile",
			"docker-compose.yml",
		},
		readmeDir: "modules/01-foundation",
		nextStep:  "Continue into OPSL.2 after the server boots with validated configuration.",
	},
	{
		id:    "OPSL.2",
		title: "Database and Models",
		testPkgs: []string{
			"./internal/db/...",
		},
		requiredFiles: []string{
			"internal/db/migrations.go",
			"internal/db/repository.go",
			"internal/models/tenant.go",
			"internal/models/user.go",
			"internal/models/order.go",
			"internal/models/payment.go",
		},
		readmeDir: "modules/02-database",
		nextStep:  "Continue into OPSL.3 after the schema and repository layer are stable.",
	},
	{
		id:    "OPSL.3",
		title: "Authentication and Tenant Isolation",
		testPkgs: []string{
			"./internal/auth/...",
		},
		requiredFiles: []string{
			"internal/auth/token.go",
			"internal/auth/password.go",
			"internal/auth/service.go",
			"internal/auth/middleware.go",
			"internal/auth/context.go",
		},
		readmeDir: "modules/03-auth",
		nextStep:  "Continue into OPSL.4 after tenant identity is trusted middleware state.",
	},
	{
		id:    "OPSL.4",
		title: "HTTP API Layer",
		testPkgs: []string{
			"./internal/handlers/...",
			"./internal/middleware/...",
		},
		requiredFiles: []string{
			"internal/handlers/handlers.go",
			"internal/handlers/api.go",
			"internal/middleware/middleware.go",
		},
		readmeDir: "modules/04-http-api",
		nextStep:  "Continue into OPSL.5 to replace CRUD-shaped writes with order workflow logic.",
	},
	{
		id:    "OPSL.5",
		title: "Order Processing",
		testPkgs: []string{
			"./internal/services/...",
		},
		requiredFiles: []string{
			"internal/services/order.go",
			"internal/services/inventory.go",
			"internal/services/validation.go",
		},
		readmeDir: "modules/05-order-processing",
		nextStep:  "Build the order service and state machine before starting OPSL.6.",
	},
	{
		id:    "OPSL.6",
		title: "Payment Pipeline",
		testPkgs: []string{
			"./internal/payment/...",
			"./internal/services/...",
			"./internal/handlers/...",
		},
		requiredFiles: []string{
			"internal/payment/gateway.go",
			"internal/payment/worker.go",
			"internal/services/payment.go",
		},
		readmeDir: "modules/06-payment-pipeline",
		nextStep:  "Continue into OPSL.7 after payment retries and reconciliation are provable.",
	},
	{
		id:    "OPSL.7",
		title: "Event Bus and Worker Pools",
		testPkgs: []string{
			"./internal/events/...",
			"./internal/workers/...",
		},
		requiredFiles: []string{
			"internal/events/bus.go",
			"internal/events/types.go",
			"internal/workers/pool.go",
			"internal/workers/order_processor.go",
			"internal/workers/payment_processor.go",
			"internal/workers/notification_worker.go",
		},
		readmeDir: "modules/07-event-workers",
		nextStep:  "Start bounded async processing after payment workflow logic exists.",
	},
	{
		id:    "OPSL.8",
		title: "Caching Layer",
		testPkgs: []string{
			"./internal/cache/...",
		},
		requiredFiles: []string{
			"internal/cache/cache.go",
			"internal/cache/store.go",
			"internal/middleware/cache.go",
		},
		readmeDir: "modules/08-caching",
		nextStep:  "Introduce cache behavior only after OPSL.7 has clear read paths to optimize.",
	},
	{
		id:    "OPSL.9",
		title: "Observability",
		testPkgs: []string{
			"./internal/logging/...",
			"./internal/metrics/...",
		},
		requiredFiles: []string{
			"internal/logging/logger.go",
			"internal/logging/middleware.go",
			"internal/metrics/metrics.go",
			"internal/tracing/tracing.go",
		},
		readmeDir: "modules/09-observability",
		nextStep:  "Add logs, metrics, and traces before the final deployment module.",
	},
	{
		id:    "OPSL.10",
		title: "Graceful Shutdown and Deployment",
		testPkgs: []string{
			"./cmd/server",
		},
		requiredFiles: []string{
			"cmd/server/shutdown.go",
		},
		requiredRepoFiles: []string{
			".github/workflows/ci.yml",
		},
		readmeDir: "modules/10-shutdown-deploy",
		nextStep:  "Finish the full flagship path by proving safe shutdown and deployment behavior.",
	},
}

type moduleResult struct {
	complete       bool
	missingFiles   []string
	failedPackages []string
}

func main() {
	root, repoRoot, err := findOpslaneRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println("Opslane module progress")
	fmt.Println(strings.Repeat("=", 72))

	results := make([]moduleResult, len(modules))
	currentIndex := -1

	for i, module := range modules {
		results[i] = checkModule(root, repoRoot, module)
		if currentIndex == -1 && !results[i].complete {
			currentIndex = i
		}
	}

	for i, module := range modules {
		status := "locked"
		switch {
		case results[i].complete:
			status = "complete"
		case i == currentIndex:
			status = "current"
		}
		fmt.Printf("%-8s %-36s %s\n", module.id, module.title, status)
	}

	if currentIndex == -1 {
		fmt.Println()
		fmt.Println("All Opslane modules are complete.")
		return
	}

	module := modules[currentIndex]
	result := results[currentIndex]

	fmt.Println()
	fmt.Printf("Current module: %s %s\n", module.id, module.title)
	fmt.Printf("Module guide   : 11-flagship/01-opslane/%s/README.md\n", module.readmeDir)

	if len(result.missingFiles) > 0 {
		fmt.Println("Missing files:")
		for _, file := range result.missingFiles {
			fmt.Printf("  - %s\n", file)
		}
	}

	if len(result.failedPackages) > 0 {
		fmt.Println("Proof checks to fix:")
		for _, pkg := range result.failedPackages {
			fmt.Printf("  - %s\n", pkg)
		}
	}

	fmt.Println("Next step:")
	fmt.Printf("  %s\n", module.nextStep)
}

func checkModule(root, repoRoot string, module moduleSpec) moduleResult {
	result := moduleResult{}

	for _, relativePath := range module.requiredFiles {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(relativePath))); os.IsNotExist(err) {
			result.missingFiles = append(result.missingFiles, relativePath)
		}
	}

	for _, relativePath := range module.requiredRepoFiles {
		if _, err := os.Stat(filepath.Join(repoRoot, filepath.FromSlash(relativePath))); os.IsNotExist(err) {
			result.missingFiles = append(result.missingFiles, "repo:"+relativePath)
		}
	}

	for _, pkg := range module.testPkgs {
		dir := strings.TrimPrefix(pkg, "./")
		dir = strings.TrimSuffix(dir, "/...")
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(dir))); os.IsNotExist(err) {
			result.failedPackages = append(result.failedPackages, pkg+" (package missing)")
			continue
		}

		cmd := exec.Command("go", "test", "-count=1", "-timeout=60s", pkg)
		cmd.Dir = root
		if err := os.MkdirAll(filepath.Join(root, ".cache", "go-build"), 0o755); err == nil {
			cmd.Env = append(os.Environ(), "GOCACHE="+filepath.Join(root, ".cache", "go-build"))
		}
		if output, err := cmd.CombinedOutput(); err != nil {
			summary := strings.TrimSpace(string(output))
			lines := strings.Split(summary, "\n")
			if len(lines) > 0 {
				summary = lines[len(lines)-1]
			}
			result.failedPackages = append(result.failedPackages, fmt.Sprintf("%s (%s)", pkg, summary))
		}
	}

	result.complete = len(result.missingFiles) == 0 && len(result.failedPackages) == 0
	return result
}

func findOpslaneRoot() (string, string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", "", fmt.Errorf("cannot determine progress.go location")
	}

	root := filepath.Dir(filepath.Dir(currentFile))
	markers := []string{
		"cmd/server/main.go",
		"internal/config/config.go",
		"docker-compose.yml",
	}

	for _, marker := range markers {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(marker))); err != nil {
			return "", "", fmt.Errorf("cannot find Opslane root from %s", currentFile)
		}
	}

	repoRoot := filepath.Dir(filepath.Dir(root))
	repoMarkers := []string{
		"go.mod",
		"curriculum.v2.json",
	}

	for _, marker := range repoMarkers {
		if _, err := os.Stat(filepath.Join(repoRoot, filepath.FromSlash(marker))); err != nil {
			return "", "", fmt.Errorf("cannot find repository root from %s", currentFile)
		}
	}

	return root, repoRoot, nil
}
