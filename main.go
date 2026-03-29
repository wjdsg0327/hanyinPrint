//go:build windows

package main

import (
	"flag"
	"fmt"
	"log"
)

// main 仅负责启动 HTTP 服务，具体打印逻辑拆分到其他文件。
func main() {
	configPath := flag.String("config", "config.json", "配置文件路径（json）")
	flag.Parse()

	appCfg, err := LoadAppConfigFromJSON(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	printer := NewPrinterService(appCfg.Printer)

	err1 := Connect("ws://127.0.0.1:30000/api/ws/"+appCfg.TenantCode, printer)
	if err1 != nil {
		fmt.Println("adsdadsd", err1)
		return
	}

	if err := StartHTTPServerWithPrinter(printer, appCfg.HTTPAddr); err != nil {
		log.Fatal(err)
	}

}
