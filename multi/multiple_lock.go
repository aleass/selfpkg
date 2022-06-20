package multi

import "sync"

const bucketLen = 10 //default count of lock

//分布锁
type date struct {
	info interface{}
}

type DataManger struct {
	live int
	*sync.Mutex
	datas []*date
}

//MulLockDataType LockData[key] = []*DataInfo
type MulLockDataType [bucketLen]sync.Map

func NewMulLockData() *MulLockDataType {
	return &MulLockDataType{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}}
}

func (m *MulLockDataType) Save(key, val string) (index int) {
	var curBucket uint8
	//获取桶号
	if key != "" {
		curBucket = key[len(key)-1] % bucketLen
	}
	oldDataInter, ok := m[curBucket].Load(key)
	if !ok {
		temp := DataManger{
			live:  1,
			Mutex: &sync.Mutex{},
			datas: []*date{{
				val,
			}},
		}
		m[curBucket].Store(key, temp)
		index = 0
	} else {
		oldData, ok := oldDataInter.(DataManger)
		if !ok {
			index = -1
			return
		}
		oldData.Lock()
		defer oldData.Unlock()
		oldData.datas = append(oldData.datas, &date{
			info: val,
		})
		oldData.live++
		m[curBucket].Store(key, oldData)
		index = len(oldData.datas) - 1
	}
	return
}

//Del delete data with key,update key to nil
/*
@ key string
@ index int
*/
func (m *MulLockDataType) Del(key string, index int) {
	var curBucket uint8
	//获取桶号
	if key != "" {
		curBucket = key[len(key)-1] % bucketLen
	}
	oldDataInter, ok := m[curBucket].Load(key)
	if !ok {
		return
	}
	oldData, ok := oldDataInter.(DataManger)
	if !ok {
		return
	}
	oldData.Lock()
	defer oldData.Unlock()
	if len(oldData.datas) < index {
		panic("index error")
	}
	oldData.datas[index] = nil
}

//Add a key for val,if exist a nil data,replace it.
func (m *MulLockDataType) Add(key string, val interface{}) int {
	var curBucket uint8
	//获取桶号
	if key != "" {
		curBucket = key[len(key)-1] % bucketLen
	}
	oldDataInter, ok := m[curBucket].Load(key)
	if !ok {
		return -1
	}
	oldData, ok := oldDataInter.(DataManger)
	if !ok {
		return -1
	}
	oldData.live++ //live count incr
	oldData.Lock() //lock data
	defer oldData.Unlock()
	//check have a nil key
	for i, v := range oldData.datas {
		if v == nil {
			oldData.datas[i] = &date{val}
			m[curBucket].Store(key, oldData)
			return i
		}
	}

	//no nil key
	oldData.datas = append(oldData.datas, &date{val})
	m[curBucket].Store(key, oldData)
	return len(oldData.datas) - 1
}
