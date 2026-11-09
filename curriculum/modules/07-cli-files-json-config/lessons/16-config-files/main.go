package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type AppConfig struct {
	ServerPort int      `json:"server_port"`
	LogLevel   string   `json:"log_level"`
	Database   DBConfig `json:"database"`
	Debug      bool     `json:"debug"`
}

func DefaultConfig() AppConfig {
	return AppConfig{
		ServerPort: 8080,
		LogLevel:   "info",
		Database: DBConfig{
			Host: "localhost",
			Port: 5432,
		},
		Debug: false,
	}
}

func LoadAppConfig(paths ...string) (AppConfig, error) {
	cfg := DefaultConfig()
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return cfg, err
		}
		if err := json.Unmarshal(data, &cfg); err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}

func main() {
	// Create temporary config files for demonstration
	os.WriteFile("_test_base.json", []byte(`{"server_port":9090,"database":{"host":"prod.example.com"}}`), 0644)
	os.WriteFile("_test_override.json", []byte(`{"debug":true,"database":{"port":5433}}`), 0644)
	defer os.Remove("_test_base.json")
	defer os.Remove("_test_override.json")

	cfg, err := LoadAppConfig("_test_base.json", "_test_override.json", "_test_nonexistent.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("ServerPort: %d\n", cfg.ServerPort)
	fmt.Printf("LogLevel: %s\n", cfg.LogLevel)
	fmt.Printf("Debug: %v\n", cfg.Debug)
	fmt.Printf("DB Host: %s\n", cfg.Database.Host)
	fmt.Printf("DB Port: %d\n", cfg.Database.Port)
}
