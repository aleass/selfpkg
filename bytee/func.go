package bytee

import (
	"fmt"
	"time"
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

func delete(){
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
