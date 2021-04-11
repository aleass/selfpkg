package main

import (
	"fmt"
	"selfpkg/driver"
)

func main() {
	RedisServer, err := driver.NewRdids("127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	RedisServer.Set("a", "b", 0)
}
