//go:build windows

package main

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// StartHTTPServer 启动 Gin 服务。
//
//   - POST /print?copies=1
//     Body: []ZhyhField
func StartHTTPServer(cfg PrintConfig, addr string) error {
	var printMu sync.Mutex

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

		printMu.Lock()
		info, err := PrintFields(cfg, fields, copies)
		printMu.Unlock()

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
