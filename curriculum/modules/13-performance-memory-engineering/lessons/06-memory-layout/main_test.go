package main

import (
	"testing"
	"unsafe"
)

func TestBadOrderedSize(t *testing.T) {
	size := unsafe.Sizeof(BadOrdered{})
	if size != 24 {
		t.Errorf("BadOrdered size = %d; want 24", size)
	}
}

func TestGoodOrderedSize(t *testing.T) {
	size := unsafe.Sizeof(GoodOrdered{})
	if size != 16 {
		t.Errorf("GoodOrdered size = %d; want 16", size)
	}
}

func TestSizeReduction(t *testing.T) {
	badSize := unsafe.Sizeof(BadOrdered{})
	goodSize := unsafe.Sizeof(GoodOrdered{})
	if goodSize >= badSize {
		t.Errorf("GoodOrdered (%d) should be smaller than BadOrdered (%d)", goodSize, badSize)
	}
}

func TestOffsets(t *testing.T) {
	var good GoodOrdered
	if unsafe.Offsetof(good.Amount) != 0 {
		t.Errorf("Amount offset should be 0, got %d", unsafe.Offsetof(good.Amount))
	}
	if unsafe.Offsetof(good.Counter) != 8 {
		t.Errorf("Counter offset should be 8, got %d", unsafe.Offsetof(good.Counter))
	}
	if unsafe.Offsetof(good.Flag) != 12 {
		t.Errorf("Flag offset should be 12, got %d", unsafe.Offsetof(good.Flag))
	}
}

func TestUserSize(t *testing.T) {
	size := unsafe.Sizeof(User{})
	if size != 32 {
		t.Errorf("User size = %d; want 32", size)
	}
}

func TestFieldAccessWorksAfterReordering(t *testing.T) {
	u := User{
		ID:    42,
		Name:  "test",
		Score: 98.5,
	}
	if u.ID != 42 || u.Name != "test" || u.Score != 98.5 {
		t.Errorf("User fields not set correctly: %+v", u)
	}
}

type alignTest struct {
	A byte
	B int64
	C byte
}

type alignTestOptimized struct {
	B int64
	A byte
	C byte
}

func TestAlignmentPadding(t *testing.T) {
	tests := []struct {
		name   string
		typ    interface{}
		expect uintptr
	}{
		{"alignTest", alignTest{}, 16},
		{"alignTestOptimized", alignTestOptimized{}, 16},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			size := unsafe.Sizeof(tc.typ)
			if size != tc.expect {
				t.Errorf("size = %d; want %d", size, tc.expect)
			}
		})
	}
}
