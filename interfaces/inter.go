package interfaces

type TaskInt interface {
	Put(tasks []any) bool //放入任务
	Get() int             //获取任务数量
	Run() bool            //运行任务
	Cancel() bool         //取消任务
	IsDone() error        //是否结束，阻塞。
}

type SelfError struct {
	error string //错误
	has   bool   //是否有错误
}

func (e SelfError) Error() string {
	return e.error
}

//Put 错误收集
func (e *SelfError) Put(err error) {
	if err == nil {
		return
	}
	e.has = true
	e.error += err.Error() + "\r\n"
}
