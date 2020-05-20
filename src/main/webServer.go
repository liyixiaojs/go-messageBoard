package main

import (
	"net/http"
	"text/template"
	"strconv"
	"os"
	"bufio"
	"encoding/json"
)

type SaveFormat struct {
	Id string `json:"id"`
	Time int64 `json:"time"`
	Name string `json:"name"`
	Content string `json:"content"`
}

type Resp struct {
	Success bool `json:"success"`
}

type ResError struct {
	Success bool `json:"success"`
	Msg string `json:"msg"`
}

type List []SaveFormat

type Listp struct {
	Success bool `json:"success"`
	Data List `json:"data"`
}

func (s SaveFormat) toString() string {
	return s.Id + Splitter + strconv.FormatInt(s.Time, 10) + Splitter + s.Name + Splitter + s.Content
}

const Splitter = "<<<&&&>>>>>"

// 写入文件index.txt
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

func webServerCallback(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("./web/index.html")
	if err != nil {
		// debugLog.SetPrefix("[Error]")
		// debugLog.Println("访问错误" + err.Error())
		webLogPrintln("[Error]", "访问错误" + err.Error())
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
		// debugLog.SetPrefix("[Error]")
		// debugLog.Println("访问错误" + err.Error())
		webLogPrintln("[Error]", "访问错误" + err.Error())
		return nil
	}
	return msg
}

func webInterface() {
	http.HandleFunc("/save", saveCallback)
	http.HandleFunc("/queryList", queryList)
	http.HandleFunc("/delete", deleteData)
	http.HandleFunc("/error", errorCallback)
}

func webServer() {
	http.HandleFunc("/", webServerCallback)
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web"))))
	webInterface()
}