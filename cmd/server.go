/*
* @Author: Jeffery
* @Date:   2020-04-21 18:29:00
* @Last Modified by:   Jeffery
* @Last Modified time: 2020-05-09 15:17:42
 */
package main

import (
	"flag"
	"fmt"
	"github.com/ty4z2008/image-downloader/utils"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

var (
	version   string = "0.0.1"
	uri       string
	directory string
	help      bool
)

func main() {
	logger.Init("image-downloader", true)
	logger.Info("start the world")
	parseFlag()
	//validate url
	isvalid := validateUrl(uri)
	if !isvalid {
		logger.Errorf("your url %s is not support", uri)
		return
	}
	//crawler page
	crawlerPage(uri)
}

//crawler page
func crawlerPage(rawurl string) {
	logger.Infof("start crawler %s page", rawurl)
	resp, err := http.Get(rawurl)
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()

	var bodyStr string
	if resp.StatusCode == 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Fatal(err)
		}
		bodyStr = string(bodyBytes)
	}
	reg := regexp.MustCompile(` <img\s[^>]*?src\s*=\s*['\"]([^'\"]*?)['\"][^>]*?>`)
	matchImgUrl := reg.FindAllStringSubmatch(bodyStr, -1)
	for _, item := range matchImgUrl {
		fmt.Println(item[1])
	}
}

//validate crawler target is alive
//the method has timeout
func validateUrl(rawurl string) bool {
	if rawurl == "" {
		logger.Error("Your url is empty")
		return false
	}

	logger.Infof("dial %s", rawurl)
	resp, err := http.Get(rawurl)
	if err != nil {
		logger.Error(err)
		return false
	}
	if resp.StatusCode != 200 {
		logger.Infof("%s response code is %d.", rawurl, resp.StatusCode)
		return false
	}
	return true
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
