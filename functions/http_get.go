package functions

import (
	"context"
	"io"
	"net/http"
	"os"
)

// get url 获取

var buffer = 32 * 1024

type FileInfo struct {
	Url    string
	Name   string
	buffer int
}

func GetUrl(urls, paths string, cancel context.Context) error {
	//client := &http.Client{}
	//req, err := http.NewRequest("get", urls, nil)
	//req.Header.Add("Accept", `*/*`)
	//req.Header.Add("Accept-Encoding", `identity`)
	//req.Header.Add("Host", `d5.wzip.ru`)
	//req.Header.Add("User-Agent", `Wget/1.21.3`)
	//req.Header.Add("Connection", `Keep-Alive`)
	//resp, err := client.Do(req)
	resp, err := http.Get(urls)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	out, err := os.Create(paths)
	if err != nil {
		return err
	}
	defer out.Close()

	if buffer < 0 {
		buffer = 32 * 1024 //32kb
	}
	buffers := make([]byte, buffer)

	// 循环读取文件内容 此处为了及时取消，所以手动拷贝。
	for true {
		select {
		case <-cancel.Done():
			os.Remove(paths)
			return cancel.Err()
		default:
			rn, err := resp.Body.Read(buffers)
			if io.EOF == err {
				return nil
			}
			if nil != err {
				return err
			}
			_, err = out.Write(buffers[:rn])
			if err == io.EOF { // 读到文件末尾就不再往后读取
				return nil
			}
		}
	}

	return nil
}
