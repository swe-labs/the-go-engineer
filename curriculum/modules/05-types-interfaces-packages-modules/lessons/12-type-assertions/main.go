package main

import (
	"fmt"
	"io"
	"strings"
)

func printType(v interface{}) {
	s, ok := v.(string)
	if ok {
		fmt.Println("string:", s)
		return
	}
	i, ok := v.(int)
	if ok {
		fmt.Println("int:", i)
		return
	}
	b, ok := v.(bool)
	if ok {
		fmt.Println("bool:", b)
		return
	}
	fmt.Printf("unknown type: %T\n", v)
}

func classify(items []interface{}) {
	for _, item := range items {
		switch {
		case item == nil:
			fmt.Println("nil")
		case func() bool { _, ok := item.(string); return ok }():
			s, _ := item.(string)
			fmt.Printf("string: %q\n", s)
		case func() bool { _, ok := item.(int); return ok }():
			i, _ := item.(int)
			fmt.Printf("int: %d\n", i)
		case func() bool { _, ok := item.(bool); return ok }():
			b, _ := item.(bool)
			fmt.Printf("bool: %v\n", b)
		case func() bool { _, ok := item.(float64); return ok }():
			f, _ := item.(float64)
			fmt.Printf("float64: %f\n", f)
		default:
			fmt.Println("other")
		}
	}
}

func main() {
	printType("hello")
	printType(42)
	printType(true)
	printType(3.14)

	var r io.Reader = strings.NewReader("hello")
	sr := r.(*strings.Reader)
	fmt.Printf("concrete: %T, len: %d\n", sr, sr.Len())

	if w, ok := r.(io.WriterTo); ok {
		fmt.Println("also implements WriterTo")
		_ = w
	} else {
		fmt.Println("does not implement WriterTo")
	}

	v := interface{}("must be string")
	s := v.(string)
	fmt.Println("safe assertion:", s)

	items := []interface{}{"hello", 42, true, 3.14, nil, []int{1, 2}}
	classify(items)
}
