package core

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
)

// These are a consistent set of error codes instead of using random non-zero integers
const (
	FATAL     = 5
	ERROR     = 4
	WARN      = 3
	IMPORTANT = 2
	INFO      = 1
	DEBUG     = 0
)

// LogColors sets the color for each type of logging output
var LogColors = map[int]*color.Color{
	FATAL:     color.New(color.FgRed).Add(color.Bold),
	ERROR:     color.New(color.FgRed),
	WARN:      color.New(color.FgYellow),
	IMPORTANT: color.New(color.Bold),
	DEBUG:     color.New(color.FgCyan).Add(color.Faint),
}

// Logger holds specific configuration data for the logging
type Logger struct {
	sync.Mutex

	debug  bool
	silent bool
}

// SetSilent will configure the logger to not display any realtime output to stdout
func (l *Logger) SetSilent(s bool) {
	l.silent = s
}

// SetDebug will configure the logger to enable debug output to be set to stdout
func (l *Logger) SetDebug(d bool) {
	l.debug = d
}

// Log is a generic printer for sending data to stdout. It does not do traditional syslog logging
func (l *Logger) Log(level int, format string, args ...interface{}) {
	l.Lock()
	defer l.Unlock()
	if level == DEBUG && l.debug == false {
		return
	} else if level < ERROR && l.silent == true {
		return
	}

	if c, ok := LogColors[level]; ok {
		_, _ = c.Printf(format, args...)
	} else {
		fmt.Printf(format, args...)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

// Fatal prints a fatal level log message to stdout
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.Log(FATAL, format, args...)
}

// Error prints an error level log message to stdout
func (l *Logger) Error(format string, args ...interface{}) {
	l.Log(ERROR, format, args...)
}

// Warn prints a warn level log message to stdout
func (l *Logger) Warn(format string, args ...interface{}) {
	l.Log(WARN, format, args...)
}

// Important prints an important level log message to stdout
func (l *Logger) Important(format string, args ...interface{}) {
	l.Log(IMPORTANT, format, args...)
}

// Info prints an info level log message to stdout
func (l *Logger) Info(format string, args ...interface{}) {
	l.Log(INFO, format, args...)
}

// Debug prints a debug level log message to stdout
func (l *Logger) Debug(format string, args ...interface{}) {
	l.Log(DEBUG, format, args...)
}

// InitLogger will initialize the logger for the session
func (s *Session) InitLogger() {
	s.Out = &Logger{}
	s.Out.SetDebug(s.Debug)
	s.Out.SetSilent(s.Silent)
}
