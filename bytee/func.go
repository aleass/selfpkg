package bytee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"unsafe"
)

func copyss() {
	var s1 = [2][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9}}
	var s = make([]int, 9)
	var l int
	for _, v := range s1 {
		copy(s[l:], v)
		l = len(v)
	}
}

func delete() {
	var ids = []int{}
	for i := 0; i < 10; i++ {
		ids = append(ids, i)
	}
	fmt.Println(ids)
	var l = len(ids)
	for i := 0; i < l; i++ {
		if ids[i]%2 == 0 {
			ids[i], ids[l-1] = ids[l-1], ids[i]
			i--
			l--
		}
	}
	fmt.Println(ids[:l])
}

// 显式id使用
const (
	Byte32 = 0xA6C7A5CB

	Rand1      = 0x3CA5 << 16
	Rand2      = 0xA5C3
	unevenByte = 1 //奇偶判断
)

// EncodeGuid 返回id的guid
// mysql自增最大为1<<32-1,因此32-64放id
func EncodeGuid(id int64) int64 {
	now := time.Now().Unix()
	if now&unevenByte == 0 {
		now = now << 16
		return (id << 32) + int64(int32(now|Rand2))
	}
	now = int64(int16(now))
	id = (id << 32) + int64(int32(now|Rand1))
	return id
}

// DecodeGuid 返回原id
func DecodeGuid(guid int64) int64 {
	return guid >> 32
}

func EncodeGuid2(id int64) int64 {
	id = id << 32
	return id + Byte32
}

/*
oss + url
*/

var domain = []byte("https://stock-1309579221.file.myqcloud.com")
var temp = [20][1024]byte{}

func AddOssUrlSlow(data []byte) []byte {
	var temp []string
	json.Unmarshal(data, &temp)
	var ss = make([]string, len(temp))
	for i := 0; i < len(temp); i++ {
		s := temp[i]
		b := make([]byte, len(s)+len(domain))
		n := copy(b, domain)
		copy(b[n:], s)
		ss[i] = string(b)
	}
	data, _ = json.Marshal(ss)
	return data
}

func AddOssUrlFast(data []byte, num int) []byte {
	if len(data) > 1024 {
		return AddOssUrlSlow(data)
	}
	var ok bool
	var l, n int
	var bytess []byte
	if num < 20 {
		bytess = temp[num][:]
	} else {
		bytess = make([]byte, 1024)
	}
	for i, v := range data {
		if v == '"' {
			ok = true
			continue
		}
		if v == '/' && ok {
			n += copy(bytess[n:], data[l:i])
			l = i
			n += copy(bytess[n:], domain)
		}
		ok = false
	}
	n += copy(bytess[n:], data[l:])
	if bytess[n-1] == 93 { //]  说明长度大于1024
		return bytess[:n]
	}
	return AddOssUrlSlow(data)
}

type Buffer struct {
	buf []byte // contents are the bytes buf[off : len(buf)]
}

func GetRWIO(b []byte) (buff *bytes.Buffer) {
	buff = &bytes.Buffer{}
	var p = (*Buffer)(unsafe.Pointer(&buff))
	(*p).buf = b
	return
}
