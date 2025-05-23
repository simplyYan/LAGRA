package lagra

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// LogType represents the type of log message.
type LogType string

const (
	Info  LogType = "INFO"
	Warn  LogType = "WARN"
	Error LogType = "ERROR"
)

// LagraConfig holds the configuration for the Lagra logger.
type LagraConfig struct {
	LogFile    string // Path to the log file.
	EnableColor bool   // Enable colored terminal output.
}

// Lagra is the main struct for the logger.
type Lagra struct {
	config      LagraConfig
	logFile     *os.File
	logBuffer   []string
	logCounter  int32
	logMutex    sync.Mutex
	stopChannel chan struct{}
}

// NewLagra initializes a new Lagra instance.
func NewLagra(config LagraConfig) (*Lagra, error) {
	lagra := &Lagra{
		config:      config,
		logBuffer:   make([]string, 0, 100),
		stopChannel: make(chan struct{}),
	}

	if config.LogFile != "" {
		file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
		lagra.logFile = file
		go lagra.periodicFlush()
	}

	return lagra, nil
}

// Close shuts down the logger and flushes any remaining logs.
func (l *Lagra) Close() {
	close(l.stopChannel)
	l.flushLogBuffer()
	if l.logFile != nil {
		l.logFile.Close()
	}
}

// periodicFlush periodically flushes the log buffer to the file.
func (l *Lagra) periodicFlush() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.flushLogBuffer()
		case <-l.stopChannel:
			return
		}
	}
}

// flushLogBuffer writes buffered log messages to the log file.
func (l *Lagra) flushLogBuffer() {
	l.logMutex.Lock()
	defer l.logMutex.Unlock()

	if l.logFile != nil && len(l.logBuffer) > 0 {
		for _, log := range l.logBuffer {
			l.logFile.WriteString(log)
		}
		l.logBuffer = l.logBuffer[:0]
		atomic.StoreInt32(&l.logCounter, 0)
	}
}

// colorize applies ANSI color codes to a text.
func (l *Lagra) colorize(text string, logType LogType) string {
	if !l.config.EnableColor {
		return text
	}

	var colorCode string
	switch logType {
	case Info:
		colorCode = "\033[32m" // Green
	case Warn:
		colorCode = "\033[33m" // Yellow
	case Error:
		colorCode = "\033[31m" // Red
	default:
		colorCode = "\033[0m"  // Reset
	}

	return fmt.Sprintf("%s%s\033[0m", colorCode, text)
}

// send logs a message with the specified type.
func (l *Lagra) send(logType LogType, message string, customLogPath ...string) error {
	l.logMutex.Lock()
	defer l.logMutex.Unlock()

	logMessage := fmt.Sprintf("%s - %s - %s\n", time.Now().Format("15:04.05.9 - 02/01/2006"), logType, message)
	logMessageColored := l.colorize(logMessage, logType)

	fmt.Print(logMessageColored)

	// Use custom log file path if provided, otherwise use the configured log file
	logFilePath := l.config.LogFile
	if len(customLogPath) > 0 && customLogPath[0] != "" {
		logFilePath = customLogPath[0]
	}

	if logFilePath != "" {
		if l.logFile != nil {
			l.logBuffer = append(l.logBuffer, logMessage) // Save plain text (no ANSI) to file
			atomic.AddInt32(&l.logCounter, 1)
			if l.logCounter >= 100 {
				l.flushLogBuffer()
			}
		} else {
			fmt.Println("Log file is not set. Message will not be logged to a file.")
		}
	}

	return nil
}

// Info logs an informational message.
func (l *Lagra) Info(message string) {
	l.send(Info, message)
}

// Warn logs a warning message.
func (l *Lagra) Warn(message string) {
	l.send(Warn, message)
}

// Error logs an error message.
func (l *Lagra) Error(message string) {
	l.send(Error, message)
}

// InfoWithCustomFile logs an informational message to a custom log file.
func (l *Lagra) InfoWithCustomFile(message, customLogPath string) {
	l.send(Info, message, customLogPath)
}

// WarnWithCustomFile logs a warning message to a custom log file.
func (l *Lagra) WarnWithCustomFile(message, customLogPath string) {
	l.send(Warn, message, customLogPath)
}

// ErrorWithCustomFile logs an error message to a custom log file.
func (l *Lagra) ErrorWithCustomFile(message, customLogPath string) {
	l.send(Error, message, customLogPath)
}
