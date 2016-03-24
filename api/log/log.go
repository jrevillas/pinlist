package log

import (
	"io"
	"os"

	"github.com/op/go-logging"
)

const logFormat = "%{color}%{shortfile} ▶ [%{level}%]{color:reset}: %{message}"

var (
	logFile = os.Getenv("LOGFILE")
	log     *logging.Logger
)

func init() {
	if logFile == "" {
		logFile = os.DevNull
	}

	log = setupLogger(logFileWriter(logFile))
}

func logFileWriter(logFile string) io.Writer {
	w, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return w
}

func setupLogger(w io.Writer) *logging.Logger {
	logging.SetFormatter(logging.MustStringFormatter(logFormat))

	stderrBackend := logging.AddModuleLevel(logging.NewLogBackend(os.Stderr, "", 0))
	stderrBackend.SetLevel(logging.WARNING, "")
	fileBackend := logging.NewLogBackend(w, "", 0)

	logging.SetBackend(stderrBackend, fileBackend)

	return logging.MustGetLogger("pinlist")
}

func Debug(msg string, args ...interface{}) {
	log.Debugf(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	log.Warningf(msg, args...)
}

func Error(msg string, args ...interface{}) {
	log.Errorf(msg, args...)
}

func Info(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}

func Critical(msg string, args ...interface{}) {
	log.Criticalf(msg, args...)
}

func Err(err error) {
	Error(err.Error())
}
