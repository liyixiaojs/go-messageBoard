package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/satori/go.uuid"
	"time"
	"os"
	"strconv"
	"strings"
	"errors"
)

// 保存
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
		webLogPrintln("[Error]", "访问错误" + openFileErr.Error())
			// debugLog.SetPrefix("[Error]")
			// debugLog.Println("访问错误" + openFileErr.Error())
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

// 查询列表
func queryList(w http.ResponseWriter, req *http.Request) {
	f, err := os.OpenFile("./index.txt", os.O_RDONLY, 0666)
	if err != nil {
		// debugLog.SetPrefix("[Error]")
		// debugLog.Println("查询列表失败" + err.Error())
		webLogPrintln("[Error]", "查询列表失败" + err.Error())
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

// 删除数据
func deleteData(w http.ResponseWriter, req *http.Request) {
	var user map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &user)
	f, err := os.OpenFile("./index.txt", os.O_RDONLY, 0666)
	if err != nil {
		// debugLog.SetPrefix("[Error]")
		// debugLog.Println("删除留言失败" + err.Error())
		webLogPrintln("[Error]", "删除留言失败" + err.Error())
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

func errorCallback(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(getErrorMsg(errors.New("错误信息")))
}