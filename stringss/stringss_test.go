package stringss

import (
	"testing"
	"unsafe"
)

func TestCopyString(t *testing.T) {
	var orig = "hello world"
	var defaults = [32]byte{}
	var l = len(orig)
	var bs []byte
	if l > 32 {
		bs = make([]byte, len(orig))
	} else {
		bs = defaults[:l]
	}

	copy(bs, orig)
	if orig != string(bs) {
		t.Errorf("%s != %s ", orig, string(bs))
	}
}
func FuzzCopyString(f *testing.F) {
	f.Fuzz(func(t *testing.T, orig string) {
		var defaults = [32]byte{}
		var l = len(orig)
		var bs []byte
		if l > 32 {
			bs = make([]byte, len(orig))
		} else {
			bs = defaults[:l]
		}
		copy(bs, orig)
		if orig != string(bs) {
			t.Errorf("%s != %s ", orig, string(bs))
		}
	})
}
func BenchmarkLinkSlice(b *testing.B) {
	var orig = "hello world" //size < 32
	//orig = "display: block;-webkit-user-select: none;margin: auto;cursor: zoom-in;background-color: hsl(0, 0%, 90%);" //size > 32
	b.Run("copy", func(b *testing.B) {
		b.ReportAllocs()
		var defaults = [32]byte{}
		var l = len(orig)
		var bs []byte
		if l > 32 {
			bs = make([]byte, len(orig))
		} else {
			bs = defaults[:l]
		}
		for i := 0; i < b.N; i++ {
			copy(bs, orig)
		}
		_ = bs
	})

	b.Run("[]byte", func(b *testing.B) {
		b.ReportAllocs()
		var bs []byte
		for i := 0; i < b.N; i++ {
			bs = []byte(orig)
		}
		_ = bs
	})
}

func TestConnect(t *testing.T) {
	a, b, c, d := "0", "", "", ""
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
		t.Errorf("%s != %s ", str, str2)
	}
}

func FuzzConnect(f *testing.F) {
	f.Fuzz(func(t *testing.T, a, b, c, d string) {
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
