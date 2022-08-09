package flog_test

import (
	"fmt"
	"testing"

	"github.com/xinight/flog"
)

func TestWriteLog(t *testing.T) {
	logger := flog.NewLogger("./log/", flog.OPEN_DEBUG|flog.OPEN_ERROR, 0)
	logger.Start()
	logger.WriteLog(flog.L_DEBUG, "debug %s", "test")
	defer logger.Close()
}

func BenchmarkWriteLog(b *testing.B) {
	logger := flog.NewLogger("./log/", flog.OPEN_DEBUG|flog.OPEN_ERROR, 0)
	logger.Start()
	defer logger.Close()
	for i := 0; i < b.N; i++ {
		logger.WriteLog(flog.L_DEBUG, "debug %s%d", "test", i)
	}
}

func ExampleNewLogger() {
	logger := flog.NewLogger("./log/", flog.OPEN_DEBUG|flog.OPEN_ERROR, 0)
	fmt.Printf("%T\n", logger)

	// Output:
	// *flog.FLogger
}

func ExampleFLogger_Start() {
	logger := flog.NewLogger("./log/", flog.OPEN_DEBUG|flog.OPEN_ERROR, 0)
	defer logger.Close()
	logger.Start()

	// Output:
}
func ExampleFLogger_WriteLog() {
	logger := flog.NewLogger("./log/", flog.OPEN_DEBUG|flog.OPEN_ERROR, 0)
	defer logger.Close()
	err := logger.WriteLog(flog.L_DEBUG, "debug %s", "test1")
	fmt.Println(err.Error())
	logger.Start()
	err = logger.WriteLog(flog.L_DEBUG, "debug %s", "test1")
	fmt.Println(err)

	// Output:
	// logger do not start
	// <nil>
}
