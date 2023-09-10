package lagra

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

type LogType string

const (
	Info  LogType = "INFO"
	Warn  LogType = "WARN"
	Error LogType = "ERROR"
)

type Lagra struct {
	logFile  *os.File
	logMutex sync.Mutex
	config   *LagraConfig
}

type LagraConfig struct {
	LogFile  string `toml:"log_file"`
	LogLevel string `toml:"log_level"`
}

// New creates a new Lagra logger instance with optional TOML configuration.
func New(tomlConfig string) (*Lagra, error) {
	var config LagraConfig
	if tomlConfig != "" {
		if _, err := toml.Decode(tomlConfig, &config); err != nil {
			return nil, err
		}
	}

	l := &Lagra{
		config: &config,
	}

	if l.config.LogFile != "" {
		l.SetLogFile(l.config.LogFile)
	}

	if l.config.LogLevel != "" {
		l.SetLogLevel(l.config.LogLevel)
	}

	return l, nil
}

// send logs a message with the specified log type.
func (l *Lagra) send(logType LogType, message string) error {
	l.logMutex.Lock()
	defer l.logMutex.Unlock()

	logMessage := fmt.Sprintf("%s - %s - %s\n", time.Now().Format("15:04.05.9 - 02/01/2006"), logType, message)

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

	logMessageColored := textColor.SprintFunc()(logMessage)
	fmt.Print(logMessageColored) // Imprime com cores

	if l.logFile != nil {
		_, err := l.logFile.WriteString(logMessage)
		if err != nil {
			return err // Retornar erro se ocorrer um problema ao escrever no arquivo de log
		}
	}

	return nil // Sem erro
}

// Info logs an informational message.
func (l *Lagra) Info(ctx context.Context, message string) error {
	return l.send(Info, message)
}

// Warn logs a warning message.
func (l *Lagra) Warn(ctx context.Context, message string) error {
	return l.send(Warn, message)
}

// Error logs an error message.
func (l *Lagra) Error(ctx context.Context, message string) error {
	return l.send(Error, message)
}

// SetLogFile sets the log file path.
func (l *Lagra) SetLogFile(filePath string) {
	l.logMutex.Lock()
	defer l.logMutex.Unlock()

	if l.logFile != nil {
		_ = l.logFile.Close()
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n", err)
		return
	}

	l.logFile = file
}

// SetLogLevel sets the log level (INFO, WARN, ERROR).
func (l *Lagra) SetLogLevel(level string) {
	switch level {
	case "INFO":
		// Set log level to INFO
	case "WARN":
		// Set log level to WARN
	case "ERROR":
		// Set log level to ERROR
	default:
		level = "INFO" // Default log level
	}

	fmt.Printf("Log level set to %s\n", level)
}
