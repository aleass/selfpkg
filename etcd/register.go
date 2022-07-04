package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

//the detail of service 定义服务结构，唯一id加ip地址
type ServiceInfo struct {
	ID uint64
	IP string
}

type Service struct {
	Name       string
	ServerPath string
	Info       ServiceInfo
	stop       chan error
	leaseid    clientv3.LeaseID
	client     *clientv3.Client
}

func main() {

	go func() {
		server, err := NewService("r1", ServiceInfo{
			ID: 0,
			IP: "127.0.0.1",
		}, []string{"127.0.0.1:2379"})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		server.Start()
	}()
	time.Sleep(time.Second * 2)
	dis, err := NewMaster([]string{"127.0.0.1:2379"}, "r1")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dis.WatchNodes()
}

func NewService(name string, info ServiceInfo, endpoints []string) (*Service, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Service{
		Name:       name,
		ServerPath: "services/",
		Info:       info,
		stop:       make(chan error),
		client:     cli,
	}, err
}

func (s *Service) Start() error {
	ch, err := s.keepAlive()
	if err != nil {
		log.Fatal(err)
		return err
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case err := <-c:
			s.clear()
			return errors.New(err.String())
		case err := <-s.stop:
			s.clear()
			return err
		case <-s.client.Ctx().Done():
			return s.clear()
		case ka, ok := <-ch:
			if !ok {
				log.Println("keep alive channel closed")
				s.clear()
				return nil
			}
			log.Printf("Recv reply from service: %s, ttl:%d", s.Name, ka.TTL)

		}
	}
}

func (s *Service) Stop() {
	s.stop <- nil
}

func (s *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	info := &s.Info
	key := s.ServerPath + s.Name
	value, _ := json.Marshal(info)

	// minimum lease TTL is 5-second
	resp, err := s.client.Grant(context.Background(), 30)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = s.client.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	s.leaseid = resp.ID

	return s.client.KeepAlive(context.TODO(), resp.ID)
}

func (s *Service) clear() error {
	_, err := s.client.Revoke(context.TODO(), s.leaseid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("servide:" + strconv.Itoa(int(s.leaseid)) + "stop lease:")

	_, err = s.client.Delete(context.TODO(), s.ServerPath, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("delete servide all ")

	s.client.Close()
	fmt.Println("season sto ")
	return err

}
