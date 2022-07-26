package main

import (
	"encoding/binary"
	"fmt"
	"selfpkg/nets"
	//"github.com/Terry-Mao/goim/pkg/encoding/binary"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3101")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var data = []byte("hello world")
	packSize := tcpp.RawHeaderSize + len(data)
	var packData = make([]byte, packSize)

	binary.BigEndian.PutUint32(packData[tcpp.HeaderOffset:], uint32(packSize))
	binary.BigEndian.PutUint16(packData[tcpp.VerOffset:], uint16(17895))
	copy(packData[tcpp.RawHeaderSize:], data)

	_, err = conn.Write(packData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn.Close()
}
