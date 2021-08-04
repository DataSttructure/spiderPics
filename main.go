package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	var keyWord string
	fmt.Println("Please input a keyword to download")
	fmt.Scanln(&keyWord)
	var searchURL string
	for i := 1; i < 10; i++ {
		searchURL = "https://www.bizhizu.cn/search/" + keyWord + "/" + strconv.Itoa(int(i)) + ".html"
		// fmt.Println(searchURL)
		mainUrl := searchImageURL(searchURL, keyWord)
		for i := range mainUrl {
			saveToFolder(mainUrl[i], "/home/liting/go/src/pics/")
		}
	}
}

func searchImageURL(url, keyWord string) (urlSlice []string) {

	// searchURL := url + "/search/" + keyWord
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		href = strings.TrimSpace(href)
		if len(href) > 2 {
			if strings.HasPrefix(href, "https://www.bizhizu.cn/pic") || strings.HasPrefix(href, "https://www.bizhizu.cn/bizhi") {
				if IsUrl(href) {
					dup := false
					for idx := range urlSlice {
						if href == urlSlice[idx] {
							dup = true
						}
					}
					if !dup {
						urlSlice = append(urlSlice, href)
						// fmt.Println("修改之后的url：", href)
					}
				}
			}
		}
	})
	return urlSlice
}

func getHtml(url string) (urlSlice []string) {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		href = strings.TrimSpace(href)
		if len(href) > 2 {

			if strings.HasPrefix(href, "/ziran/") {
				href = SamePathUrl("https://www.bizhizu.cn/", href, 1)
				if IsUrl(href) {
					dup := false
					for idx := range urlSlice {
						if href == urlSlice[idx] {
							dup = true
						}
					}
					if !dup {
						urlSlice = append(urlSlice, href)
						// fmt.Println("修改之后的url：", href)
					}
				}
			}
		}
	})
	return urlSlice
}

func getSubURL(url string) (urlSlice []string) {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		href = strings.TrimSpace(href)
		if len(href) > 2 && strings.HasPrefix(href, "https://www.bizhizu.cn/pic/") {
			if IsUrl(href) {
				dup := false
				for idx := range urlSlice {
					if href == urlSlice[idx] {
						dup = true
					}
				}
				if !dup {
					urlSlice = append(urlSlice, href)
					// fmt.Println("修改之后的url：", href)
				}
			}
		}
	})

	// doc.Find("a").Each(func(i int, selection *goquery.Selection) {
	// 	title = selection.Text()
	// 	// fmt.Println(title)
	// })

	return urlSlice
}

func IsUrl(str string) bool {
	if strings.HasPrefix(str, "#") || strings.HasPrefix(str, "//") || strings.HasSuffix(str, ".exe") || strings.HasSuffix(str, ":void(0);") {
		return false
	} else if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
		return false
	} else if strings.EqualFold(str, "javascript:;") {
		return false
	} else {
		return true
	}
	return true
}

func SamePathUrl(preUrl string, url string, mark int) (newUrl string) {
	last := strings.LastIndex(preUrl, "/")
	if last == 6 {
		newUrl = preUrl[:last] + url
	} else {
		if mark == 1 {
			newUrl = preUrl[:last] + url
		} else {
			newPreUrl := preUrl[:last]
			newLast := strings.LastIndex(newPreUrl, "/")
			newUrl = newPreUrl[:newLast] + url
		}
	}
	return newUrl
}

func saveToFolder(url, path string) {
	// url := "https://www.bizhizu.cn/pic/35413.html"
	title, imgNum := getTitleAndImgNum(url)

	dir := path + title + strconv.Itoa(imgNum) + "/"
	isExist, err := PathExists(dir)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if isExist {
			fmt.Println(dir + "文件夹已存在！")
		} else {
			// 文件夹名称，权限
			err := os.Mkdir(dir, os.ModePerm)
			if err != nil {
				fmt.Println(dir+"文件夹创建失败：", err.Error())
			} else {
				fmt.Println(dir + "文件夹创建成功！")
			}
		}
	}

	for i := 0; i < imgNum; i++ {
		// url1 := "https://www.bizhizu.cn/pic/35413-" + strconv.Itoa(i) + ".html"
		url1 := url[:len(url)-5] + "-" + strconv.Itoa(i) + ".html"
		imageUrl, err := getImageURL(url1)
		if err != nil {
			fmt.Println(err)
			return
		}
		saveImage(imageUrl, path, title, i, imgNum)
	}
}

func getTitleAndImgNum(url string) (string, int) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	var title string
	doc.Find("h1").Each(func(i int, selection *goquery.Selection) {
		title = selection.Text()
	})
	var idx, idxNum1, idxNum2 int
	for i := range title {
		if title[i] == '(' {
			idx = i
		}
		if title[i] == '/' {
			idxNum1 = i
		}
		if title[i] == ')' {
			idxNum2 = i
		}
	}
	imgNum, err := strconv.Atoi(title[idxNum1+1 : idxNum2])
	if err != nil {
		fmt.Println(err.Error())
	}
	title = title[0:idx]
	return title, imgNum
}

func getImageURL(url string) (string, error) {
	var imageUrl string
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".pic > a > img").Each(func(i int, selection *goquery.Selection) {
		imageUrl, _ = selection.Attr("src")
	})
	return imageUrl, nil
}

func saveImage(url, path, title string, i, imgNum int) {
	path = path + title + strconv.Itoa(imgNum) + "/" + strconv.Itoa(i) + ".jpg"
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http err:", err)
		return
	}

	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		//写入文件
		f.Write(buf[:n])
	}
	fmt.Println(path + " saved!")
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
