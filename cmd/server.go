/*
* @Author: Jeffery
* @Date:   2020-04-21 18:29:00
* @Last Modified by:   Jeffery
* @Last Modified time: 2020-05-11 10:39:42
 */
package main

import (
	"flag"
	"fmt"
	"github.com/ty4z2008/image-downloader/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	version   string = "0.0.1"
	uri       string
	directory string
	help      bool
)

func main() {
	utils.Init("image-downloader", true)
	utils.Info("start the world")
	parseFlag()
	//validate url
	isvalid := validateUrl(uri)
	if !isvalid {
		utils.Errorf("your url %s is not support", uri)
		return
	}
	//create directory
	createDownloadDir(directory)
	//crawler page
	imgUrls := crawlerPageImage(uri)
	if len(imgUrls) == 0 {
		utils.Warning("empty images resource")
		return
	}
	utils.Infof("find %d images", len(imgUrls))
	var wg sync.WaitGroup
	for _, imgUrl := range imgUrls {
		wg.Add(1)
		go download(imgUrl, &wg)
	}
	wg.Wait()
	utils.Info("download image task complete.")
}

//crawler page
func crawlerPageImage(rawurl string) []string {
	var urls []string
	utils.Infof("start crawler %s page", rawurl)
	resp, err := http.Get(rawurl)
	if err != nil {
		utils.Error(err)
		return urls
	}
	defer resp.Body.Close()

	var bodyStr string

	if resp.StatusCode == 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			utils.Fatal(err)
		}
		bodyStr = string(bodyBytes)
	}
	reg := regexp.MustCompile(`<img.*src="((\/|(https?)).*?)".*\/?>`)
	matchImgUrl := reg.FindAllStringSubmatch(bodyStr, -1)

	u, _ := url.Parse(rawurl)
	scheme := u.Scheme
	host := u.Host
	//use map extract duplicate
	var urlsMap = make(map[string]string, len(matchImgUrl))
	for _, item := range matchImgUrl {
		imgUrl := item[1]
		if strings.HasPrefix(imgUrl, "//") {
			imgUrl = fmt.Sprintf("%s:%s", scheme, imgUrl)
		} else if strings.HasPrefix(imgUrl, "/") {
			imgUrl = fmt.Sprintf("%s://%s%s", scheme, host, imgUrl)
		}
		urlsMap[imgUrl] = ""
	}

	for imgUrl := range urlsMap {
		urls = append(urls, imgUrl)
	}
	return urls
}

//download image
func download(imgUrl string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName, err := getUrlFileName(imgUrl)
	if err != nil {
		return
	}
	path := filepath.Join(directory, fileName)
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		utils.Warning(err)
		return
	}
	client := httpClient()
	if client == nil {
		utils.Warning("get http client fail")
		return
	}
	resp, err := client.Get(imgUrl)
	defer resp.Body.Close()
	if err != nil {
		utils.Warning(err)
		return
	}
	size, err := io.Copy(file, resp.Body)
	if err != nil {
		utils.Warning(err)
		return
	}
	utils.Infof("Downloaded a file %s with size %s ", fileName, utils.ByteCount(size, "KB"))
}

//http client
func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

//get file name by image url
func getUrlFileName(imgUrl string) (string, error) {
	fileUrl, err := url.Parse(imgUrl)
	if err != nil {
		utils.Warning(err)
		return "", err
	}
	path := fileUrl.Path
	segments := strings.Split(path, "/")
	return segments[len(segments)-1], nil
}

//validate crawler target is alive
//the method has timeout
func validateUrl(rawurl string) bool {
	if rawurl == "" {
		utils.Error("Your url is empty")
		os.Exit(0)
		return false
	}

	utils.Infof("dial %s", rawurl)
	resp, err := http.Get(rawurl)
	if err != nil {
		utils.Error(err)
		return false
	}
	if resp.StatusCode != 200 {
		utils.Infof("%s response code is %d.", rawurl, resp.StatusCode)
		return false
	}
	return true
}

//create save image directory
func createDownloadDir(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.FileMode(0755))
	}
	return path
}

//parse flag
func parseFlag() {
	defer flag.Parse()

	flag.BoolVar(&help, "h", false, "show usage")
	flag.StringVar(&directory, "d", "images", "save the directory of image")
	flag.StringVar(&uri, "uri", "", "Your want download website url for image ")
	flag.Usage = usage
	if help {
		flag.Usage()
	}

}

// Print the tools usage
func usage() {
	fmt.Fprintf(os.Stderr, `download image tools for website version %s
Usage: download [-h] [-uri https://github.com] 

Options:
`, version)
	flag.PrintDefaults()
}
