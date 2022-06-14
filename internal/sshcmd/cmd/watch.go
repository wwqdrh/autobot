package cmd

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/wwqdrh/logger"
)

const oneMBByte = 1024 * 1024

//WatchFileSize is
func LoggerFileSize(filename string, size int) {
	t := time.NewTicker(3 * time.Second) //every 3s check file
	defer t.Stop()
	for {
		if <-t.C; true {
			length := String("/bin/sh", "-c", "ls -l "+filename+" | awk '{print $5}'", "")
			length = strings.Replace(length, "\n", "", -1)
			length = strings.Replace(length, "\r", "", -1)
			lengthByte, _ := strconv.Atoi(length)
			if lengthByte == size {
				t.Stop()
			}
			lengthFloat := float64(lengthByte)
			value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", lengthFloat/oneMBByte), 64)
			logger.DefaultLogger.Info(fmt.Sprintf("[os][%s]transfer total size is: %.2f%s", filename, value, "MB"))
		}
	}
}

//RemoteFileExist is
func IsFileExist(filepath string) bool {
	// if remote file is
	// ls -l | grep aa | wc -l
	fileName := path.Base(filepath) // aa
	fileDirName := path.Dir(filepath)
	fileCommand := fmt.Sprintf("ls -l %s | grep %s | wc -l", fileDirName, fileName)
	data := strings.Replace(String("/bin/sh", "-c", fileCommand), "\r", "", -1)
	data = strings.Replace(data, "\n", "", -1)
	count, err := strconv.Atoi(strings.TrimSpace(data))
	defer func() {
		if r := recover(); r != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("[os][%s]RemoteFileExist:%s", filepath, err.Error()))
		}
	}()
	if err != nil {
		panic(1)
	}
	return count != 0
}
