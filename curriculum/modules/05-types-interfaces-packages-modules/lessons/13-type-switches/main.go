package main

import (
	"fmt"
	"time"
)

func describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("integer %d", x)
	case string:
		return fmt.Sprintf("string %q (len %d)", x, len(x))
	case bool:
		if x {
			return "true"
		}
		return "false"
	case time.Duration:
		return fmt.Sprintf("duration %v (%d ns)", x, int64(x))
	default:
		return fmt.Sprintf("unknown type %T", v)
	}
}

func mathSwitch(v any) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case int:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case float32:
		return float64(x * 2), nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}

func main() {
	inputs := []any{nil, 42, "hello", true, time.Second * 3, 3.14}
	for _, v := range inputs {
		fmt.Println(describe(v))
	}

	testVals := []any{int(10), int64(20), float32(3.5), float64(2.5), "hello"}
	for _, v := range testVals {
		result, err := mathSwitch(v)
		if err != nil {
			fmt.Printf("mathSwitch(%v) error: %v\n", v, err)
		} else {
			fmt.Printf("mathSwitch(%v) = %v\n", v, result)
		}
	}
}
