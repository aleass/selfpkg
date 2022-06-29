package main

import (
	"fmt"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:3101")
	if err != nil {
		fmt.Println("net.ResolveTCPAddr(tcp, %s) error(%v)", err)
		return
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("net.ListenTCP(tcp, %s) error(%v)", err)
		return
	}
	conn, err := listener.AcceptTCP()
	if err != nil {
		fmt.Println("net.ListenTCP(tcp, %s) error(%v)", err)
		return
	}
	var header = make([]byte, 32)
	for {
		_, err = conn.Read(header)
		if err != nil {
			fmt.Println("net.ListenTCP(tcp, %s) error(%v)", err)
			return
		}
		fmt.Println(string(header))
		_, err = conn.Write(header)
		if err != nil {
			fmt.Println("net.ListenTCP(tcp, %s) error(%v)", err)
			return
		}
	}
}
