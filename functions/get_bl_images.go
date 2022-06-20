package functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	VideoImage = "https://api.bilibili.com/x/space/arc/search?mid=%d&ps=30&tid=0&pn=%d&keyword=&order=pubdate&jsonp=jsonp"
	ImageAlbum = "https://api.bilibili.com/x/dynamic/feed/draw/doc_list?uid=%d&page_num=%d&page_size=30&biz=all&jsonp=jsonp"
)

type VideoImageStruct struct {
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
type ImageAlbumStruct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Items []struct {
			Pictures []struct {
				ImgSrc string `json:"img_src"`
			} `json:"pictures"`
		} `json:"items"`
	} `json:"data"`
}

//ImagePage 当是相册时pages：0，视频封面图片pages：1
func ImagePage(maxPage, mid int) {
	if maxPage == 0 || mid == 0 {
		return
	}
	var images = make([]string, 0, 1000)

	for i := 1; i < maxPage+1; i++ {
		_urls := fmt.Sprintf(VideoImage, mid, i)
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
		var is VideoImageStruct
		json.Unmarshal(response, &is)
		if is.Code != 0 {
			fmt.Println(is.Code, is.Message, i)
			return
		}
		for _, image := range is.Data.List.Vlist {
			images = append(images, image.Pic)
		}
		time.Sleep(time.Millisecond * 500)

	}

	if len(images) == 0 {
		fmt.Println("images is 0 ")
		return
	}
	save("/Users/tuski/Downloads/temp/", images, true)
}

func ImageAlbumFc(maxPage, mid int) {
	if maxPage == 0 || mid == 0 {
		return
	}
	var images = make([]string, 0, 1000)

	for i := 2; i < maxPage; i++ {
		_urls := fmt.Sprintf(ImageAlbum, mid, i)
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
		var is ImageAlbumStruct
		json.Unmarshal(response, &is)
		if is.Code != 0 {
			fmt.Println(is.Code, is.Message, i)
			return
		}
		for _, multiImages := range is.Data.Items {
			for _, image := range multiImages.Pictures {
				images = append(images, image.ImgSrc)
			}
		}
		time.Sleep(time.Millisecond * 500)
	}

	if len(images) == 0 {
		fmt.Println("images is 0 ")
		return
	}
	save("/Users/tuski/Downloads/temp/", images, true)
}

func save(path string, images []string, isCopy bool) {
	bytes, _ := json.Marshal(images)
	var err error
	if isCopy {
		err = ioutil.WriteFile("a.txt", bytes, 0755)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	for i, urls := range images {
		res, err := http.Get(urls)
		if err != nil {
			fmt.Println(err.Error(), urls)
			continue
		}
		datas, _ := ioutil.ReadAll(res.Body)
		err = ioutil.WriteFile(path+strconv.Itoa(i)+".jpg", datas, 0777)
		if err != nil {
			fmt.Println(err.Error())
		}
		res.Body.Close()
		time.Sleep(time.Millisecond * 500)
	}
}
