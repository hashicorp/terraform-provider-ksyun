package logger

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

const skip = 1
const ReqFormat = "ACTION:%s;REQ:%+v"
const RespFormat = "ACTION:%s;REQ:%+v;RESP:%+v"
const ErrFormat = "ACTION:%s;REQ:%+v;ERR:%+v"
const AllFormat = "ACTION:%s;REQ:%+v;RESP:%+v;ERR:%+v"

func Debug(format string, action string, req interface{}, v ...interface{}) {
	_, file, line, _ := runtime.Caller(skip)
	start := strings.LastIndex(file, "/")
	if start != -1 {
		file = file[start+1:]
	}
	message := fmt.Sprintf("[DEBUG] {%v:%v}", file, line)
	log.Printf(message+format, action, req, v)
}
func DebugInfo(format string, info interface{}) {
	_, file, line, _ := runtime.Caller(skip)
	start := strings.LastIndex(file, "/")
	if start != -1 {
		file = file[start+1:]
	}
	message := fmt.Sprintf("[DEBUG] {%v:%v}", file, line)
	log.Printf(message+format, info)
}

func Info(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(skip)
	start := strings.LastIndex(file, "/")
	if start != -1 {
		file = file[start+1:]
	}
	message := fmt.Sprintf("[INFO] {%v:%v}", file, line)
	log.Printf(message+format, v)
}
