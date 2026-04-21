// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

const (
	Sunday = iota + 1
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type LogLevel int

const (
	LogError LogLevel = iota
	LogWarn
	LogInfo
	LogDebug
	LogFatal
)

func (l LogLevel) String() string {
	switch l {
	case LogError:
		return "ERROR"
	case LogWarn:
		return "WARN"
	case LogInfo:
		return "INFO"
	case LogDebug:
		return "DEBUG"
	case LogFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func main() {
	fmt.Println("=== Days of the Week (iota + 1) ===")
	fmt.Println("Sunday:   ", Sunday)
	fmt.Println("Monday:   ", Monday)
	fmt.Println("Tuesday:  ", Tuesday)
	fmt.Println("Wednesday:", Wednesday)
	fmt.Println("Thursday: ", Thursday)
	fmt.Println("Friday:   ", Friday)
	fmt.Println("Saturday: ", Saturday)

	fmt.Println()

	fmt.Println("=== Log Levels (type-safe enum) ===")
	fmt.Println("LogError:", LogError)
	fmt.Println("LogWarn: ", LogWarn)
	fmt.Println("LogInfo: ", LogInfo)
	fmt.Printf("LogError as int: %d\n", int(LogError))

	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: LB.4 application-logger")
	fmt.Println("Current: LB.3 (enums)")
	fmt.Println("---------------------------------------------------")
}
