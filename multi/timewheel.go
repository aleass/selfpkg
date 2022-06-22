package multi

import (
	"fmt"
	"time"
)

type taskList struct {
	f     []func()
	count int
}

type timeWheel [60]*taskList

func NewTimeWheel() *timeWheel {
	return &timeWheel{}
}
func (t *timeWheel) Add(index int, fl []func()) {
	t[index] = &taskList{f: fl}
}
func (t *timeWheel) Run() {
	fmt.Println("start :", time.Now().Unix())
	var index int
	for true {
		if cur := t[index]; cur != nil {
			go func() {
				for _, f := range cur.f {
					f()
				}
				cur.count++
			}()
		}
		index++
		if index == 59 {
			index = 0
		}
		time.Sleep(time.Second)
	}
}
