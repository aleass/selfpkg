package https

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

//限速
var limit int64

//kb
const kb = 1024

func httpLimit() {
	limit = 500 * kb
	data, err := http.Get("https://down.sandai.net/mac/thunder_5.0.1.65498.dmg")
	if err != nil {
		fmt.Println(err.Error())
	}
	var buff = make([]byte, kb*10)
	var temp int64
	var f, _ = os.OpenFile("qq.exe", os.O_CREATE|os.O_RDONLY|os.O_TRUNC, 0655)
	t := time.Now().Unix()
	//测试模仿途中修改限速
	go func() {
		time.Sleep(time.Second * 36)
		atomic.StoreInt64(&limit,200 *kb)
	}()
	var now = t
	for true {
		n, err := data.Body.Read(buff)
		if err != nil {
			f.Write(buff[:n])
			data.Body.Close()
			break
		}
		f.Write(buff[:n])
		temp += int64(n)
		if temp >= limit {
			if (time.Now().Unix() - t) <= 1 {
				time.Sleep(time.Second)
				t = time.Now().Unix()
			}
			temp = 0
		}
	}
	f.Close()
	fmt.Println(time.Now().Unix()-now, limit)
}
