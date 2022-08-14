package testt

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestT(t *testing.T) {
	fmt.Println("t")
}

func TestT2(t *testing.T) {
	fmt.Println("t2")
}

func BenchmarkName2221(b *testing.B) {
	a := 01222
	a = 89
	_ = a
}

func FuzzConnect2(f *testing.F) {
	f.Fuzz(func(t *testing.T, a, b, c, d string) {
		fmt.Println(a, "2--------")
		var s = []string{a, b, c, d}
		var l = len(s[0])
		for i := 1; i < len(s); i++ {
			l += len(s[i])
		}
		var bytes = make([]byte, l)               //申请长度的byte
		str := *(*string)(unsafe.Pointer(&bytes)) //绑定同一个地址
		var str2 string
		for _, v := range s {
			copy(bytes, v)
			bytes = bytes[len(v):]
			str2 += v
		}

		if str != str2 {
			t.Errorf("%s != %s,{a:%s}{b:%s}{c:%s}{d:%s} ", str, str2, a, b, c, d)
		}
	})
}
func FuzzConnect(f *testing.F) {
	f.Fuzz(func(t *testing.T, a, b, c, d string) {
		fmt.Println(a, "1--------")
		var s = []string{a, b, c, d}
		var l = len(s[0])
		for i := 1; i < len(s); i++ {
			l += len(s[i])
		}
		var bytes = make([]byte, l)               //申请长度的byte
		str := *(*string)(unsafe.Pointer(&bytes)) //绑定同一个地址
		var str2 string
		for _, v := range s {
			copy(bytes, v)
			bytes = bytes[len(v):]
			str2 += v
		}

		if str != str2 {
			t.Errorf("%s != %s,{a:%s}{b:%s}{c:%s}{d:%s} ", str, str2, a, b, c, d)
		}
	})
}
