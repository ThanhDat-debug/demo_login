package logger

import (
	"log"
	"os"
)

// Logger là một wrapper đơn giản cho log package của Go
type Logger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
}

// NewLogger tạo một instance mới của Logger
func NewLogger() *Logger {
	return &Logger{
		InfoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		DebugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info ghi log thông tin
func (l *Logger) Info(v ...interface{}) {
	l.InfoLogger.Println(v...)
}

// Error ghi log lỗi
func (l *Logger) Error(v ...interface{}) {
	l.ErrorLogger.Println(v...)
}

// Debug ghi log debug
func (l *Logger) Debug(v ...interface{}) {
	l.DebugLogger.Println(v...)
}

// Infof ghi log thông tin với format
func (l *Logger) Infof(format string, v ...interface{}) {
	l.InfoLogger.Printf(format, v...)
}

// Errorf ghi log lỗi với format
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.ErrorLogger.Printf(format, v...)
}

// Debugf ghi log debug với format
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.DebugLogger.Printf(format, v...)
}
