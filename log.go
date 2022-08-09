package flog

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type logLv int

const (
	L_DEBUG logLv = iota
	L_Info
	L_WARING
	L_ERROR
	L_FATAL
	L_PANIC
)

const (
	DEFAULTSIZE  = 100
	DEFAULTSLEEP = 500 * time.Millisecond
)

const (
	OPEN_DEBUG  int = syscall.O_RDONLY
	OPEN_INFO   int = syscall.O_WRONLY
	OPEN_WARING int = syscall.O_RDWR
	OPEN_ERROR  int = syscall.O_APPEND
	OPEN_FATAL  int = syscall.O_CREAT
	OPEN_PANIC  int = syscall.O_EXCL
)

var level = map[logLv]string{
	L_DEBUG:  "debug",
	L_Info:   "info",
	L_WARING: "waring",
	L_ERROR:  "error",
	L_FATAL:  "fatal",
	L_PANIC:  "panic",
}

var ch chan *logC

type logC struct {
	File  string
	Level logLv
	Time  string
	Line  int
	Msg   string
}

var files = map[logLv]*os.File{
	L_DEBUG:  nil,
	L_Info:   nil,
	L_WARING: nil,
	L_ERROR:  nil,
	L_FATAL:  nil,
	L_PANIC:  nil,
}

type FLogger struct {
	// log root path
	rootPath string
	// goroutine number
	num int
	max int
	min int

	flag int
	// channel size
	size int

	// logger status
	start bool
	// goroutine wait group
	wg *sync.WaitGroup
}

var Logger *FLogger

// new logger
func NewLogger(path string, flag int, chanSize int) *FLogger {
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	if chanSize <= 0 {
		chanSize = DEFAULTSIZE
	}
	Logger = &FLogger{
		rootPath: path,
		num:      1,
		max:      1,
		min:      1,
		flag:     flag,
		size:     chanSize,
		start:    false,
		wg:       &sync.WaitGroup{},
	}
	return Logger
}

// start goroutine to write log
func (l *FLogger) Start() {
	l.start = true
	l.wg.Add(l.num)
	ch = make(chan *logC, l.size)
	for i := 0; i < l.num; i++ {
		go func(c chan *logC) {
			for {
				select {
				case log, ok := <-c:
					if !ok {
						l.wg.Done()
						return
					}
					log.write(l.rootPath)
				default:
					time.Sleep(DEFAULTSLEEP)
				}
			}
		}(ch)
	}
}

// send log struct to channel
func (l *FLogger) WriteLog(lv logLv, formatMsg string, other ...interface{}) error {
	if !l.start {
		return errors.New("logger do not start")
	}
	info := newLog(2)
	info.Level = lv
	info.Msg = fmt.Sprintf(formatMsg, other...)
	select {
	case ch <- info:
	default:
	}
	return nil
}

// create log struct
func newLog(caller int) *logC {
	_, f, ln, _ := runtime.Caller(caller)
	return &logC{
		Time: time.Now().Format("2006-01-02 15:04:05"),
		File: f,
		Line: ln,
	}
}

// write log
func (c *logC) write(path string) {
	f := getFile(c.Level, path)
	f.WriteString(fmt.Sprintf("[%s][%s]\n[%s:%d]\n[message]%s\n", level[c.Level], c.Time, c.File, c.Line, c.Msg))
}

// get file for write
func getFile(lv logLv, path string) *os.File {
	fName := level[lv] + "." + time.Now().Format("20060102") + ".log"
	f := files[lv]
	if f == nil {
		f, _ = os.OpenFile(path+fName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		files[lv] = f
	}
	if f.Name() != fName {
		f.Close()
		f = nil
		f, _ = os.OpenFile(path+fName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		files[lv] = f
	}
	return f
}

// close logger
func (l *FLogger) Close() {
	close(ch)
	l.wg.Wait()
	l.start = false
	for i := L_DEBUG; i < 6; i++ {
		if files[i] != nil {
			files[i].Close()
		}
	}
}
