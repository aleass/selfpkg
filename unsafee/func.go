package unsafee

import (
	"fmt"
	"unsafe"
)

func pointers() {
	var a = []int{1898, 2, 3}
	fmt.Println(**(**int)(unsafe.Pointer(&a)))
	fmt.Println(*(*int)(unsafe.Add(unsafe.Pointer(&a), 8)))
}
