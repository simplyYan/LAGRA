package lagra

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type LogType string

const (
	Info  LogType = "INFO"
	Warn  LogType = "WARN"
	Error LogType = "ERROR"
)

type LagraConfig struct {
	LogFile        string            // Path to log file.
	EnableColor    bool              // Enable colored output.
	MinLevel       LogType           // Minimum level to log.
	ContextFields  map[string]string // Extra fields to include in all logs.
	JSONFormat     bool              // Output logs as JSON.
	TimeFormat     string            // Custom time format.
	CustomColors   map[LogType]string// Custom ANSI colors.
	Async          bool              // Asynchronous logging.
}

type Lagra struct {
	config      LagraConfig
	logFile     *os.File
	logChan     chan string
	stopChan    chan struct{}
	wg          sync.WaitGroup
	hooks       []func(LogType, string)
	mu          sync.Mutex
}

func NewLagra(config LagraConfig) (*Lagra, error) {
	lagra := &Lagra{
		config:   config,
		logChan:  make(chan string, 100),
		stopChan: make(chan struct{}),
	}

	if config.LogFile != "" {
		file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		lagra.logFile = file
	}

	if config.Async {
		lagra.wg.Add(1)
		go lagra.processLogs()
	}

	return lagra, nil
}

func (l *Lagra) Close() {
	close(l.stopChan)
	if l.config.Async {
		l.wg.Wait()
	}
	if l.logFile != nil {
		l.logFile.Close()
	}
}

func (l *Lagra) AddHook(hook func(LogType, string)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hooks = append(l.hooks, hook)
}

func (l *Lagra) log(level LogType, msg string) {
	if !l.shouldLog(level) {
		return
	}

	entry := l.formatLog(level, msg)

	// Run hooks
	for _, hook := range l.hooks {
		hook(level, msg)
	}

	if l.config.Async {
		l.logChan <- entry
	} else {
		fmt.Print(entry)
		if l.logFile != nil {
			l.logFile.WriteString(entry)
		}
	}
}

func (l *Lagra) processLogs() {
	defer l.wg.Done()
	for {
		select {
		case entry := <-l.logChan:
			fmt.Print(entry)
			if l.logFile != nil {
				l.logFile.WriteString(entry)
			}
		case <-l.stopChan:
			return
		}
	}
}

func (l *Lagra) shouldLog(level LogType) bool {
	order := map[LogType]int{Info: 1, Warn: 2, Error: 3}
	return order[level] >= order[l.config.MinLevel]
}

func (l *Lagra) formatLog(level LogType, msg string) string {
	timestamp := time.Now().Format(l.config.TimeFormat)
	logData := map[string]interface{}{
		"time":    timestamp,
		"level":   level,
		"message": msg,
	}

	for k, v := range l.config.ContextFields {
		logData[k] = v
	}

	var output string
	if l.config.JSONFormat {
		bytes, _ := json.Marshal(logData)
		output = string(bytes) + "\n"
	} else {
		colored := l.colorize(fmt.Sprintf("%s [%s] %s", timestamp, level, msg), level)
		output = colored + "\n"
	}

	return output
}

func (l *Lagra) colorize(text string, level LogType) string {
	if !l.config.EnableColor {
		return text
	}

	color := l.config.CustomColors[level]
	if color == "" {
		switch level {
		case Info:
			color = "\033[32m"
		case Warn:
			color = "\033[33m"
		case Error:
			color = "\033[31m"
		default:
			color = "\033[0m"
		}
	}

	return fmt.Sprintf("%s%s\033[0m", color, text)
}

func (l *Lagra) Info(msg string)  { l.log(Info, msg) }
func (l *Lagra) Warn(msg string)  { l.log(Warn, msg) }
func (l *Lagra) Error(msg string) { l.log(Error, msg) }
