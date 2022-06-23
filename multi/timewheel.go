package multi

import (
	"fmt"
	"time"
)

type timewheels struct {
	Wheel [3][60]*taskList //second min hour
}
type taskList struct {
	f     []func()
	count int
}

func NewTimeWheel() *timewheels {
	return &timewheels{}
}
func (t *timewheels) Add(index int, fl []func(), types uint) {
	if types >= 3 {
		return
	}
	t.Wheel[types][index].f = append(t.Wheel[types][index].f, fl...)
}
func (t timewheels) Run() {
	fmt.Println("start :", time.Now().Unix())
	var sed, min, hour int
	for true {
		if cur := t.Wheel[0][sed]; cur != nil {
			go func() {
				for _, f := range cur.f {
					f()
				}
				cur.count++
			}()
		}

		sed++
		if sed == 59 {
			min++
		}

		if min == 59 {
			hour++
		}

		time.Sleep(time.Second)
	}
}
