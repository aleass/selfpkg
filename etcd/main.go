package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	c "go.etcd.io/etcd/client/v3/concurrency"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

func main() {

}

func distributedTes() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 3 * time.Second,
	})
	defer cli.Close()
	if err != nil {
		log.Fatal(err)
	}
	//前缀锁
	var lockKey = "lock_num"
	var num = 500
	for i := 0; i < num; i++ {
		response, e := cli.Grant(context.Background(), 30)
		if e != nil {
			log.Fatal(e.Error())
		}
		//newsession 自动keep live租约
		s, err := c.NewSession(cli, c.WithLease(response.ID))
		if err != nil {
			log.Fatal("new session err:", err)
		}
		m2 := c.NewMutex(s, lockKey)
		err = m2.Lock(context.Background())
		if err != nil {
			fmt.Println("lock err:", err.Error())
		}
		data, _ := ioutil.ReadFile("temp")
		count, _ := strconv.Atoi(string(data))
		count++
		ioutil.WriteFile("temp", []byte(strconv.Itoa(count)), 0644)
		err = m2.Unlock(context.Background()) //此处会删除锁key
		if err != nil {
			fmt.Println("unlock err:", err.Error())
		}
		s.Close() //此处会删除锁key
	}
}
