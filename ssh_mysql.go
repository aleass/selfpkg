package main

import (
	"app-api/model"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	sql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

type Dialer struct {
	client *ssh.Client
}

type SSH struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Port     int    `json:"port"`
	Type     string `json:"type"`
	Password string `json:"password"`
	KeyFile  string `json:"key"`
}

func (v *Dialer) Dial(address string) (net.Conn, error) {
	return v.client.Dial("tcp", address)
}

func (s *SSH) DialWithPassword() (*ssh.Client, error) {
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return ssh.Dial("tcp", address, config)
}

func (s *SSH) DialWithKeyFile() (*ssh.Client, error) {
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	config := &ssh.ClientConfig{
		User:            s.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if k, err := ioutil.ReadFile(s.KeyFile); err != nil {
		return nil, err
	} else {
		signer, err := ssh.ParsePrivateKey(k)
		if err != nil {
			return nil, err
		}
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}
	return ssh.Dial("tcp", address, config)
}

var dial *ssh.Client

func DevDb(charset string) {
	client := SSH{
		Host:    "42.192.40.225",
		User:    "guanziliang",
		Port:    22,
		KeyFile: "C:\\Users\\Administrator\\Desktop\\id_rsa_4096_new",
		Type:    "KEY", // PASSWORD or KEY
	}
	var (
		err error
	)

	switch client.Type {
	case "KEY":
		dial, err = client.DialWithKeyFile()
	case "PASSWORD":
		dial, err = client.DialWithPassword()
	default:
		panic("unknown ssh type.")
	}
	if err != nil {
		log.Fatalf("ssh connect error: %s", err.Error())
		return
	}

	// 注册ssh代理
	mysql.RegisterDial("mysql+ssh", (&Dialer{client: dial}).Dial)

	config := Configer.RedMysql

	dsn := fmt.Sprintf(
		"%s:%s@mysql+ssh(%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.DbName,
		charset,
	)
	w := log.New(os.Stdout, "\r\n", log.LstdFlags) // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
	level := logger.Silent
	if Configer.App.Debug {
		level = logger.Info
	}

	newLogger := logger.New(
		w,
		logger.Config{
			SlowThreshold:             time.Second * 6, // 慢 SQL 阈值
			LogLevel:                  level,           // 日志级别
			IgnoreRecordNotFoundError: true,            // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,           // 禁用彩色打印
		},
	)
	RedMy, err = gorm.Open(sql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("mysql 启动失败!,原因:" + err.Error())
	}

	config = Configer.StockMysql
	dsn = fmt.Sprintf(
		"%s:%s@mysql+ssh(%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.DbName,
		charset,
	)
	StockMy, err = gorm.Open(sql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("mysql 启动失败!,原因:" + err.Error())
	}

	var i int64
	err = RedMy.Model(&model.BaseBoard{}).Exec("select count(*)c from `red_bill`.`base_board` ").Count(&i).Error
	if err != nil {
		log.Fatalf("mysql connect error: %s", err.Error())
		return
	}
}
