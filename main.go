//go:build windows

package main

import (
	"flag"
	"os"

	"go.uber.org/zap"
)

// main 仅负责启动 HTTP 服务，具体打印逻辑拆分到其他文件。
func main() {
	configPath := flag.String("config", "config.json", "配置文件路径（json）")
	flag.Parse()

	appCfg, err := LoadAppConfigFromJSON(*configPath)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	if err := InitLogger(appCfg.Log); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	defer func() {
		_ = SyncLogger()
	}()

	L().Info("service starting",
		zap.String("http_addr", appCfg.HTTPAddr),
		zap.String("tenant_code", appCfg.TenantCode),
		zap.String("printer_model", appCfg.Printer.Model),
		zap.String("printer_port", appCfg.Printer.Port),
		zap.String("log_file", appCfg.Log.FilePath),
	)

	printer := NewPrinterService(appCfg.Printer)

	if err := Connect("ws://127.0.0.1:30000/api/ws/"+appCfg.TenantCode, printer); err != nil {
		L().Error("websocket connect failed", zap.Error(err), zap.String("tenant_code", appCfg.TenantCode))
		return
	}

	if err := StartHTTPServerWithPrinter(printer, appCfg.HTTPAddr); err != nil {
		L().Fatal("http server exited", zap.Error(err), zap.String("http_addr", appCfg.HTTPAddr))
	}
}
