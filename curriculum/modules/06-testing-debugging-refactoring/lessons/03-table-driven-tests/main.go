package main

import (
	"fmt"
	"strings"
)

func ParseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "1":
		return true, nil
	case "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean: %q", s)
	}
}

func main() {
	fmt.Println(ParseBool("true"))
	fmt.Println(ParseBool("FALSE"))
	fmt.Println(ParseBool("1"))
	fmt.Println(ParseBool("yes"))
}
