package mylogger

import "log"

type AppLogger struct {
	loglevel int
}
type LogLevel int

var logNames = [...]string{
	"Info",
	"Debug",
	"Error",
}

const (
	Info  LogLevel = 0
	Debug LogLevel = 1
	Error LogLevel = 2
)

func Init(level int) *AppLogger {
	return &AppLogger{level}
}

func (l *AppLogger) Log(caller string, level LogLevel, message interface{}) {

	log.Printf("[%v][%v] - %v", logNames[level], caller, message)
}
