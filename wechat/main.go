package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kbinani/screenshot"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Text struct {
	Content string `json:"content"`
}
type Image struct {
	Base64 string `json:"base64"`
	Md5    string `json:"md5"`
}

type Images struct {
	Msgtype string `json:"msgtype"`
	Image   Image  `json:"image"`
}

func main() {

	f, _ := os.Open("1.png")
	info, _ := ioutil.ReadAll(f)

	base := base64.StdEncoding.EncodeToString(info)
	hash := md5.New()
	_, _ = io.WriteString(hash, string(info))

	jsonStr := Images{
		"image",
		Image{
			Base64: base,
			Md5:    hex.EncodeToString(hash.Sum(nil)),
		},
	}
	url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=9ee81f4c-a128-43d4-bb15-996b59a03645"

	str, _ := json.Marshal(jsonStr)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(str))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	statuscode := resp.StatusCode
	fmt.Println(statuscode)
}

// save *image.RGBA to filePath with PNG format.
func save(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}
func main2() {

	////自定义截图
	//img, err := screenshot.Capture(0, 0, 500, 500)
	//if err != nil {
	//	panic(err)
	//}
	//save(img, "自定义.png")

	//获取所有活动屏幕
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		panic("没有发现活动的显示器")
	}

	//全屏截取所有活动屏幕
	if n > 0 {
		for i := 0; i < n; i++ {
			img, err := screenshot.CaptureDisplay(i)
			if err != nil {
				panic(err)
			}
			fileName := fmt.Sprintf("第%d屏幕截图.png", i)
			save(img, fileName)
		}
	}

	////使用 Rectangle 自定义截图
	//// 获取第一个屏幕显示范围
	//var new image.Rectangle = image.Rect(0, 0, 600, 700)
	//img, err = screenshot.CaptureRect(new)
	//if err != nil {
	//	panic(err)
	//}
	//save(img, "new.png")
	//
	////使用 GetDisplayBounds获取指定屏幕显示范围，全屏截图
	//bounds := screenshot.GetDisplayBounds(0)
	//img, err = screenshot.CaptureRect(bounds)
	//if err != nil {
	//	panic(err)
	//}
	//save(img, "all.png")
	//
	////使用Union获取指定屏幕 矩形最小点和最大点
	//var all image.Rectangle = image.Rect(0, 0, 0, 0)
	//bounds1 := screenshot.GetDisplayBounds(0)
	//all = bounds1.Union(all)
	//fmt.Println(all.Min.X, all.Min.Y, all.Dx(), all.Dy())

}
