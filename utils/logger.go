package utils

import (
	"github.com/google/logger"
	ioutil "io/ioutil"
	"log"
)

func Init(name string, verbose bool) {
	defer logger.Init(name, verbose, true, ioutil.Discard).Close()
	logger.SetFlags(log.Ldate | log.Ltime | log.LUTC)

}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}
func Error(v ...interface{}) {
	logger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func Errorln(v ...interface{}) {
	logger.Errorln(v...)
}

func Info(v ...interface{}) {
	logger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}
func Warning(v ...interface{}) {
	logger.Warning(v...)
}

func Warningf(format string, v ...interface{}) {
	logger.Warningf(format, v...)
}
