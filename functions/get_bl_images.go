package functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ImageStruct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List struct {
			Vlist []struct {
				Pic string `json:"pic"`
			} `json:"vlist"`
		} `json:"list"`
	} `json:"data"`
}

func ImagePage(maxPage, mid int) {
	//ImagePage(23, 39101587)
	if maxPage == 0 || mid == 0 {
		return
	}
	var images = make([]string, 0, 1000)
	var urls = "https://api.bilibili.com/x/space/arc/search?mid=%d&ps=30&tid=0&pn=%d&keyword=&order=pubdate&jsonp=jsonp"
	for i := 1; i < maxPage+1; i++ {
		_urls := fmt.Sprintf(urls, mid, i)
		res, err := http.Get(_urls)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		response, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		var is ImageStruct
		json.Unmarshal(response, &is)
		if is.Code != 0 {
			fmt.Println(is.Code, is.Message, i)
			return
		}
		for _, image := range is.Data.List.Vlist {
			images = append(images, image.Pic)
		}
		time.Sleep(time.Second)
	}

	if len(images) == 0 {
		fmt.Println("images is 0 ")
		return
	}

	bytes, _ := json.Marshal(images)
	err := ioutil.WriteFile("a.txt", bytes, 0755)
	if err != nil {
		fmt.Println(err.Error())
	}

	for i, v := range images {
		ioutil.WriteFile("images/"+strconv.Itoa(i)+".jpg", []byte(v), 0777)
		time.Sleep(time.Millisecond * 500)
	}
}
