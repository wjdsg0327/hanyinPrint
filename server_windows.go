//go:build windows

package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PrinterService struct {
	cfg PrintConfig
	mu  sync.Mutex
}

func NewPrinterService(cfg PrintConfig) *PrinterService {
	return &PrinterService{cfg: normalizePrintConfig(cfg)}
}

func (s *PrinterService) PrintFields(fields []ZhyhField, copies int32) (PrinterInfo, error) {
	effectiveCopies := copies
	if effectiveCopies <= 0 {
		effectiveCopies = 1
	}

	start := time.Now()
	s.mu.Lock()
	info, err := PrintFields(s.cfg, fields, effectiveCopies)
	s.mu.Unlock()

	logFields := []zap.Field{
		zap.Int("field_count", len(fields)),
		zap.Int32("copies", effectiveCopies),
		zap.Duration("latency", time.Since(start)),
		zap.String("used_model", info.UsedModel),
		zap.Int32("printer_status", info.Status),
		zap.String("firmware", info.Firmware),
		zap.String("sn", info.SN),
	}
	if err != nil {
		logFields = append(logFields, zap.Error(err))
		L().Error("print failed", logFields...)
		return info, err
	}

	L().Info("print succeeded", logFields...)
	return info, nil
}

// StartHTTPServer 启动 Gin 服务。
//
//   - POST /print?copies=1
//     Body: []ZhyhField
func StartHTTPServer(cfg PrintConfig, addr string) error {
	return StartHTTPServerWithPrinter(NewPrinterService(cfg), addr)
}

func StartHTTPServerWithPrinter(printer *PrinterService, addr string) error {
	if printer == nil {
		return errors.New("printer service is nil")
	}

	g := gin.New()
	g.Use(httpRequestLogger(), httpRecoveryLogger())
	g.POST("/print", func(c *gin.Context) {
		var fields []ZhyhField
		if err := c.ShouldBindJSON(&fields); err != nil {
			L().Warn("http print bind failed",
				zap.Error(err),
				zap.String("client_ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusBadRequest, gin.H{
				"ok":    false,
				"error": err.Error(),
			})
			return
		}

		copies := int32(1)
		if v := strings.TrimSpace(c.Query("copies")); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				copies = int32(n)
			}
		}

		info, err := printer.PrintFields(fields, copies)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"ok":    false,
				"error": err.Error(),
				"info":  info,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"info": info,
		})
	})

	return g.Run(addr)
}

func httpRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
		}
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		L().Info("http request", fields...)
	}
}

func httpRecoveryLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				path := c.FullPath()
				if path == "" {
					path = c.Request.URL.Path
				}
				L().Error("http panic",
					zap.Any("panic", recovered),
					zap.String("method", c.Request.Method),
					zap.String("path", path),
					zap.String("client_ip", c.ClientIP()),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
