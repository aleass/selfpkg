package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var s sync.WaitGroup
var num int64

func MulSse(count int) {
	go func() {
		for true {
			fmt.Println("count:", num)
			time.Sleep(time.Second * 1)
		}
	}()
	s.Add(count)
	for i := 0; i < count; i++ {
		go seeClient()
		time.Sleep(time.Millisecond * 1)
	}
	s.Wait()
	fmt.Println(num)
}

func sendData(ctx *gin.Context) {
	_, err := ctx.Writer.Write([]byte("str"))
	if err != nil {
		fmt.Println(err.Error())
	}
	ctx.Writer.Flush()
	//m := map[string]interface{}{
	//	"types": "data",
	//}
	//jsonBytes, _ := json.Marshal(m)
	//// 测试加一个data层
	//jsonStr := string(jsonBytes)
	//if len(jsonStr) != 0 {
	//	str := "data:" + jsonStr + "\n\n"
	//	_, err := ctx.Writer.Write([]byte(str))
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//	ctx.Writer.Flush()
	//}
}
func sseServer() {
	gin.SetMode("debug")
	r := gin.Default()
	r.GET("/a", func(ctx *gin.Context) {
		ctx.Header("X-Accel-Buffering", "no")
		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		close := ctx.Request.Context()
		for i := 0; i < 9999; i++ {
			select {
			case <-close.Done():
				return
			default:
			}
			sendData(ctx)
			time.Sleep(time.Second * 1)
		}
	})
	r.Run("127.0.0.1:" + "1880")
}

func seeClient() {
	defer s.Done()
	//url := "https://29.push2.eastmoney.com/api/qt/stock/trends2/sse?secid=1.600705&ndays=1&ut=fa5fd1943c7b386f172d6893dbfba10b&fields1=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13&fields2=f51,f52,f53,f54,f55,f56,f57,f58"
	url := "http://127.0.0.1:1880/sse/a/index/all?type=key&stock=sz-000031,sz-000041,sz-000032"
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error() + "******" + strconv.Itoa(int(num)))
		return
	}
	atomic.AddInt64(&num, 1)
	buf := make([]byte, 12563)
	//读取 Reader 对象中的内容到 []byte 类型的 buf 中
	reader := bufio.NewReader(res.Body)
	info, err := reader.Read(buf)
	for err == nil && info != 0 {
		info, err = reader.Read(buf)
	}
	res.Body.Close()
	if info == 0 {
		fmt.Println("client close----------------------")
	}
	if err != nil {
		fmt.Println(err.Error() + "******" + strconv.Itoa(int(num)))
		return
	}
}
