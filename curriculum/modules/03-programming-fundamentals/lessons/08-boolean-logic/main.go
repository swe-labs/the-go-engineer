package main

import "fmt"

func canAccess(role string, isAdmin, isActive bool, hour int) bool {
	if role == "" {
		return false
	}
	if role == "banned" {
		return false
	}
	if !isActive {
		return false
	}
	if isAdmin {
		return true
	}
	return hour >= 9 && hour < 17
}

func main() {
	cases := []struct {
		role     string
		isAdmin  bool
		isActive bool
		hour     int
	}{
		{"admin", true, true, 3},
		{"admin", true, false, 3},
		{"user", false, true, 10},
		{"user", false, true, 20},
		{"user", false, false, 10},
		{"banned", false, true, 10},
		{"", false, true, 10},
		{"guest", false, true, 9},
		{"guest", false, true, 17},
	}
	for _, c := range cases {
		access := canAccess(c.role, c.isAdmin, c.isActive, c.hour)
		fmt.Printf("role=%q admin=%t active=%t hour=%d → %t\n",
			c.role, c.isAdmin, c.isActive, c.hour, access)
	}
}
