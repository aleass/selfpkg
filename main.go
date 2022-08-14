package main

import (
	"fmt"
	"unsafe"
)

type info struct {
	id    int
	infos *info
}

func main() {
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
