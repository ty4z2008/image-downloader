A page image downloader

This project inspired from [https://github.com/kkdai/project52](https://github.com/kkdai/project52)

Logic

- parse the crawle url and download directory from command args
- get the body of response 
- use regex match image url with suffix
- use go concurrency download the image to your directory
- exit