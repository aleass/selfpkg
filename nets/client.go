package main

import (
	"encoding/binary"
	"fmt"
	"time"

	//"github.com/Terry-Mao/goim/pkg/encoding/binary"
	"net"
)

type s3 struct {
}

func (s s3) Read(p []byte) (n int, err error) {
	p = []byte("1234567")
	return len(p), nil
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3101")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var header = make([]byte, 16)

	binary.BigEndian.PutUint32(header[_packOffset:], 38)
	binary.BigEndian.PutUint16(header[_headerOffset:], uint16(16))
	binary.BigEndian.PutUint16(header[6:8], uint16(169))
	binary.BigEndian.PutUint32(header[8:12], uint32(7))
	binary.BigEndian.PutUint32(header[12:16], uint32(999))
	header = append(header, []byte("fjdskfjsafjdsfjlsdfjls")...)
	_, err = conn.Write(header)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go func() {
		for {
			_, err = conn.Write(header)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			time.Sleep(time.Second)
		}
	}()
	var headers = make([]byte, 4096)
	for {
		num, err := conn.Read(headers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(headers), num)
	}

}

const (
	MaxBodySize = int32(1 << 12)

	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_heartOffset  = _seqOffset + _seqSize
)
