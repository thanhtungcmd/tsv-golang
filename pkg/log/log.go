package log

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	file   *os.File
	logger *log.Logger
}

var (
	instance *Logger
	once     sync.Once
)

func getLogger() *Logger {
	once.Do(func() {
		logPath := os.Getenv("LOG_PATH")
		if logPath == "" {
			logPath = "logs/log-{time}.log"
		}

		// Replace {time}
		logPath = strings.Replace(logPath, "{time}", time.Now().Format("06-01-02"), 1)

		// Ensure directory exists
		if err := os.MkdirAll(getDir(logPath), 0755); err != nil {
			panic(fmt.Sprintf("failed to create log dir: %v", err))
		}

		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Sprintf("failed to open log file: %v", err))
		}

		instance = &Logger{
			file:   file,
			logger: log.New(file, "", log.LstdFlags|log.Lshortfile),
		}
	})
	return instance
}

func getDir(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx == -1 {
		return "."
	}
	return path[:idx]
}

func Info(messages ...interface{}) {
	getLogger().log("INFO", messages...)
}

func Error(messages ...interface{}) {
	getLogger().log("ERROR", messages...)
}

func Debug(messages ...interface{}) {
	getLogger().log("DEBUG", messages...)
}

func (l *Logger) log(level string, messages ...interface{}) {
	prefix := fmt.Sprintf("[%s]", level)
	msg := fmt.Sprint(messages...)
	// file
	err := l.logger.Output(3, fmt.Sprintf("%s %s", prefix, msg))
	if err != nil {
		return
	}
	// console
	fmt.Printf("%s %s\n", prefix, msg)
}
