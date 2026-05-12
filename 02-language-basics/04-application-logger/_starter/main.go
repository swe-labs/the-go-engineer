// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// LogLevel (Type): names the log level concept so the lesson can pass it as a first-class value.
type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
)

// levelNames (Slice): keeps ordered lesson state visible for iteration and comparison.
var levelNames = []string{"Trace", "Debug", "Info", "Warning", "Error"}

// LogLevel.String (Method): applies the string operation to receiver state at a visible boundary.
func (l LogLevel) String() string {
	if l < LevelTrace || l > LevelError {
		return "Unknown"
	}
	return levelNames[l]
}

// printLogLevel (Function): runs the print log level step and keeps its inputs, outputs, or errors visible.
func printLogLevel(level LogLevel) {
	fmt.Printf("Log level: %d %s\n", level, level.String())
}

func main() {
	printLogLevel(LevelTrace)
	printLogLevel(LevelDebug)
	printLogLevel(LevelInfo)
	printLogLevel(LevelWarning)
	printLogLevel(LevelError)

	printLogLevel(10)

	fmt.Println()
	fmt.Println("Section 02 complete! Ready for Control Flow.")
}
