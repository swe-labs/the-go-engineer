package main

import (
	"fmt"
)

func LookupSafe(data map[string]interface{}, key string, expectedType string) (interface{}, bool, error) {
	val, exists := data[key]
	if !exists {
		return nil, false, nil
	}

	switch expectedType {
	case "string":
		v, ok := val.(string)
		return v, ok, nil
	case "int":
		v, ok := val.(int)
		return v, ok, nil
	case "float64":
		v, ok := val.(float64)
		return v, ok, nil
	case "bool":
		v, ok := val.(bool)
		return v, ok, nil
	default:
		return nil, false, fmt.Errorf("unknown type: %s", expectedType)
	}
}

func main() {
	data := map[string]interface{}{
		"name":   "Alice",
		"age":    30,
		"score":  95.5,
		"active": true,
	}

	if v, ok, err := LookupSafe(data, "name", "string"); ok {
		fmt.Printf("name: %s\n", v.(string))
	} else if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("name not found")
	}

	if v, ok, err := LookupSafe(data, "age", "string"); ok {
		fmt.Printf("age as string: %s\n", v.(string))
	} else if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("age not found")
	}

	if v, ok, err := LookupSafe(data, "score", "float64"); ok {
		fmt.Printf("score: %.1f\n", v.(float64))
	} else if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("score not found")
	}

	if _, ok, err := LookupSafe(data, "missing", "string"); ok {
		fmt.Println("found missing key")
	} else if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("missing key correctly reported as not found")
	}
}
