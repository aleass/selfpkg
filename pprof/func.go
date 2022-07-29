package pprof

import (
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"sync"
)

func main() {
	//model := &StockSimple{
	//	Name:         "天顺风能",
	//	Code:         "sz-002531",
	//	SimpleName:   "",
	//	ExchangeCode: 2,
	//	CreateTime:   0,
	//	UpdateTime:   0,
	//	Kind:         1,
	//}
	//Draw(model.Code, 17.4, model)
	//return

	f, _ := os.Create("./cpu3.prof")
	f2, _ := os.Create("./men3.prof")
	pprof.StartCPUProfile(f)

	var sy = sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		sy.Add(1)
		go func() {
			model := &StockSimple{
				Name:         "天顺风能",
				Code:         "sz-002531",
				SimpleName:   "",
				ExchangeCode: 2,
				CreateTime:   0,
				UpdateTime:   0,
				Kind:         1,
			}
			Draw(model.Code, 17.4, model)
			sy.Done()
		}()
	}
	sy.Wait()
	pprof.StopCPUProfile()
	pprof.WriteHeapProfile(f2)
	f.Close()
	f2.Close()
}
