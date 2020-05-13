/*
* @Author: Jeffery
* @Date:   2020-04-21 18:29:00
* @Last Modified by:   Jeffery
* @Last Modified time: 2020-05-13 11:03:02
 */
package main

import (
	"flag"
	"fmt"
	"github.com/ty4z2008/image-downloader/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	version   string = "0.1.0"
	uri       string
	directory string
	cookie    string
	help      bool
)

func main() {
	utils.Init("image-downloader", true)
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
	resp, err := httpClient(rawurl)
	if err != nil {
		utils.Error(err)
		return nil
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
	reg := regexp.MustCompile(`<img.*src="(([a-zA-Z0-9]|\/|(https?)).*?)".*\/?>`)
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
		} else if !strings.HasPrefix(imgUrl, "http") {
			imgUrl = fmt.Sprintf("%s://%s%s%s", scheme, host, "/", imgUrl)
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
	resp, err := httpClient(imgUrl)
	if err != nil {
		utils.Error(err)
		return
	}
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
func httpClient(rawurl string) (res *http.Response, err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		utils.Fatal(err)
		return nil, err
	}

	client := http.Client{
		Jar: jar,
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		utils.Error("Make a request fail")
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		utils.Error("request fatal:")
		utils.Error(err)
		return nil, err
	}
	return resp, nil
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
	resp, err := httpClient(rawurl)
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

	flag.BoolVar(&help, "h", false, "Show usage")
	flag.StringVar(&directory, "d", "images", "Save the directory of image")
	flag.StringVar(&cookie, "cookie", "", "Send cookies from string")
	flag.StringVar(&uri, "uri", "", "Your want download website url for image ")
	flag.Usage = usage
	if help {
		flag.Usage()
	}

}

// Print the tools usage
func usage() {
	fmt.Fprintf(os.Stderr, `download image tools for website version %s
Usage: download [-h] [-uri https://github.com] [-d] images

Options:
`, version)
	flag.PrintDefaults()
}
