package interfaces

import (
	"selfpkg/functions"
	"strconv"
	"testing"
	"time"
)

//TestMulti 测试并发下载
func TestMulti(t *testing.T) {
	var Url = "https://dlie.sogoucdn.com/se/sogou_explorer_11.0.1.34700_0000.exe"
	tdl := NewTaskDl(100, 8, true, "./")
	var dl = make([]any, 5)
	for i := 0; i < 100; i++ {
		dl[i] = functions.FileInfo{
			Url:  Url,
			Name: strconv.Itoa(i) + ".qq.dmg",
		}
	}
	tdl.Put(dl)
	go tdl.Run()
	tdl.IsDone()
}

//TestCancel 测试取消
func TestCancel(t *testing.T) {
	var Url = "https://dlie.sogoucdn.com/se/sogou_explorer_11.0.1.34700_0000.exe"
	tdl := NewTaskDl(100, 2, true, "./")
	var dl = make([]any, 5)
	for i := 0; i < 3; i++ {
		dl[i] = functions.FileInfo{
			Url:  Url,
			Name: strconv.Itoa(i) + ".qq.dmg",
		}
	}
	tdl.Put(dl)
	go tdl.Run()
	go func() {
		time.Sleep(time.Second * 3)
		println("cancel !!!")
		println("tasks : ", tdl.Get())
		tdl.Cancel()
	}()
	tdl.IsDone()
}
