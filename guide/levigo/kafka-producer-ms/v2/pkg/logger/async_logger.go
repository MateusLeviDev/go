package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/MateusLeviDev/pkg/constants"
	"github.com/petermattis/goid"
)

var (
	Async   *log.Logger
	logCh   chan logMessage
	closeCh chan struct{}
	flushCh chan struct{}
)

type logMessage struct {
	level string
	msg   string
}

func init() {
	file, err := os.OpenFile(constants.LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	Async = log.New(io.MultiWriter(os.Stdout, file), constants.PREFIX, log.Ldate|log.Ltime|log.Lmicroseconds)
	logCh = make(chan logMessage, 100)
	closeCh = make(chan struct{})
	flushCh = make(chan struct{})

	go logWorker()
}

func logWorker() {
	for {
		select {
		case logMsg := <-logCh:
			Async.SetPrefix(logMsg.level)
			Async.Println(logMsg.msg)
		case <-closeCh:
			for len(logCh) > 0 {
				logMsg := <-logCh
				Async.SetPrefix(logMsg.level)
				Async.Println(logMsg.msg)
			}
			flushCh <- struct{}{}
			return
		}
	}
}

func logLocationAsync() string {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return constants.WTF
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = funcName[strings.LastIndex(funcName, "/")+1:]
	return fmt.Sprintf("[%s] [line:%d] [function:%s] [GoroutineID: %d] ", path.Base(file), line, funcName, goid.Get())
}

func logAsync(level string, v ...interface{}) {
	logCh <- logMessage{
		level: level,
		msg:   fmt.Sprintf("%s%s", logLocationAsync(), fmt.Sprint(v...)),
	}
}

func DebugAsync(v ...interface{}) {
	logAsync(constants.DEBUG, v...)
}

func InfoAsync(v ...interface{}) {
	logAsync(constants.INFO, v...)
}

func WarningAsync(v ...interface{}) {
	logAsync(constants.WARNING, v...)
}

func ErrorAsync(v ...interface{}) {
	logAsync(constants.ERROR, v...)
}

func Close() {
	close(closeCh)
	<-flushCh
}
