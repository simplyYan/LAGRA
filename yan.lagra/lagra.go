package lagra

import (
	"fmt"
	"os"
	"time"
)

type Lagra struct {
	logFile *os.File
}

type LogType string

const (
	Info  LogType = "INFO"
	Warn  LogType = "WARN"
	Error LogType = "ERROR"
)

func New() *Lagra {
	return &Lagra{}
}

func (l *Lagra) send(logType LogType, message string) {
	currentTime := time.Now().Format("15:04:05.000 - 02/01/2006")
	fmt.Printf("%s [%s] %s\n", currentTime, logType, message)
}

func (l *Lagra) save() {
	if l.logFile == nil {
		var err error
		l.logFile, err = os.Create("log.lagra")
		if err != nil {
			fmt.Printf("Failed to create log file: %v\n", err)
			return
		}
	}
}

func (l *Lagra) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}
