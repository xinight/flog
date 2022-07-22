package flog_test

import (
	"flog"
	"testing"
)

func TestLog(t *testing.T) {
	logger := flog.NewLogger("./log/", flog.OPEN_DEBUG|flog.OPEN_ERROR, 0)
	logger.Start()
	logger.WriteLog(flog.L_DEBUG, "debug %s", "test")
	defer logger.Close()
}

func BenchmarkLog(b *testing.B) {
	logger := flog.NewLogger("./log/", flog.OPEN_DEBUG|flog.OPEN_ERROR, 0)
	logger.Start()
	defer logger.Close()
	for i := 0; i < b.N; i++ {
		logger.WriteLog(flog.L_DEBUG, "debug %s%d", "test", i)
	}
}
