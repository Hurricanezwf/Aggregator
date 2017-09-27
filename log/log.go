// Package log is for debuging rdb only.
// You can set _ENABLE_LOG=false to disable it.
package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

const (
	LOG_LEVEL_ERROR = 1
	LOG_LEVEL_WARN  = 2
	LOG_LEVEL_INFO  = 3
	LOG_LEVEL_DEBUG = 4
)

var (
	ENABLE_LOG bool  = true
	_LOG_LEVEL int64 = LOG_LEVEL_DEBUG
)

func init() {
	if ENABLE_LOG {
	}
}

func Debug(format string, v ...interface{}) {
	if ENABLE_LOG && _LOG_LEVEL >= LOG_LEVEL_DEBUG {
		nowTime := time.Now().Format("2006/01/02 15:04:05")
		fileName, lineNo := callerInfo()
		fmt.Println(fmt.Sprintf("%s [%s:%d] [D] %s", nowTime, fileName, lineNo, fmt.Sprintf(format, v...)))
	}
}

func Info(format string, v ...interface{}) {
	if ENABLE_LOG && _LOG_LEVEL >= LOG_LEVEL_INFO {
		nowTime := time.Now().Format("2006/01/02 15:04:05")
		fileName, lineNo := callerInfo()
		fmt.Println(fmt.Sprintf("%s [%s:%d] [I] %s", nowTime, fileName, lineNo, fmt.Sprintf(format, v...)))
	}
}

func Warn(format string, v ...interface{}) {
	if ENABLE_LOG && _LOG_LEVEL >= LOG_LEVEL_WARN {
		nowTime := time.Now().Format("2006/01/02 15:04:05")
		fileName, lineNo := callerInfo()
		fmt.Println(fmt.Sprintf("%s [%s:%d] [W] %s", nowTime, fileName, lineNo, fmt.Sprintf(format, v...)))
	}
}

func Error(format string, v ...interface{}) {
	if ENABLE_LOG && _LOG_LEVEL >= LOG_LEVEL_ERROR {
		nowTime := time.Now().Format("2006/01/02 15:04:05")
		fileName, lineNo := callerInfo()
		fmt.Println(fmt.Sprintf("%s [%s:%d] [E] %s", nowTime, fileName, lineNo, fmt.Sprintf(format, v...)))
	}
}

func callerInfo() (fileName string, lineNo int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return fileName, lineNo
	}

	return filepath.Base(file), line
}
