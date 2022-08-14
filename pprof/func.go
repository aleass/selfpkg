package pprof

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

func a(loop bool) {
	var s string
	if !loop {
		for i := 0; i < 10000; i++ {
			s += strconv.Itoa(i)
		}
	} else {
		for i := 0; true; i++ {
			s += strconv.Itoa(i)
			time.Sleep(time.Second * 5)
		}
	}

}
func netProf() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
	a(true)
}

func runtimeProf() {
	f, _ := os.Create("./cpu3.prof")
	f2, _ := os.Create("./men3.prof")
	pprof.StartCPUProfile(f)
	a(false)
	pprof.StopCPUProfile()
	pprof.WriteHeapProfile(f2)
	f.Close()
	f2.Close()
}
