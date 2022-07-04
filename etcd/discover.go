package main

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type Master struct {
	Path       string
	ServerPath string
	Nodes      map[string]*Node
	Client     *clientv3.Client
}

//node is a client
type Node struct {
	State bool
	Key   string
	Info  map[uint64]*ServiceInfo
}

func NewMaster(endpoints []string, watchPath string) (*Master, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	master := &Master{
		Path:       watchPath,
		ServerPath: "services/",
		Nodes:      make(map[string]*Node),
		Client:     cli,
	}
	master.AddAllServer()
	go master.WatchNodes()
	return master, err
}

func (m *Master) AddAllServer() {
	res, err := m.Client.Get(context.TODO(), m.ServerPath, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, v := range res.Kvs {
		info := &ServiceInfo{}
		err := json.Unmarshal(v.Value, info)
		if err != nil {
			log.Fatal(err)
			continue
		}
		m.AddNode(v.String(), info)
	}

}
func (m *Master) AddNode(key string, info *ServiceInfo) {
	node, ok := m.Nodes[key]
	if !ok {
		node = &Node{
			State: true,
			Key:   key,
			Info:  map[uint64]*ServiceInfo{info.ID: info},
		}
		m.Nodes[node.Key] = node
	} else {
		node.Info[info.ID] = info
	}
}

func (m *Master) DeleteNode(key string, info *ServiceInfo) {
	node, ok := m.Nodes[key]
	if !ok {
		return
	} else {
		delete(node.Info, info.ID)
	}
}

func GetServiceInfo(ev *clientv3.Event) *ServiceInfo {
	info := &ServiceInfo{}
	err := json.Unmarshal(ev.Kv.Value, info)
	if err != nil {
		log.Println(err)
	}
	return info
}

func (m *Master) WatchNodes() {
	rch := m.Client.Watch(context.Background(), m.Path, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				fmt.Printf("[%s] %q : %qn", ev.Type, ev.Kv.Key, ev.Kv.Value)
				info := GetServiceInfo(ev)
				m.AddNode(string(ev.Kv.Key), info)
			case clientv3.EventTypeDelete:
				fmt.Printf("[%s] %q : %qn", ev.Type, ev.Kv.Key, ev.Kv.Value)
				info := GetServiceInfo(ev)
				m.DeleteNode(string(ev.Kv.Key), info)
			}
		}
	}
}
