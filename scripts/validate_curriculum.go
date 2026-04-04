package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Lesson struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type Section struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Lessons []Lesson `json:"lessons"`
}

type Curriculum struct {
	Sections []Section `json:"sections"`
}

func main() {
	data, err := os.ReadFile("curriculum.json")
	if err != nil {
		fmt.Printf("❌ Failed to read curriculum.json: %v\\n", err)
		os.Exit(1)
	}

	var cur Curriculum
	if err := json.Unmarshal(data, &cur); err != nil {
		fmt.Printf("❌ Failed to parse curriculum.json: %v\\n", err)
		os.Exit(1)
	}

	errors := 0
	lessonCount := 0
	for _, s := range cur.Sections {
		for _, l := range s.Lessons {
			lessonCount++
			if l.Path == "" {
				fmt.Printf("⚠️  Unmapped lesson: %s - %s\\n", l.ID, l.Name)
				errors++
				continue
			}

			fullPath := l.Path
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				fmt.Printf("❌ Path does not exist: %s (%s - %s)\\n", fullPath, l.ID, l.Name)
				errors++
			}
		}
	}

	if errors == 0 {
		fmt.Printf("✅ Success! All %d lessons correctly mapped to the filesystem.\\n", lessonCount)
	} else {
		fmt.Printf("❌ Found %d mapping errors.\\n", errors)
		os.Exit(1)
	}
}
