package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type ThreadItem struct {
	url     string
	content string
	imgs    string
}

var (
	//图片正则表达式
	imageItemExp = regexp.MustCompile(`src="//i\.4cdn\.org/s/[0123456789]+s\.jpg"`)
	//帖子路径正则表达式
	threadItemExp = regexp.MustCompile(`"thread/[0123456789]+"`)
)

func main() {
	content, err2 := httpGet("http://boards.4chan.org/s/")
	if err2 != 200 {
		fmt.Printf("err2==", err2)
		return
	}
	items := chooseImageItem(content)
	// items := chooseThreadItem(content)
	err := os.Mkdir("mydir", os.ModePerm)

	if err != nil {
		fmt.Printf("\nerr=", err)
		// return
	}
	//写入页面文件
	// for _, mod := range items {
	// 	fileName := mod.imgs
	// 	nameSlice := strings.Split(fileName, "/")
	// 	i := 0
	// 	for i, _ = range nameSlice {

	// 	}

	// 	file, err := os.Create("mydir" + "/" + nameSlice[i])
	// 	if err != nil {
	// 		fmt.Printf("\nerr==", err)
	// 		return
	// 	}

	// 	file.WriteString(content)
	// }
	//写入图片文件
	for _, mod := range items {
		fileName := mod.imgs
		nameSlice := strings.Split(fileName, "/")
		i := 0
		for i, _ = range nameSlice {
		}
		file, err := os.Create("mydir" + "/" + nameSlice[i])
		if err != nil {
			fmt.Println(err)
			return
		}
		imgByte := httpImgGet(fileName)
		go file.Write(imgByte)
	}
}

func httpGet(url string) (content string, statusCode int) {
	resp, err := http.Get(url)
	if err != nil {
		statusCode = -100
		return
	}
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		statusCode = -200
		return
	}
	fmt.Print("这里是200")
	fmt.Print(resp)
	statusCode = resp.StatusCode
	content = string(data)
	return
}

func httpImgGet(url string) (data []byte) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("img Err get", err, url)
	}
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println("read  body err2")
	}
	return data
}

func chooseThreadItem(content string) (items []ThreadItem) {
	tds := threadItemExp.FindAllStringSubmatch(content, 10000)
	var tdstr = make([]string, 0)
	for _, t := range tds {
		var n = strings.Replace(t[0], "\"", "", -1)
		tdstr = append(tdstr, n)
	}
	var threads = make([]ThreadItem, 0)
	for _, t := range tdstr {
		threads = append(threads, ThreadItem{url: "http://boards.4chan.org/s/" + t})
	}
	return threads
}

func chooseImageItem(content string) (items []ThreadItem) {
	tds := imageItemExp.FindAllStringSubmatch(content, 10000)
	var tdstr = make([]string, 0)
	for _, path := range tds {
		path1 := strings.TrimLeft(path[0], "src=\"//")
		path2 := strings.TrimRight(path1, "\"")
		path3 := fmt.Sprintf("http://%s", path2)
		tdstr = append(tdstr, path3)
	}
	var threads = make([]ThreadItem, 0)
	for _, t := range tdstr {
		threads = append(threads, ThreadItem{imgs: t})
	}
	return threads
}
