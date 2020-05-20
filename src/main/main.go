package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
)

func initFile() {
	if !Exists("./index.txt") {
		webLogPrintln("[Warning]", "文件 'index.txt' 不存在")
		f,err := os.Create("index.txt")
		defer f.Close()
		if err != nil {
			log.Fatalln("open file error")
		}
	}
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err.(string))
			webLogPrintln("[Error]", err.(string))
		}
	}()
	if Exists("./web") && IsDir("./web") {
		if Exists("./web/index.html") && IsFile("./web/index.html") {
			initFile()
			webServer()
			webLogPrintln("[Info]", "服务器即将开启，访问地址 http://localhost:8080")
			fmt.Println("服务器即将开启，访问地址 http://localhost:8080")
			err := http.ListenAndServe(":8080", nil)
			if err != nil {
				panic("服务器开启错误: " + err.Error())
			}
		} else {
			panic("文件 'web/index.html' 不存在")
		}
	} else {
		panic("文件夹 'web' 不存在")
	}
}