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

func (s *PrinterService) GetStatus() (PrinterInfo, error) {
	cfg := normalizePrintConfig(s.cfg)
	start := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	sdk, err := newTSPLSDK(cfg.SDKPath)
	if err != nil {
		L().Warn("load sdk failed for status probe",
			zap.Error(err),
			zap.String("sdk_path", cfg.SDKPath),
			zap.Duration("latency", time.Since(start)),
		)
		return PrinterInfo{}, err
	}
	sdk.sdkInit()
	defer sdk.sdkDeInit()

	handle, usedModel, err := tryCreatePrinter(sdk, cfg.Model)
	if err != nil {
		L().Warn("create printer handle failed for status probe",
			zap.Error(err),
			zap.String("model", cfg.Model),
			zap.Duration("latency", time.Since(start)),
		)
		return PrinterInfo{}, err
	}
	defer func() {
		if err := sdk.printerDestroy(handle); err != nil {
			L().Warn("destroy printer handle failed after status probe",
				zap.Error(err),
				zap.String("used_model", usedModel),
			)
		}
	}()

	if err := sdk.portOpen(handle, cfg.Port); err != nil {
		L().Warn("open printer port failed for status probe",
			zap.Error(err),
			zap.String("port", cfg.Port),
			zap.String("used_model", usedModel),
			zap.Duration("latency", time.Since(start)),
		)
		return PrinterInfo{}, err
	}
	defer func() {
		if err := sdk.portClose(handle); err != nil {
			L().Warn("close printer port failed after status probe",
				zap.Error(err),
				zap.String("port", cfg.Port),
				zap.String("used_model", usedModel),
			)
		}
	}()

	info, err := readPrinterInfo(sdk, handle, usedModel)
	if err != nil {
		L().Warn("read printer info failed for status probe",
			zap.Error(err),
			zap.String("used_model", usedModel),
			zap.Duration("latency", time.Since(start)),
		)
		return PrinterInfo{}, err
	}

	L().Info("printer status probe succeeded",
		zap.String("port", cfg.Port),
		zap.String("used_model", info.UsedModel),
		zap.Int32("printer_status", info.Status),
		zap.String("firmware", info.Firmware),
		zap.String("sn", info.SN),
		zap.Duration("latency", time.Since(start)),
	)
	return info, nil
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
//   - GET /printer/status
//   - POST /print?copies=1
//     Body: []ZhyhField
func StartHTTPServer(cfg PrintConfig, appConfig AppConfig) error {
	return StartHTTPServerWithPrinter(NewPrinterService(cfg), appConfig)
}

func StartHTTPServerWithPrinter(printer *PrinterService, appConfig AppConfig) error {
	if printer == nil {
		return errors.New("printer service is nil")
	}

	g := gin.New()
	g.Use(httpRequestLogger(), httpRecoveryLogger())
	g.GET("/printer/status", func(c *gin.Context) {
		info, err := printer.GetStatus()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				//"ok":        true,
				"connected": false,
				"error":     err.Error(),
				"设备编码":      appConfig.TenantCode,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			//"ok":      true,
			"打印机是否连接": true,
			"info":    info,
			"设备编码":    appConfig.TenantCode,
		})
	})
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

	return g.Run(appConfig.HTTPAddr)
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
