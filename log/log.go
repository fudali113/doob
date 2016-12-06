package log

import (
	"fmt"
	"time"
)

const (
	LOG_TEMPLATE = "%s  %s  %s \n"
	ERROR_COLOR  = "\033[1;31;40m %s \033[0m"
	INFO_COLOR   = "\033[1;32;40m %s \034"
	NOTICE_COLOR = "\033[1;37;40m %s \034"
	WARN_COLOR   = "\033[1;35;40m %s \034"
	DEBUG_COLOR  = "\033[1;34;40m %s \034"
)

type Log interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Notice(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
}

func GetLog(name string) Log {
	return &SimpleLog{name: name}
}

type SimpleLog struct {
	name string
}

func (this *SimpleLog) Error(format string, args ...interface{}) {
	template := fmt.Sprintf(format, args...)
	fmt.Printf(fmt.Sprintf(ERROR_COLOR, LOG_TEMPLATE), getTimeString(), this.name+" [ERROR]", template)
}

func (this *SimpleLog) Info(format string, args ...interface{}) {
	template := fmt.Sprintf(format, args...)
	fmt.Printf(fmt.Sprintf(INFO_COLOR, LOG_TEMPLATE), getTimeString(), this.name+" [INFO]", template)
}

func (this *SimpleLog) Notice(format string, args ...interface{}) {
	template := fmt.Sprintf(format, args...)
	fmt.Printf(fmt.Sprintf(NOTICE_COLOR, LOG_TEMPLATE), getTimeString(), this.name+" [NOTICE]", template)
}

func (this *SimpleLog) Warn(format string, args ...interface{}) {
	template := fmt.Sprintf(format, args...)
	fmt.Printf(fmt.Sprintf(WARN_COLOR, LOG_TEMPLATE), getTimeString(), this.name+" [WARN]", template)
}

func (this *SimpleLog) Debug(format string, args ...interface{}) {
	template := fmt.Sprintf(format, args...)
	fmt.Printf(fmt.Sprintf(DEBUG_COLOR, LOG_TEMPLATE), getTimeString(), this.name+" [DEBUG]", template)
}

func getTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
