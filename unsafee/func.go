package unsafee

import (
	"fmt"
	"unsafe"
)

type info struct {
	id    int
	infos *info
}

func unsafeStruct() {
	var s = info{
		id:    123,
		infos: &info{id: 233},
	}
	var p = unsafe.Pointer(&s)
	l := *(*int)(p) //获取长度
	println(l)
	infos := *(**uintptr)(unsafe.Add(p, 8))  //获取infos地址
	index1 := *(*int)(unsafe.Pointer(infos)) //根据数据地址拿第一个数据
	fmt.Println(index1)
}
func unsafeSlice() {
	var s = make([]int, 12, 18)
	s[0] = 99
	var p = unsafe.Pointer(&s)
	l := *(*int)(unsafe.Add(p, 8)) //获取长度
	println(l)
	c := *(*int)(unsafe.Add(p, 16)) //获取cap
	println(c)
	ptr := *(**uintptr)(p)                 //获取数据地址
	index1 := *(*int)(unsafe.Pointer(ptr)) //根据数据地址拿第一个数据
	fmt.Println(index1)

	arr := *(*[3]int)(unsafe.Pointer(ptr)) //转为数组展示
	fmt.Println(arr)
}
