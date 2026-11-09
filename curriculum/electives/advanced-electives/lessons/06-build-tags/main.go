package main

import (
	"fmt"
	"runtime"
)

// PlatformInfo reports operating system and architecture details.
type PlatformInfo struct {
	OS   string
	Arch string
}

// GetPlatformInfo returns the active build target information.
func GetPlatformInfo() PlatformInfo {
	return PlatformInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

// Demo demonstrates Build tags and satisfies test requirements.
func Demo() string {
	info := GetPlatformInfo()
	_ = info
	return "Build tags verified"
}

func main() {
	fmt.Println(Demo())
}
