//go:build windows

package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type PrinterService struct {
	cfg PrintConfig
	mu  sync.Mutex
}

func NewPrinterService(cfg PrintConfig) *PrinterService {
	return &PrinterService{cfg: normalizePrintConfig(cfg)}
}

func (s *PrinterService) PrintFields(fields []ZhyhField, copies int32) (PrinterInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return PrintFields(s.cfg, fields, copies)
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

	g := gin.Default()
	g.POST("/print", func(c *gin.Context) {
		var fields []ZhyhField
		if err := c.ShouldBindJSON(&fields); err != nil {
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
