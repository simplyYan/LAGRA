package lagra

import (
	"context"
	"fmt"
	"os"
        "strings"
	"sync"
	"time"
	"sync/atomic"
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
	logBuffer []string
	logCounter int32
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
		logBuffer: make([]string, 0, 100),
	}

	if l.config.LogFile != "" {
		l.SetLogFile(l.config.LogFile)
	}

	if l.config.LogLevel != "" {
		l.SetLogLevel(l.config.LogLevel)
	}

	go l.flushLogBuffer()

	return l, nil
}

// send logs a message with the specified log type and an optional custom log file path.
func (l *Lagra) send(logType LogType, message string, customLogPath ...string) error {
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
	fmt.Print(logMessageColored) // Print with colors

	// Use custom log file path if provided, otherwise use the configured log file
	logFilePath := l.config.LogFile
	if len(customLogPath) > 0 && customLogPath[0] != "" {
		logFilePath = customLogPath[0]
	}

	if logFilePath != "" {
		if l.logFile != nil {
			l.logBuffer = append(l.logBuffer, logMessage)
			atomic.AddInt32(&l.logCounter, 1)
			if l.logCounter >= 100 {
				l.flushLogBuffer()
			}
		} else {
			fmt.Println("Log file is not set. Message will not be logged to a file.")
		}
	}

	return nil // No error
}

// Info logs an informational message with optional custom log file path.
func (l *Lagra) Info(ctx context.Context, message string, customLogPath ...string) error {
	return l.send(Info, message, customLogPath...)
}

// Warn logs a warning message with optional custom log file path.
func (l *Lagra) Warn(ctx context.Context, message string, customLogPath ...string) error {
	return l.send(Warn, message, customLogPath...)
}

// Error logs an error message with optional custom log file path.
func (l *Lagra) Error(ctx context.Context, message string, customLogPath ...string) error {
	return l.send(Error, message, customLogPath...)
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

// flushLogBuffer writes the log buffer to the log file and resets the buffer and counter.
func (l *Lagra) flushLogBuffer() {
	l.logMutex.Lock()
	defer l.logMutex.Unlock()

	if l.logFile != nil && len(l.logBuffer) > 0 {
		_, err := l.logFile.WriteString(strings.Join(l.logBuffer, ""))
		if err != nil {
			fmt.Printf("Failed to write to log file: %v\n", err)
		}
		l.logBuffer = l.logBuffer[:0]
		atomic.StoreInt32(&l.logCounter, 0)
	}
}

type ErrorCollector struct {
    errors []error
}

func New() *ErrorCollector {
    return &ErrorCollector{}
}

func (ec *ErrorCollector) N(err error) {
    if err != nil {
        ec.errors = append(ec.errors, err)
    }
}

func (ec *ErrorCollector) Handle() bool {
    return len(ec.errors) > 0
}

func (ec *ErrorCollector) Errors() []error {
    return ec.errors
}

type StrSelect struct {
    strMap map[string]string
}

func NewStrSelect() *StrSelect {
    return &StrSelect{
        strMap: make(map[string]string),
    }
}

func (s *StrSelect) SetStr(strName, value string) {
    s.strMap[strName] = value
}

func (s *StrSelect) SelectStr(strName, delimiter string) string {
    if str, ok := s.strMap[strName]; ok {
        startIdx := strings.Index(str, delimiter)
        endIdx := strings.LastIndex(str, delimiter)
        if startIdx != -1 && endIdx != -1 && startIdx < endIdx {
            return str[startIdx+len(delimiter) : endIdx]
        }
    }
    return ""
}
