// Package logging provides leveled, timezone-aware logging to stderr in
// text or JSON format, per SPEC/06-configuration.md.
package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

func ParseLevel(s string) Level {
	switch s {
	case "debug":
		return Debug
	case "warn":
		return Warn
	case "error":
		return Error
	default:
		return Info
	}
}

func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	default:
		return "INFO"
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	jsonFmt  bool
	loc      *time.Location
}

func New(out io.Writer, minLevel Level, jsonFormat bool, loc *time.Location) *Logger {
	if loc == nil {
		loc = time.UTC
	}
	return &Logger{out: out, minLevel: minLevel, jsonFmt: jsonFormat, loc: loc}
}

func (l *Logger) log(level Level, msg string) {
	if level < l.minLevel {
		return
	}
	now := time.Now().In(l.loc)
	if l.jsonFmt {
		enc := json.NewEncoder(l.out)
		_ = enc.Encode(map[string]string{
			"timestamp": now.Format("2006-01-02 15:04:05"),
			"level":     level.String(),
			"message":   msg,
		})
		return
	}
	fmt.Fprintf(l.out, "[%s] [%s] %s\n", now.Format("2006-01-02 15:04:05"), level.String(), msg)
}

func (l *Logger) Debugf(format string, args ...any) { l.log(Debug, fmt.Sprintf(format, args...)) }
func (l *Logger) Infof(format string, args ...any)  { l.log(Info, fmt.Sprintf(format, args...)) }
func (l *Logger) Warnf(format string, args ...any)  { l.log(Warn, fmt.Sprintf(format, args...)) }
func (l *Logger) Errorf(format string, args ...any) { l.log(Error, fmt.Sprintf(format, args...)) }
