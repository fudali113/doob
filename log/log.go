package log

import (
	"fmt"
	"time"
)

const (
	LOG_TEMPLATE = "%s  %s  %s \n"
	ERROR_COLOR  = "\033[31;40m %s \033[0m"
	INFO_COLOR   = "\033[32;40m %s \033[0m"
	NOTICE_COLOR = "\033[37;40m %s \033[0m"
	WARN_COLOR   = "\033[35;40m %s \033[0m"
	DEBUG_COLOR  = "\033[34;40m %s \033[0m"
)

type Log interface {
	Get(string) Log
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Notice(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
}

var (
	looger Log = &SimpleLog{}
)

func GetLog(name string) Log {
	return &SimpleLog{name: name}
}

func SetLog(log Log) {
	looger = log
}

type SimpleLog struct {
	name string
}

func (this *SimpleLog) Get(name string) Log {
	return &SimpleLog{name: name}
}

func (this *SimpleLog) Error(format string, args ...interface{}) {
	template := fmt.Sprintf(format, args...)
	fmt.Printf(fmt.Sprintf(ERROR_COLOR, LOG_TEMPLATE), getTimeString(), getLogInfo(this.name, "ERROR"), template)
}

func (this *SimpleLog) Info(format string, args ...interface{}) {
	template := fmt.Sprintf(format, args...)
	fmt.Printf(fmt.Sprintf(INFO_COLOR, LOG_TEMPLATE), getTimeString(), getLogInfo(this.name, "INFO"), template)
}

func (this *SimpleLog) Notice(format string, args ...interface{}) {
	template := fmt.Sprintf(format, args...)
	fmt.Printf(fmt.Sprintf(NOTICE_COLOR, LOG_TEMPLATE), getTimeString(), getLogInfo(this.name, "NOTICE"), template)
}

func (this *SimpleLog) Warn(format string, args ...interface{}) {
	template := fmt.Sprintf(format, args...)
	fmt.Printf(fmt.Sprintf(WARN_COLOR, LOG_TEMPLATE), getTimeString(), getLogInfo(this.name, "WARN"), template)
}

func (this *SimpleLog) Debug(format string, args ...interface{}) {
	template := fmt.Sprintf(format, args...)
	fmt.Printf(fmt.Sprintf(DEBUG_COLOR, LOG_TEMPLATE), getTimeString(), getLogInfo(this.name, "DEBUG"), template)
}

func getTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func getLogInfo(name, class string) string {
	return " [" + class + "] " + "{name:" + name + "} "
}
