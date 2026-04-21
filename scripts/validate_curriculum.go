//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/rasel9t6/the-go-engineer/scripts/internal/curriculumvalidator"
)

func main() {
	result, err := curriculumvalidator.Validate(".", func(message string) {
		fmt.Println(message)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if result.ErrorCount == 0 {
		if result.HasV2 {
			if result.PlaceholderCount > 0 {
				if result.LessonCount > 0 {
					fmt.Printf("Success! %d legacy lessons mapped, %d files with run commands validated, and %d v2 sections plus %d v2 items checked (%d placeholder warnings).\n",
						result.LessonCount, result.FilesScanned, result.V2SectionCount, result.V2ItemCount, result.PlaceholderCount)
				} else {
					fmt.Printf("Success! %d files with run commands validated, and %d v2 sections plus %d v2 items checked (%d placeholder warnings).\n",
						result.FilesScanned, result.V2SectionCount, result.V2ItemCount, result.PlaceholderCount)
				}
			} else {
				if result.LessonCount > 0 {
					fmt.Printf("Success! %d legacy lessons mapped, %d files with run commands validated, and %d v2 sections plus %d v2 items checked.\n",
						result.LessonCount, result.FilesScanned, result.V2SectionCount, result.V2ItemCount)
				} else {
					fmt.Printf("Success! %d files with run commands validated, and %d v2 sections plus %d v2 items checked.\n",
						result.FilesScanned, result.V2SectionCount, result.V2ItemCount)
				}
			}
			return
		}

		fmt.Printf("Success! All %d lessons mapped and %d files with run commands validated.\n", result.LessonCount, result.FilesScanned)
		return
	}

	fmt.Printf("Found %d validation errors.\n", result.ErrorCount)
	os.Exit(1)
}
