package multi

import "time"

var sequenceId int64
var workId int64
var datacenterId int64

func init() {
	workId = 1
	datacenterId = 1
}

func SnowFlake() (snowflake int64) {
	if sequenceId == 4096 {
		sequenceId = 0
	}
	milli := time.Now().UnixMilli()
	snowflake += milli << 22        //毫秒
	snowflake += datacenterId << 12 //datacenter id
	snowflake += workId << 17       // worker id
	snowflake += sequenceId         //sequence id
	sequenceId++
	return
}
