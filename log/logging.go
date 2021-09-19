package log

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jsawatzky/go-common/internal"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelSevere
	LevelGlobal = -1
)

var (
	globalLevel   = LevelInfo
	defaultLogger = New("root")
	loggers       = make(map[string]Logger)
)

type Logger interface {
	Debug(f string, v ...interface{})
	Info(f string, v ...interface{})
	Warn(f string, v ...interface{})
	Error(f string, v ...interface{})
	Fatal(f string, v ...interface{})
	Panic(f string, v ...interface{})
	SetLevel(lvl int)
}

func New(name string) Logger {
	return &logger{
		Logger: log.New(os.Stdout, name+" ", log.Ldate|log.Ltime|log.LUTC|log.Lmsgprefix),
		level:  LevelGlobal,
	}
}

func SetGlobalLevel(level int) {
	if level < 0 || level > LevelSevere {
		panic("invalid log level")
	}
	globalLevel = level
}

func GetLogger(name string) Logger {
	if logger, ok := loggers[name]; ok {
		return logger
	}
	logger := New(name)
	loggers[name] = logger
	return logger
}

func Debug(f string, v ...interface{}) {
	defaultLogger.Debug(f, v...)
}

func Info(f string, v ...interface{}) {
	defaultLogger.Info(f, v...)
}

func Warn(f string, v ...interface{}) {
	defaultLogger.Warn(f, v...)
}

func Error(f string, v ...interface{}) {
	defaultLogger.Error(f, v...)
}

func Fatal(f string, v ...interface{}) {
	defaultLogger.Fatal(f, v...)
}

func Panic(f string, v ...interface{}) {
	defaultLogger.Panic(f, v...)
}

type logger struct {
	*log.Logger
	level int
}

func (l *logger) shouldLog(level int) bool {
	lvl := l.level
	if lvl < 0 {
		lvl = globalLevel
	}
	return lvl <= level
}

func (l *logger) Debug(f string, v ...interface{}) {
	if l.shouldLog(LevelDebug) {
		l.Printf("[DEBUG] %s", fmt.Sprintf(f, v...))
	}
}

func (l *logger) Info(f string, v ...interface{}) {
	if l.shouldLog(LevelInfo) {
		l.Printf("[INFO] %s", fmt.Sprintf(f, v...))
	}
}

func (l *logger) Warn(f string, v ...interface{}) {
	if l.shouldLog(LevelWarn) {
		l.Printf("[WARN] %s", fmt.Sprintf(f, v...))
	}
}

func (l *logger) Error(f string, v ...interface{}) {
	if l.shouldLog(LevelError) {
		l.Printf("[ERROR] %s", fmt.Sprintf(f, v...))
	}
}

func (l *logger) Fatal(f string, v ...interface{}) {
	l.Fatalf("[FATAL] %s", fmt.Sprintf(f, v...))
}

func (l *logger) Panic(f string, v ...interface{}) {
	l.Panicf("[PANIC] %s", fmt.Sprintf(f, v...))
}

func (l *logger) SetLevel(lvl int) {
	l.level = lvl
}

func Middleware(h http.Handler) http.Handler {
	logger := GetLogger("http")
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		resp := internal.RecordResponse(rw)
		h.ServeHTTP(resp, r)

		logger.Info("http request: \"%s %s %s\" %d %v %d", r.Method, r.URL.Path, r.Proto, resp.Status(), time.Since(start), resp.ResponseSize())
	})
}
