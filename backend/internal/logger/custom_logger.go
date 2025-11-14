package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type CustomLogger struct {
	curDay           int
	prefix           string
	wg               sync.WaitGroup
	ctx              context.Context
	cancel           context.CancelFunc
	consoleOutLogger *log.Logger
	consoleErrLogger *log.Logger
	fileLogger       *log.Logger
	logFile          *os.File
}

var customLogger *CustomLogger
var maxAgeDays int

func NewCustomLogger(prefix string, maxAge int) (*CustomLogger, error) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := &CustomLogger{
		prefix: prefix,
		ctx:    ctx,
		cancel: cancel,
	}

	maxAgeDays = maxAge

	if err := logger.setLogFileAndLogger(prefix); err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *CustomLogger) setLogFileAndLogger(prefix string) error {
	now := time.Now()
	filepath := getFilepath(now)
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	//consoleOutLog := log.New(os.Stdout, prefix, log.Ldate|log.Ltime)
	consoleErrLog := log.New(os.Stderr, prefix, log.Ldate|log.Ltime)
	fileLog := log.New(file, prefix, log.Ldate|log.Ltime)

	l.logFile = file
	//l.consoleOutLogger = consoleOutLog
	l.consoleErrLogger = consoleErrLog
	l.fileLogger = fileLog
	l.curDay = now.Day()
	return nil
}

func getFilepath(now time.Time) string {
	dir := filepath.Base("log")
	fileName := fmt.Sprintf("%d_%02d_%02d", now.Year(), now.Month(), now.Day())
	return filepath.Join(dir, fileName)
}

func SetLogger(logger *CustomLogger) {
	if customLogger != nil {
		return
	}
	customLogger = logger
}

func Println(v ...any) {
	//customLogger.consoleOutLogger.Println(v...)
	customLogger.consoleErrLogger.Println(v...)
	customLogger.fileLogger.Println(v...)
}

func Printf(format string, v ...any) {
	//customLogger.consoleOutLogger.Println(v...)
	customLogger.consoleErrLogger.Printf(format, v...)
	customLogger.fileLogger.Printf(format, v...)
}

func Close() error {
	return customLogger.logFile.Close()
}

func StartCleaning() {
	customLogger.wg.Add(1)
	oneMinuteTicker := time.NewTicker(1 * time.Minute)
	defer func() {
		oneMinuteTicker.Stop()
		customLogger.wg.Done()
	}()

	cleanOldLogs()

	for {
		select {
		case <-customLogger.ctx.Done():
			Println("[Log Cleanup] Log cleanup goroutine terminated")
			return
		case <-oneMinuteTicker.C:
			now := time.Now()
			if customLogger.curDay != now.Day() {
				Close()
				customLogger.setLogFileAndLogger(customLogger.prefix)
				cleanOldLogs()
				customLogger.curDay = now.Day()
			}

		}
	}
}

func cleanOldLogs() {
	cutoffTime := time.Now().AddDate(0, 0, maxAgeDays)

	dir := filepath.Base("log")

	files, err := os.ReadDir(dir)
	if err != nil {
		Printf("[Log Cleanup] Failed to read log directory:\n\t %v", err)
		return
	}

	var deletedCount int

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			Printf("[Log Cleanup] Failed to get file information:\n\t %v", err)
			continue
		}

		if info.ModTime().Before(cutoffTime) {

			filePath := filepath.Join(dir, file.Name())
			Printf("[Log Cleanup] Old log files deleted: %s (Deleted date: %s)",
				filePath, info.ModTime().Format("2006-01-02"))

			if err := os.Remove(filePath); err != nil {
				Printf("[Log Cleanup] Failed to delete file:\n\t %v", err)
			} else {
				deletedCount++
			}
		}
	}

	Printf("[Log Cleanup] Log cleanup completed: %d files deleted in total", deletedCount)
}

func Shutdown() {
	customLogger.cancel()
	customLogger.wg.Wait()
	Close()
}
