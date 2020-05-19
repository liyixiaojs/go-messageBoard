
package main

import (
	// "io"
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"text/template"
	"encoding/json"
	// "bytes"
	"github.com/satori/go.uuid"
	"time"
	"os"
	"strconv"
	"strings"
	"bufio"
	"errors"
	// "unicode/utf8"
	// "unsafe"
)

type SaveFormat struct {
	Id string `json:"id"`
	Time int64 `json:"time"`
	Name string `json:"name"`
	Content string `json:"content"`
}

type List []SaveFormat

type Listp struct {
	Success bool `json:"success"`
	Data List `json:"data"`
}

type Resp struct {
	Success bool `json:"success"`
}

type ResError struct {
	Success bool `json:"success"`
	Msg string `json:"msg"`
}

const Splitter = "<<<&&&>>>>>"

func (s SaveFormat) toString() string {
	return s.Id + Splitter + strconv.FormatInt(s.Time, 10) + Splitter + s.Name + Splitter + s.Content
}

func webServer(w http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("./web/index.html")
		if err != nil {
			debugLog.SetPrefix("[Error]")
			debugLog.Println("访问错误" + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := map[string]string{
		}

		t.Execute(w, data)
}

func getErrorMsg(err error) []byte {
	profile := ResError{
		Success: false,
		Msg: err.Error(),
	}
	msg, err := json.Marshal(profile)
	if err != nil {
		debugLog.SetPrefix("[Error]")
		debugLog.Println("访问错误" + err.Error())
		return nil
	}
	return msg
}

func errorCallback(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(getErrorMsg(errors.New("错误信息")))
}

func saveCallback(w http.ResponseWriter, req *http.Request) {
	var user map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &user)
	u1 := uuid.Must(uuid.NewV4())
	param := SaveFormat{
		Id: u1.String(),
		Time: time.Now().UnixNano() / 1000000,
		Name: user["name"].(string),
		Content: user["content"].(string),
	}

	f, openFileErr := os.OpenFile("./index.txt", os.O_RDWR|os.O_APPEND, 0666)
	if openFileErr != nil {
			debugLog.SetPrefix("[Error]")
			debugLog.Println("访问错误" + openFileErr.Error())
			http.Error(w, openFileErr.Error(), http.StatusInternalServerError)
			return
	}

	defer f.Close()

	profile := Resp{
		Success: true,
	}
	if len(param.Name) > 0 && len(param.Content) > 0 {
		data := param.toString()
		f.Write([]byte(data + "\n"))
	} else {
		profile.Success = false
	}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func queryList(w http.ResponseWriter, req *http.Request) {
	f, err := os.OpenFile("./index.txt", os.O_RDONLY, 0666)
	if err != nil {
		debugLog.SetPrefix("[Error]")
		debugLog.Println("查询列表失败" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer f.Close()

	temp := make([]byte, 1024 * 4)
	fileLen, _ := f.Read(temp)
	data := string(temp[:fileLen])
	var list List
	kv := strings.Split(data, "\n")
	for i := 0; i < len(kv); i++ {
		items := strings.Split(kv[i], Splitter)
		if len(items) >= 4 {
			time, _ := strconv.ParseInt(items[1], 10, 64)
			list = append(list, SaveFormat{
				Id: items[0],
				Time: time,
				Name: items[2],
				Content: items[3],
			})
		}
	}

	resData := Listp{
		Success: true,
		Data: list,
	}

	jsonData, err := json.Marshal(resData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func writeToFile(outPut []byte) error {
	f, err := os.OpenFile("index.txt", os.O_WRONLY|os.O_TRUNC, 0666)
	defer f.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	_, err = writer.Write(outPut)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func deleteData(w http.ResponseWriter, req *http.Request) {
	var user map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &user)
	f, err := os.OpenFile("./index.txt", os.O_RDONLY, 0666)
	if err != nil {
		debugLog.SetPrefix("[Error]")
		debugLog.Println("删除留言失败" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	id := user["id"]
	name := user["name"]
	temp := make([]byte, 1024 * 4)
	fileLen, _ := f.Read(temp)
	data := string(temp[:fileLen])
	kv := strings.Split(data, "\n")
	var index int
	for i := 0; i < len(kv); i++ {
		items := strings.Split(kv[i], Splitter)
		if len(items) >= 4 && id == items[0] && name == items[2] {
			index = i
			break
		}
	}
	if index > -1 {
		kv = append(kv[:index], kv[index+1:]...)
	}
	newData := strings.Join(kv, "\n")
	writeToFile([]byte(newData))
	profile := Resp{
		Success: true,
	}
	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func webInterface() {
	http.HandleFunc("/save", saveCallback)
	http.HandleFunc("/queryList", queryList)
	http.HandleFunc("/delete", deleteData)
	http.HandleFunc("/error", errorCallback)
}

func initFile() {
	if !Exists("./index.txt") {
		debugLog.SetPrefix("[Warning]")
		debugLog.Println("文件 'index.txt' 不存在")
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

var debugLog *log.Logger

func main() {
	fileName := "index.log"
	logFile,err  := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
			log.Fatalln("open file error")
	}
	debugLog = log.New(logFile,"[Info]",log.Llongfile)
	debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)
	defer func() {
		if err := recover(); err != nil {
			debugLog.SetPrefix("[Error]")
			debugLog.Println(errors.New(err.(string)))
		}
	}()
	if Exists("./web") && IsDir("./web") {
		if Exists("./web/index.html") && IsFile("./web/index.html") {
			initFile()
			http.HandleFunc("/", webServer)
			http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web"))))
			debugLog.SetPrefix("[Info]")
			debugLog.Println("服务器即将开启，访问地址 http://localhost:8080")
			fmt.Println("服务器即将开启，访问地址 http://localhost:8080")
		
			webInterface()
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