package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger is a structured logger that writes to both file and stdout
type Logger struct {
	logger     *log.Logger
	file       *os.File
	logDir     string
	logFile    string
	logMaxSize int64
}

// LogLevel represents the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// Fields represents structured log fields
type Fields map[string]interface{}

// New creates a new Logger instance
func New(logDir, logFile string, maxSizeMB int64) (*Logger, error) {
	// Create log directory if not exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logPath := filepath.Join(logDir, logFile)

	// Open log file
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create multi-writer for both file and stdout
	writer := io.MultiWriter(os.Stdout, file)

	logger := log.New(writer, "", 0)
	logger.SetFlags(0) // We use our own format

	return &Logger{
		logger:     logger,
		file:      file,
		logDir:    logDir,
		logFile:   logFile,
		logMaxSize: maxSizeMB * 1024 * 1024,
	}, nil
}

// Close closes the logger file
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// checkRotation checks if log file needs rotation
func (l *Logger) checkRotation() error {
	if l.file == nil {
		return nil
	}

	info, err := l.file.Stat()
	if err != nil {
		return err
	}

	if info.Size() >= l.logMaxSize {
		// Rotate log file
		now := time.Now().Format("20060102-150405")
		backupPath := filepath.Join(l.logDir, fmt.Sprintf("%s.%s", l.logFile, now))

		l.file.Close()

		if err := os.Rename(filepath.Join(l.logDir, l.logFile), backupPath); err != nil {
			return fmt.Errorf("failed to rotate log: %w", err)
		}

		newFile, err := os.OpenFile(filepath.Join(l.logDir, l.logFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to create new log file: %w", err)
		}

		l.file = newFile

		// Update logger writer
		writer := io.MultiWriter(os.Stdout, l.file)
		l.logger.SetOutput(writer)
	}

	return nil
}

// formatLog formats a log entry with consistent format
func formatLog(level LogLevel, msg string, fields Fields) string {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	levelStr := levelNames[level]

	if fields == nil {
		return fmt.Sprintf("[%s] [%s] %s\n", now, levelStr, msg)
	}

	fieldStr := ""
	for k, v := range fields {
		fieldStr += fmt.Sprintf(" %s=%v", k, v)
	}

	return fmt.Sprintf("[%s] [%s] %s%s\n", now, levelStr, msg, fieldStr)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields Fields) {
	l.checkRotation()
	l.logger.Print(formatLog(DEBUG, msg, fields))
}

// Info logs an info message
func (l *Logger) Info(msg string, fields Fields) {
	l.checkRotation()
	l.logger.Print(formatLog(INFO, msg, fields))
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields Fields) {
	l.checkRotation()
	l.logger.Print(formatLog(WARN, msg, fields))
}

// Error logs an error message
func (l *Logger) Error(msg string, fields Fields) {
	l.checkRotation()
	l.logger.Print(formatLog(ERROR, msg, fields))
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields Fields) {
	l.checkRotation()
	l.logger.Print(formatLog(FATAL, msg, fields))
}

// WithFields returns a new entry with fields attached
func (l *Logger) WithFields(fields Fields) *Entry {
	return &Entry{logger: l, fields: fields}
}

// Entry is a log entry with pre-attached fields
type Entry struct {
	logger *Logger
	fields Fields
}

// Debug logs a debug message with fields
func (e *Entry) Debug(msg string) {
	e.logger.Debug(msg, e.fields)
}

// Info logs an info message with fields
func (e *Entry) Info(msg string) {
	e.logger.Info(msg, e.fields)
}

// Warn logs a warning message with fields
func (e *Entry) Warn(msg string) {
	e.logger.Warn(msg, e.fields)
}

// Error logs an error message with fields
func (e *Entry) Error(msg string) {
	e.logger.Error(msg, e.fields)
}

// Fatal logs a fatal message with fields and exits
func (e *Entry) Fatal(msg string) {
	e.logger.Fatal(msg, e.fields)
}