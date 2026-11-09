package main

import (
	"fmt"
	"unsafe"
)

type BadOrdered struct {
	Flag    bool    // 1 byte + 7 padding
	Amount  float64 // 8 bytes
	Counter int32   // 4 bytes + 4 padding
} // total: 24 bytes

type GoodOrdered struct {
	Amount  float64 // 8 bytes
	Counter int32   // 4 bytes
	Flag    bool    // 1 byte + 3 padding
} // total: 16 bytes

type User struct {
	ID    int64   // 8 bytes
	Name  string  // 16 bytes (pointer + len on 64-bit)
	Score float64 // 8 bytes
} // total: 32 bytes

func main() {
	fmt.Println("BadOrdered size:", unsafe.Sizeof(BadOrdered{}))
	fmt.Println("GoodOrdered size:", unsafe.Sizeof(GoodOrdered{}))

	var bad BadOrdered
	var good GoodOrdered

	fmt.Printf("BadOrdered offsets: Flag=%d, Amount=%d, Counter=%d\n",
		unsafe.Offsetof(bad.Flag), unsafe.Offsetof(bad.Amount), unsafe.Offsetof(bad.Counter))

	fmt.Printf("GoodOrdered offsets: Amount=%d, Counter=%d, Flag=%d\n",
		unsafe.Offsetof(good.Amount), unsafe.Offsetof(good.Counter), unsafe.Offsetof(good.Flag))

	fmt.Println("User size:", unsafe.Sizeof(User{}))
}
