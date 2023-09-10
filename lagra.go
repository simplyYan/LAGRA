package lagra

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
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

func (l *Lagra) send(logType LogType) func(ctx context.Context, message string) error {
	var textColor *color.Color

	switch logType {
	case Info:
		textColor = color.New(color.FgGreen)
	case Warn:
		textColor = color.New(color.FgYellow)
	case Error:
		textColor = color.New(color.FgRed)
	default:
		textColor = color.New(color.Reset)
	}

	return func(ctx context.Context, message string) error {
		logMessage := fmt.Sprintf("%s - %s - %s\n", time.Now().Format("15:04.05.9 - 02/01/2006"), logType, message)
		textColor.Print(logMessage)
		fmt.Print(logMessage) // Imprime sem cores no arquivo de log

		if l.logFile != nil {
			_, err := l.logFile.WriteString(logMessage)
			if err != nil {
				return err // Retornar erro se ocorrer um problema ao escrever no arquivo de log
			}
		}

		return nil // Sem erro
	}
}

// Info logs an informational message.
func (l *Lagra) Info(ctx context.Context, message string) error {
	return l.send(Info)(ctx, message)
}

// Warn logs a warning message.
func (l *Lagra) Warn(ctx context.Context, message string) error {
	return l.send(Warn)(ctx, message)
}

// Error logs an error message.
func (l *Lagra) Error(ctx context.Context, message string) error {
	return l.send(Error)(ctx, message)
}

func (l *Lagra) save() error {
	if l.logFile == nil {
		var err error
		l.logFile, err = os.Create("log.lagra")
		if err != nil {
			return err // Retornar erro se ocorrer um problema ao criar o arquivo de log
		}
	}
	return nil // Sem erro
}

func (l *Lagra) Close() error {
	if l.logFile != nil {
		err := l.logFile.Close()
		if err != nil {
			return err // Retornar erro se ocorrer um problema ao fechar o arquivo de log
		}
	}
	return nil // Sem erro
}
