// date: 2019-03-12
package log

import (
	"github.com/Jarvens/Exchange-Agent/util/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
)

func LumberJackLogger(filePath string, maxSize int, maxBackups int, maxAge int) *lumberjack.Logger {
	return &lumberjack.Logger{Filename: filePath,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups}
}

func InitLogToStdoutDebug() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func InitLogToStdout() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.WarnLevel)
}

func InitLogToFile() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	out := LumberJackLogger(config.WsErrorLogFilePath+config.WsErrorLogFileExtension,
		config.WsErrorLogMaxSize,
		config.WsErrorLogMaxBackups,
		config.WsErrorLogMaxAge)
	logrus.SetOutput(out)

	access := LumberJackLogger(config.WsAccessLogFilePath+config.WsAccessLogFileExtension,
		config.WsAccessLogMaxSize,
		config.WsAccessLogMaxBackups,
		config.WsAccessLogMaxAge)
	logrus.SetLevel(logrus.InfoLevel)
	//成功日志输出
	logrus.SetOutput(access)
	//错误日志输出
	logrus.SetLevel(logrus.InfoLevel)
}

func Init(environment string) {
	switch environment {
	case "DEVELOPMENT":
		InitLogToStdoutDebug()
	case "TEST":
		InitLogToFile()
	case "PRODUCTION":
		InitLogToFile()
	}
	logrus.Infof("Environment : %s", environment)
}

func Debug(msg string) {
	logrus.Debug(msg)
}

// Debugf logs a formatted message with debug log level.
func Debugf(msg string, args ...interface{}) {
	logrus.Debugf(msg, args...)
}

// Info logs a message with info log level.
func Info(msg string) {
	logrus.Info(msg)
}

// Infof logs a formatted message with info log level.
func Infof(msg string, args ...interface{}) {

	logrus.Infof(msg, args...)
}

// Warn logs a message with warn log level.
func Warn(msg string) {
	logrus.Warn(msg)
}

// Warnf logs a formatted message with warn log level.
func Warnf(msg string, args ...interface{}) {
	logrus.Warnf(msg, args...)
}

// Error logs a message with error log level.
func Error(msg string) {
	logrus.Error(msg)
}

// Errorf logs a formatted message with error log level.
func Errorf(msg string, args ...interface{}) {
	logrus.Errorf(msg, args...)
}

// Fatal logs a message with fatal log level.
func Fatal(msg string) {
	logrus.Fatal(msg)
}

// Fatalf logs a formatted message with fatal log level.
func Fatalf(msg string, args ...interface{}) {
	logrus.Fatalf(msg, args...)
}

// Panic logs a message with panic log level.
func Panic(msg string) {
	logrus.Panic(msg)
}

// Panicf logs a formatted message with panic log level.
func Panicf(msg string, args ...interface{}) {
	logrus.Panicf(msg, args...)
}

// log response body data for debugging
func DebugResponse(response *http.Response) string {
	bodyBuffer := make([]byte, 5000)
	var str string
	count, err := response.Body.Read(bodyBuffer)
	for ; count > 0; count, err = response.Body.Read(bodyBuffer) {
		if err != nil {
		}
		str += string(bodyBuffer[:count])
	}
	Debugf("response data : %v", str)
	return str
}

func init() {
	Init("PRODUCTION")
}
