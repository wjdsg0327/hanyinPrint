//go:build windows

package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerMu  sync.RWMutex
	logger    = zap.NewNop()
	logWriter *dailyFileWriter
)

type dailyFileWriter struct {
	mu          sync.Mutex
	basePath    string
	currentDate string
	file        *os.File
}

func InitLogger(cfg LogConfig) error {
	level := zap.InfoLevel
	if err := level.UnmarshalText([]byte(strings.TrimSpace(cfg.Level))); err != nil {
		return err
	}

	writer, err := newDailyFileWriter(cfg.FilePath)
	if err != nil {
		return err
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = formatLogTime

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		level,
	)

	newLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	loggerMu.Lock()
	oldLogger := logger
	oldWriter := logWriter
	logger = newLogger
	logWriter = writer
	loggerMu.Unlock()

	_ = oldLogger.Sync()
	if oldWriter != nil {
		_ = oldWriter.Close()
	}
	return nil
}

func newDailyFileWriter(basePath string) (*dailyFileWriter, error) {
	w := &dailyFileWriter{basePath: basePath}
	if err := w.rotateIfNeeded(time.Now()); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *dailyFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.rotateIfNeeded(time.Now()); err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

func (w *dailyFileWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return nil
	}
	return w.file.Sync()
}

func (w *dailyFileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return nil
	}
	oldFile := w.file
	w.file = nil
	w.currentDate = ""
	return oldFile.Close()
}

func (w *dailyFileWriter) rotateIfNeeded(now time.Time) error {
	date := now.Local().Format("2006-01-02")
	if w.file != nil && w.currentDate == date {
		return nil
	}

	path := dailyLogFilePath(w.basePath, date)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	oldFile := w.file
	w.file = file
	w.currentDate = date

	if oldFile != nil {
		_ = oldFile.Sync()
		_ = oldFile.Close()
	}
	return nil
}

func dailyLogFilePath(basePath, date string) string {
	ext := filepath.Ext(basePath)
	name := strings.TrimSuffix(filepath.Base(basePath), ext)
	if strings.TrimSpace(name) == "" {
		name = "log"
	}

	dir := filepath.Dir(basePath)
	if ext == "" {
		return filepath.Join(dir, name+"-"+date)
	}
	return filepath.Join(dir, name+"-"+date+ext)
}

func formatLogTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Local().Format("2006-01-02 15:04:05"))
}

func L() *zap.Logger {
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	return logger
}

func SyncLogger() error {
	loggerMu.RLock()
	current := logger
	loggerMu.RUnlock()
	return current.Sync()
}
