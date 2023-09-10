package lagra

import (
	"fmt"
	"os"
	"time"
)

type LogType string

const (
	Info  LogType = "INFO"
	Warn  LogType = "WARN"
	Error LogType = "ERROR"
)

type Lagra struct {
	logFile *os.File
}

func New() *Lagra {
	return &Lagra{}
}

func (l *Lagra) send(logType LogType) func(message string) {
	return func(message string) {
		logMessage := fmt.Sprintf("%s - %s - %s\n", time.Now().Format("15:04.05.9 - 02/01/2006"), logType, message)
		fmt.Print(logMessage)
		if l.logFile != nil {
			l.logFile.WriteString(logMessage)
		}
	}
}

func (l *Lagra) save() {
	if l.logFile == nil {
		var err error
		l.logFile, err = os.Create("log.lagra")
		if err != nil {
			fmt.Println("Error creating log file:", err)
		}
	}
}

func (l *Lagra) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}
