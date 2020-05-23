package main

import (
	// "fmt"
	"log"
	"os"
)

var debugLog *log.Logger

func webLogPrintln(prefix string, msg string) {
	debugLog.SetPrefix(prefix)
	debugLog.Println(msg)
}

func init() {
	fileName := "index.log"
	logFile,err  := os.Create(fileName)
	// defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error")
	}
	debugLog = log.New(logFile,"[Info]",log.Llongfile)
	debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)
}
