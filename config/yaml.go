package main

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"log"
)

var yamlconfig = `
cache:
  enable : false
  list : [redis,mongoDB]
mysql:
  user : root
  password : Tech2501
  host : 10.11.22.33
  port : 3306
  name : cwi
`

type Yaml1 struct {
	SQLConf   Mysql `yaml:"mysql"`
	CacheConf Cache `yaml:"cache"`
}

// Mysql struct of mysql conf
type Mysql struct {
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
}

// Cache struct of cache conf
type Cache struct {
	Enable bool     `yaml:"enable"`
	List   []string `yaml:"list,flow"`
}

func main() {
	conf := new(Yaml1)
	err := yaml.Unmarshal([]byte(yamlconfig), conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	fmt.Println("conf", conf)
}
