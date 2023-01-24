package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/hashicorp/logutils"
	"github.com/pkg/errors"
)

const (
	unknownFileName       string = "unknow file"
	logFilenameSkipFrames int    = 2
)

type causer interface {
	Cause() error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Newline
// helps control newline character
// use '\r' for lambda
// use '\n' for local debug.
//
//nolint:gochecknoglobals
var useCr bool

// SetLogLevels
//
//	@param logLevels
//	@param minLevel
//	@param crNewline
func SetLogLevels(logLevels []string, minLevel string, crNewline bool) {
	levels := []logutils.LogLevel{}
	for _, l := range logLevels {
		levels = append(levels, logutils.LogLevel(l))
	}
	filter := &logutils.LevelFilter{
		Levels:   levels,
		MinLevel: logutils.LogLevel(minLevel),
		Writer:   os.Stdout,
	}
	log.SetOutput(filter)
	useCr = crNewline
}

// Debug
//
//	@param format
//	@param v
func Debug(format string, v ...any) {
	f := getLogFile()
	s := fmt.Sprintf(format, v...)
	s = fmt.Sprintf("%s %s %s", "[DEBUG]", f, s)
	s = newline(s)
	log.Println(s)
}

// Info
//
//	@param format
//	@param v
func Info(format string, v ...any) {
	f := getLogFile()
	s := fmt.Sprintf(format, v...)
	s = fmt.Sprintf("%s %s %s", "[INFO]", f, s)
	s = newline(s)
	log.Println(s)
}

// Warn
//
//	@param format
//	@param v
func Warn(format string, v ...any) {
	f := getLogFile()
	s := fmt.Sprintf(format, v...)
	s = fmt.Sprintf("%s %s %s", "[WARN]", f, s)
	s = newline(s)
	log.Println(s)
}

// Error
//
//	@param format
//	@param err
//	@param v
func Error(format string, err error, v ...any) {
	f := getLogFile()
	s := fmt.Sprintf(format, v...)
	var t string
	if err == nil {
		t = "{nil error}"
	} else if tracer := getTracableRoot(err); tracer == nil {
		t = fmt.Sprintf("cause: %+v", err)
	} else {
		traces := []string{}
		traces = append(traces,
			fmt.Sprintf("cause: %s", err.Error()),
			"trace: [",
		)
		for _, st := range tracer.StackTrace() {
			traces = append(traces, fmt.Sprintf("  - %+v", st))
		}
		traces = append(traces, "]")
		t = strings.Join(traces, "\n")
	}
	l := newline(s + "\n" + t)
	l = fmt.Sprintf("%s %s %s", "[ERROR]", f, l)
	log.Println(l)
}

// Pretty
//
//	@param v any object
//	@return string json string
func Pretty(v interface{}) string {
	if v == nil {
		return "{nil}"
	}
	vjson, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "pretty_error"
	}
	return string(vjson)
}

// getTracableRoot
//
//	@param err
//	@return stackTracer
//
//nolint:ireturn
func getTracableRoot(err error) stackTracer {
	if err == nil {
		return nil
	}
	var tracer stackTracer
	prev := err
	for err != nil {
		//nolint:errorlint
		cause, ok := err.(causer)
		if !ok {
			break
		}
		//nolint:errorlint
		_, ok2 := err.(stackTracer)
		if ok2 {
			prev = err
		}
		err = cause.Cause()
	}
	ok := errors.As(prev, &tracer)
	if ok {
		return tracer
	}
	return nil
}

// getLogFile
//
//	@return string
func getLogFile() string {
	_, file, line, ok := runtime.Caller(logFilenameSkipFrames)
	if !ok {
		return unknownFileName
	}
	index := strings.LastIndex(file, "/")
	if index >= 0 && index+1 < len(file) {
		sf := file[index+1:]
		return sf + "#" + strconv.Itoa(line)
	}
	return unknownFileName
}

// newline change newline characters
//
//	@param text
//	@return string
func newline(text string) string {
	if !useCr {
		return text
	}
	line := strings.ReplaceAll(text, "\n", "\r")
	return line
}
