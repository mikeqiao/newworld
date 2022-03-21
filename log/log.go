package log

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	conf "github.com/mikeqiao/newworld/config"
)

const (
	debugLevel   = 0
	releaseLevel = 1
	warningLevel = 2
	errorLevel   = 3
	fatalLevel   = 4
)

const (
	printDebugLevel   = "[debug  ] "
	printReleaseLevel = "[release] "
	printWarningLevel = "[warning] "
	printErrorLevel   = "[error  ] "
	printFatalLevel   = "[fatal  ] "
)

type Logger struct {
	level      int
	baseLogger *log.Logger
	baseFile   *os.File
}

func New(strLevel string, pathName, logName string, flag int) (*Logger, error) {
	// level
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level = debugLevel
	case "release":
		level = releaseLevel
	case "warning":
		level = warningLevel
	case "error":
		level = errorLevel
	case "fatal":
		level = fatalLevel
	default:
		return nil, errors.New("unknown level: " + strLevel)
	}
	// logger
	var baseLogger *log.Logger
	var baseFile *os.File
	if pathName != "" {
		now := time.Now()
		filename := fmt.Sprintf("%d%02d%02d_%v.log",
			now.Year(),
			now.Month(),
			now.Day(),
			logName)
		file, err := os.OpenFile(path.Join(pathName, filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}

		baseLogger = log.New(file, "", flag)
		baseFile = file
	} else {
		baseLogger = log.New(os.Stdout, "", flag)
	}
	// new
	logger := new(Logger)
	logger.level = level
	logger.baseLogger = baseLogger
	logger.baseFile = baseFile
	return logger, nil
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.baseFile != nil {
		err := logger.baseFile.Close()
		if nil != err {
			fmt.Println("logger close err:%v", err)
		}
	}
	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = printLevel + format

	t := time.Now().String()
	var buffer bytes.Buffer
	buffer.WriteString(t[:26])
	_, file, line, ok := runtime.Caller(2)
	if ok {
		buffer.WriteString(" ")
		buffer.WriteString(file)
		buffer.WriteString("	line:")
		buffer.WriteString(strconv.Itoa(line))
		buffer.WriteString(":")
	}
	err := logger.baseLogger.Output(3, buffer.String()+fmt.Sprintf(format, a...))
	if nil != err {
		fmt.Println("log output, err:%v", err)
	}
	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Release(format string, a ...interface{}) {
	logger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func (logger *Logger) Warning(format string, a ...interface{}) {
	logger.doPrintf(warningLevel, printWarningLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

var gLogger, _ = New("debug", "", "", log.LstdFlags)

// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if logger != nil {
		gLogger = logger
	}
}

func Debug(format string, a ...interface{}) {
	gLogger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func Release(format string, a ...interface{}) {
	gLogger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func Warning(format string, a ...interface{}) {
	gLogger.doPrintf(warningLevel, printWarningLevel, format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func Close() {
	gLogger.Close()
}

func Init() {
	logger, err := New(conf.Conf.LogLevel, conf.Conf.LogPath, conf.Conf.SInfo.Name, int(conf.Conf.LogFlag))
	if err != nil {
		panic(err)
	}
	Export(logger)
}
