package stringss

import (
	"fmt"
	"unsafe"
)

func demo1() {
	var s1 = [2]string{"{1, 2, 3, 4, 5, 6}", "{7, 8, 9}"}
	var s = make([]byte, 27)             //申请长度的byte
	st := *(*string)(unsafe.Pointer(&s)) //绑定同一个地址
	for _, v := range s1 {
		copy(s, v)
		s = s[len(v):]
	}
	fmt.Println(string(s))
	fmt.Println(st)
}
func BenchmarkSlice() {
	data := `hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbddddddddddddddddddddddddddddddddddddddddddddddddddddddbbbbbbbbbbbbbbbbbbbbbbbcccccccccccccccccccccccccccccccccccccccccccccccc`
	var bs = make([]byte, len(data))
	copy(bs, data)

	bs = make([]byte, len(data))
	bs = []byte(data)
}
