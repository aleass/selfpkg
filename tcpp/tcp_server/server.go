package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:3101")
	if err != nil {
		fmt.Println("ResolveTCPAddr err:", err.Error())
		return
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("ListenTCP err:", err.Error())
		return
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("AcceptTCP err:", err.Error())
			return
		}
		if err = conn.SetKeepAlive(true); err != nil {
			fmt.Println("ListenTCP err:", err.Error())
			return
		}
		if err = conn.SetReadBuffer(4096); err != nil {
			fmt.Println("ListenTCP err:", err.Error())
			return
		}
		if err = conn.SetWriteBuffer(4096); err != nil {
			fmt.Println("ListenTCP err:", err.Error())
			return
		}
		go func() {
			for {
				data, err := ReadTcp(conn)
				if err != nil {
					if err != io.EOF {
						fmt.Println("ReadTcp err:", err.Error())
					}
					return
				}
				fmt.Println(string(data))
			}
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	fmt.Println(<-c)
}

func ReadTcp(conn *net.TCPConn) ([]byte, error) {
	var header = make([]byte, tcpp.RawHeaderSize)
	_, err := conn.Read(header)
	if err != nil {
		return nil, err
	}
	packSize := binary.BigEndian.Uint32(header[tcpp.HeaderOffset:])
	//ver := binary.BigEndian.Uint16(header[nets.VerOffset:])
	dataSize := packSize - tcpp.RawHeaderSize
	if dataSize <= 0 {
		return nil, errors.New("")

	}
	var data = make([]byte, dataSize)
	_, err = conn.Read(data)
	return data, err
}
