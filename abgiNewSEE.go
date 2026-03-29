package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
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
		if err := recover(); err != nil {
			err = fmt.Errorf("捕获异常:%v", err)
			fmt.Println(err)
			return
		}
	}()

	if abgiClient != nil {
		return fmt.Errorf("已经在线，请勿重复上线")
	}

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("连接失败")
	}

	abgiClient = &AbgiClient{
		Conn:    conn,
		URL:     url,
		Headers: nil,
		Printer: printer,
	}

	go abgiClient.listen()
	return nil
}

func (c *AbgiClient) listen() {
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("接收消息失败:%v\n", err)
			go c.reconnectLoop()
			return
		}

		var zhyhFields []ZhyhField
		err = json.Unmarshal(msg, &zhyhFields)
		if err != nil {
			fmt.Printf("解析消息失败:%v\n", err)
		}
		if len(zhyhFields) > 0 && c.Printer != nil {
			if _, err := c.Printer.PrintFields(zhyhFields, 1); err != nil {
				fmt.Printf("打印失败:%v\n", err)
			}
		}

		fmt.Println("收到消息", zhyhFields)
	}
}

func (c *AbgiClient) reconnectLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Conn != nil {
		c.Conn.Close()
		c.Conn = nil
	}

	for {
		dialer := websocket.DefaultDialer
		conn, _, err := dialer.Dial(c.URL, c.Headers)
		if err == nil {
			c.Conn = conn
			go c.listen()
			return
		}
		fmt.Printf("重新连接失败:%v, 1分钟后再试\n", err)
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
