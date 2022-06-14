package filesize

import (
	"crypto/tls"
	"net/http"

	"github.com/wwqdrh/logger"
)

//Do is fetch file size
func Do(url string) int64 {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	defer func() {
		if r := recover(); r != nil {
			logger.DefaultLogger.Error("[globals] get file size is error： " + r.(error).Error())
		}
	}()
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	return resp.ContentLength
}
