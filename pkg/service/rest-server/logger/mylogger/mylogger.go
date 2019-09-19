package mylogger

import "log"

type AppLogger struct {
	loglevel int
}
type LogLevel int

const (
	Info  LogLevel = 0
	Debug LogLevel = 1
	Error LogLevel = 2
)

func Init(level int) *AppLogger {
	return &AppLogger{level}
}

func (l *AppLogger) Log(caller string, level LogLevel, message interface{}) {
	var logNames = [...]string{
		"Info",
		"Debug",
		"Error",
	}
	log.Printf("[%v][%v] - %v", logNames[level], caller, message)
}
