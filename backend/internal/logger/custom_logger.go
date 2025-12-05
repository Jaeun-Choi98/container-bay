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
	curDay       int
	prefix       string
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	stdoutLogger *log.Logger
	stdoutFile   *os.File
	stderrFile   *os.File
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
	errfilepath := getFilepath("stderr", now)
	errfile, err := os.OpenFile(errfilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	outfilepath := getFilepath("stdout", now)
	outfile, err := os.OpenFile(outfilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	os.Stderr = errfile
	outLogger := log.New(outfile, prefix, log.Ldate|log.Ltime)

	l.stdoutFile = outfile
	l.stderrFile = errfile
	l.stdoutLogger = outLogger
	l.curDay = now.Day()
	return nil
}

func getFilepath(name string, now time.Time) string {
	dir := filepath.Base("log")
	fileName := fmt.Sprintf("%s_%d_%02d_%02d", name, now.Year(), now.Month(), now.Day())
	return filepath.Join(dir, fileName)
}

func SetLogger(logger *CustomLogger) {
	if customLogger != nil {
		return
	}
	customLogger = logger
}

func Println(v ...any) {
	log.Println(v...)
	customLogger.stdoutLogger.Println(v...)
}

func Printf(format string, v ...any) {
	log.Printf(format, v...)
	customLogger.stdoutLogger.Printf(format, v...)
}

func Close() (rsterr error) {
	if err := customLogger.stderrFile.Close(); err != nil {
		rsterr = fmt.Errorf("failed to close stderrFile: %v", err)
	}
	if err := customLogger.stdoutFile.Close(); err != nil {
		rsterr = fmt.Errorf("failed to close stdoutFile: %v, %v", err, rsterr)
	}
	return rsterr
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
