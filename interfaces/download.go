package interfaces

import (
	"context"
	"fmt"
	"selfpkg/functions"
	"sync"
	"time"
)

//任务下载对应实现接口

type TaskDownload struct {
	multi    uint                     //并发量
	tChan    chan *functions.FileInfo //任务存储
	filePath string                   //保存地址
	ct       context.Context          //取消信号量
	fc       func()                   //取消方法
	Done     chan struct{}            //完成信号量
	isAll    bool                     //是否一次性下载
	error    SelfError                //错误信息
	outTime  time.Duration            //超时
}

func NewTaskDl(max, multi int, isAll bool, filePath string, ot time.Duration) TaskInt {
	ct, f := context.WithCancel(context.Background())
	return &TaskDownload{
		filePath: filePath,
		tChan:    make(chan *functions.FileInfo, max),
		multi:    uint(multi),
		ct:       ct,
		fc:       f,
		isAll:    isAll,
		Done:     make(chan struct{}),
		error:    SelfError{},
		outTime:  ot,
	}
}

func (t *TaskDownload) Run() bool {
	if t.tChan == nil {
		panic("task chan is nil")
	}
	//if t.outTime > 0 {
	//	t.SetOT()
	//}
	var w sync.WaitGroup
	w.Add(int(t.multi))
	for i := uint(0); i < t.multi; i++ {
		go func() {
			defer w.Done()
			for true {
				//一次性输出并且数量为0
				if t.isAll && len(t.tChan) == 0 {
					return
				}
				select {
				case u := <-t.tChan:
					err := functions.GetUrl(u.Url, t.filePath+u.Name, t.ct, t.outTime)
					if err != nil {
						t.error.Put(err)
						//return
					}
					fmt.Println(fmt.Sprintf("not finish: %d", len(t.tChan)))
				case <-t.ct.Done(): //当为false，可以及时停止任务。
					t.error.Put(t.ct.Err())
					return
				}
			}
		}()
	}
	w.Wait()
	t.Done <- struct{}{}
	return true
}

func (t *TaskDownload) Get() (c int) {
	if t.tChan == nil {
		panic("task chan is nil")
	}
	return len(t.tChan)
}

func (t *TaskDownload) Cancel() bool {
	if t.tChan == nil {
		panic("task chan is nil")
	}
	t.fc()
	return true
}

func (t *TaskDownload) Put(task []any) bool {
	if t.tChan == nil {
		panic("task chan is nil")
	}
	for _, ts := range task {
		var f, ok = ts.(functions.FileInfo)
		if !ok {
			continue
		}
		t.tChan <- &f
	}
	return true
}

func (t *TaskDownload) IsDone() error {
	<-t.Done
	if !t.error.has {
		return nil
	}
	return t.error
}

func (t *TaskDownload) SetOT() {
	if t.outTime <= 0 {
		return
	}
	t.ct, t.fc = context.WithTimeout(t.ct, t.outTime)
}
