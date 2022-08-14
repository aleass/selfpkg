package main

import (
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)

func a() {
	var s string
	for i := 0; true; i++ {
		s += strconv.Itoa(i)
		time.Sleep(time.Second * 5)
	}
}
func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
	a()

}
