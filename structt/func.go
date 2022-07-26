package structt

import (
	"fmt"
	"unsafe"
)

type slices struct {
	p unsafe.Pointer //bbbbb
	l int
	c int
}

func main() {
	var s2 = []int{7, 98569}
	fmt.Println(s2, len(s2), cap(s2))
	s1 := (*slices)(unsafe.Pointer(&s2))
	fmt.Println(**(**int)(unsafe.Pointer(&s1.p)))

	s1 = *(**slices)(unsafe.Pointer(&s2))
	fmt.Println(*(*int)(unsafe.Pointer(&s1.p)))
	fmt.Println(*(*int)(unsafe.Pointer(s1)))
	fmt.Println(*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(s1)) + 8)))
}

/*
   [] 0 0
   7
   7
   7
   98569
*/
