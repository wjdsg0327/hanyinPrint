//go:build windows

package main

import (
	"flag"
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

	if err := StartHTTPServer(appCfg.Printer, appCfg.HTTPAddr); err != nil {
		log.Fatal(err)
	}
}
