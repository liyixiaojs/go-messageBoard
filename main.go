package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"io/ioutil"
	"strconv"
	"flag"
	yaml "gopkg.in/yaml.v2"
)

type WebConfig struct {
	WebPath string `yaml:"webPath"`
	Port int `yaml:"port"`
	SaveFile string `yaml:"saveFile"`
}

type Config struct{
	WebConfig WebConfig `yaml:"webServe"`
}

var Conf WebConfig

func initFile() {
	if !Exists(Conf.SaveFile) {
		webLogPrintln("[Warning]", "文件 '" + Conf.SaveFile + "' 不存在")
		f,err := os.Create(Conf.SaveFile)
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


func ReadYamlConfig(path string)  (*Config, error){
	conf := &Config{}
	if f, err := ioutil.ReadFile(path); err != nil {
	    return nil,err
	} else {
		yaml.Unmarshal(f, &conf)
		// yaml.NewDecoder(f).Decode(conf)
	}
	return  conf,nil
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err.(string))
			webLogPrintln("[Error]", err.(string))
		}
	}()

	conf,confErr := ReadYamlConfig("./conf.yaml")

	if confErr != nil {
		panic("解析配置文件失败")
	}

	Conf = conf.WebConfig

	var port int

	flag.IntVar(&port, "port", Conf.Port, "端口、默认" + strconv.Itoa(Conf.Port))

	var GO_PORT string
	GO_PORT = os.Getenv("GO_PORT")
	if len(GO_PORT) > 0 {
		p, pErr := strconv.Atoi(GO_PORT)
		if pErr != nil {
			panic("端口类型错误")
		}
		port = p
	}

	flag.Parse()

	if Exists(Conf.WebPath) && IsDir(Conf.WebPath) {
		if Exists(Conf.WebPath + "/index.html") && IsFile(Conf.WebPath + "/index.html") {
			initFile()
			webServer()
			webLogPrintln("[Info]", "服务器即将开启，访问地址 http://localhost:" + strconv.Itoa(port))
			fmt.Println("服务器即将开启，访问地址 http://localhost:" + strconv.Itoa(port))
			err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
			if err != nil {
				panic("服务器开启错误: " + err.Error())
			}
		} else {
			panic("文件 '" + Conf.WebPath + "/index.html" + "' 不存在")
		}
	} else {
		panic("文件夹 '" + Conf.WebPath + "' 不存在")
	}
}