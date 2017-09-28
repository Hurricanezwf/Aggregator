package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	LOG_LEVEL_ERROR = 1
	LOG_LEVEL_WARN  = 2
	LOG_LEVEL_INFO  = 3
	LOG_LEVEL_DEBUG = 4
	LOG_LEVEL_TRACE = 5
)

type Logger struct {
	level int

	w io.Writer
}

func New() *Logger {
	return &Logger{
		level: LOG_LEVEL_DEBUG,
		w:     os.Stdout,
	}
}

func (l Logger) Trace(format string, v ...interface{}) {
	if l.level >= LOG_LEVEL_TRACE {
		nowTime := time.Now().Format("2006/01/02 15:04:05")
		fileName, lineNo := l.callerInfo()
		fmt.Println(fmt.Sprintf("%s [%s:%d] [T] %s", nowTime, fileName, lineNo, fmt.Sprintf(format, v...)))
	}
}

func (l Logger) Debug(format string, v ...interface{}) {
	if l.level >= LOG_LEVEL_DEBUG {
		nowTime := time.Now().Format("2006/01/02 15:04:05")
		fileName, lineNo := l.callerInfo()
		fmt.Println(fmt.Sprintf("%s [%s:%d] [D] %s", nowTime, fileName, lineNo, fmt.Sprintf(format, v...)))
	}
}

func (l Logger) Info(format string, v ...interface{}) {
	if l.level >= LOG_LEVEL_INFO {
		nowTime := time.Now().Format("2006/01/02 15:04:05")
		fileName, lineNo := l.callerInfo()
		fmt.Println(fmt.Sprintf("%s [%s:%d] [I] %s", nowTime, fileName, lineNo, fmt.Sprintf(format, v...)))
	}
}

func (l Logger) Warn(format string, v ...interface{}) {
	if l.level >= LOG_LEVEL_WARN {
		nowTime := time.Now().Format("2006/01/02 15:04:05")
		fileName, lineNo := l.callerInfo()
		fmt.Println(fmt.Sprintf("%s [%s:%d] [W] %s", nowTime, fileName, lineNo, fmt.Sprintf(format, v...)))
	}
}

func (l Logger) Error(format string, v ...interface{}) {
	if l.level >= LOG_LEVEL_ERROR {
		nowTime := time.Now().Format("2006/01/02 15:04:05")
		fileName, lineNo := l.callerInfo()
		fmt.Println(fmt.Sprintf("%s [%s:%d] [E] %s", nowTime, fileName, lineNo, fmt.Sprintf(format, v...)))
	}
}

func (l *Logger) SetLevel(level int) bool {
	if level < LOG_LEVEL_ERROR || level > LOG_LEVEL_TRACE {
		return false
	}
	l.level = level
	return true
}

func (l *Logger) SetWriter(w io.Writer) {
	l.w = w
}

func (l Logger) callerInfo() (fileName string, lineNo int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return fileName, lineNo
	}

	return filepath.Base(file), line
}
