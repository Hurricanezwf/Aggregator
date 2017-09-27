package log_test

import (
	"testing"
	"unialarm/rdb/util/log"
)

func TestLog(t *testing.T) {
	log.Debug("Hello, This is debug log")
	log.Info("Hello, This is info log")
	log.Warn("Hello, This is warn log")
	log.Error("Hello, This is error log")
}
