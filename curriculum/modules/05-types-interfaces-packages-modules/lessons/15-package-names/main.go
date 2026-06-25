package main

import (
	"fmt"

	"github.com/rasel9t6/the-go-engineer/curriculum/modules/05-types-interfaces-packages-modules/lessons/15-package-names/mathutil"
	"github.com/rasel9t6/the-go-engineer/curriculum/modules/05-types-interfaces-packages-modules/lessons/15-package-names/strutil"
)

func main() {
	a, b := 7, 3
	fmt.Printf("%d + %d = %d\n", a, b, mathutil.Add(a, b))
	fmt.Printf("%d * %d = %d\n", a, b, mathutil.Mul(a, b))
	fmt.Println(strutil.Upper("hello, go!"))
}
