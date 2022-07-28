package stringss

import (
	"unsafe"
)

func connectString(s []string) (str string) {
	var l = len(s[0])
	for i := 1; i < len(s); i++ {
		l += len(s[i])
	}
	var bytes = make([]byte, l)              //申请长度的byte
	str = *(*string)(unsafe.Pointer(&bytes)) //绑定同一个地址
	var str2 string
	for _, v := range s {
		copy(bytes, v)
		bytes = bytes[len(v):]
		str2 += v
	}
	return
}

func string2sliceByte(data string) (bs []byte) {
	var defaults = [32]byte{} //32以内
	var l = len(data)
	if l > 32 {
		bs = make([]byte, len(data))
	} else {
		bs = defaults[:l]
	}
	copy(bs, data)
	return
}
