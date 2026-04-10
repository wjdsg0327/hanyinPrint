package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type AbgiClient struct {
	Conn    *websocket.Conn
	URL     string
	Headers http.Header
	Printer *PrinterService
	mu      sync.Mutex
}

var abgiClient *AbgiClient

func Connect(url string, printer *PrinterService) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("捕获异常:%v", recovered)
			L().Error("websocket connect panic", zap.Error(err), zap.String("url", url))
		}
	}()

	if abgiClient != nil {

		return fmt.Errorf("已经在线，请勿重复上线")
	}

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		L().Error("连接失败")
		return fmt.Errorf("连接失败: %w", err)
	}

	abgiClient = &AbgiClient{
		Conn:    conn,
		URL:     url,
		Headers: nil,
		Printer: printer,
	}

	//L().Info("websocket connected", zap.String("url", url))
	go abgiClient.listen()
	return nil
}

func (c *AbgiClient) listen() {
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			L().Error("websocket read failed", zap.Error(err), zap.String("url", c.URL))
			go c.reconnectLoop()
			return
		}

		var zhyhFields []ZhyhField
		if err := json.Unmarshal(msg, &zhyhFields); err != nil {
			L().Error("websocket payload parse failed", zap.Error(err), zap.String("url", c.URL), zap.ByteString("payload", msg))
			continue
		}

		L().Info("websocket payload received", zap.String("url", c.URL), zap.Int("field_count", len(zhyhFields)))

		if len(zhyhFields) > 0 && c.Printer != nil {
			if _, err := c.Printer.PrintFields(zhyhFields, 1); err != nil {
				L().Error("websocket print failed", zap.Error(err), zap.String("url", c.URL), zap.Int("field_count", len(zhyhFields)))
			}
		}
	}
}

func (c *AbgiClient) reconnectLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Conn != nil {
		_ = c.Conn.Close()
		c.Conn = nil
	}

	for {
		dialer := websocket.DefaultDialer
		conn, _, err := dialer.Dial(c.URL, c.Headers)
		if err == nil {
			c.Conn = conn
			L().Info("websocket reconnected", zap.String("url", c.URL))
			go c.listen()
			return
		}
		L().Warn("websocket reconnect failed", zap.Error(err), zap.String("url", c.URL), zap.Duration("retry_after", time.Minute))
		time.Sleep(time.Minute)
	}
}

func (c *AbgiClient) Send(msg string) error {
	if abgiClient == nil {
		return fmt.Errorf("WebSocket 未连接")
	}
	abgiClient.mu.Lock()
	defer abgiClient.mu.Unlock()

	return abgiClient.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
}
