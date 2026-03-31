//go:build windows

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// ensureTenantCode 引导用户在浏览器中填写 tenant_code，并保存到配置文件。
func ensureTenantCode(cfg *AppConfig, configPath string) error {
	if strings.TrimSpace(cfg.TenantCode) != "" {
		return nil
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("启动配置引导失败: %w", err)
	}

	addr := "http://" + ln.Addr().String()
	done := make(chan error, 1)

	mux := http.NewServeMux()
	successHTML := `<!doctype html><html><head><meta charset="UTF-8"><style>body{margin:0;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif;background:#f5f7fa;display:flex;align-items:center;justify-content:center;height:100vh;color:#1f2d3d;} .card{width:360px;background:#fff;border-radius:12px;box-shadow:0 10px 30px rgba(0,0,0,0.08);padding:28px;} h2{margin:0 0 12px;font-weight:700;font-size:20px;color:#1f2d3d;} p{margin:0 0 8px;color:#58677c;} .btn{display:inline-block;margin-top:12px;padding:10px 16px;border:none;border-radius:8px;background:#2d8cf0;color:#fff;font-weight:600;cursor:pointer;} .hint{color:#19be6b;font-weight:600;margin-top:8px;}</style></head><body><div class="card"><h2>设备编码已保存</h2><p class="hint">可以关闭此窗口并返回程序。</p><button class="btn" onclick="window.close()">关闭</button></div></body></html>`
	formHTML := func(errMsg string) string {
		errBlock := ""
		if errMsg != "" {
			errBlock = fmt.Sprintf(`<div class="alert">%s</div>`, errMsg)
		}
		return fmt.Sprintf(`<!doctype html><html><head><meta charset="UTF-8"><style>body{margin:0;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif;background:#f5f7fa;display:flex;align-items:center;justify-content:center;height:100vh;color:#1f2d3d;} .card{width:360px;background:#fff;border-radius:12px;box-shadow:0 10px 30px rgba(0,0,0,0.08);padding:28px;} h2{margin:0 0 12px;font-weight:700;font-size:20px;color:#1f2d3d;} p{margin:0 0 10px;color:#58677c;} form{display:flex;flex-direction:column;gap:12px;} label{font-weight:600;color:#1f2d3d;} input{padding:10px 12px;border:1px solid #d8e3f0;border-radius:8px;font-size:14px;outline:none;transition:border-color .2s,box-shadow .2s;} input:focus{border-color:#2d8cf0;box-shadow:0 0 0 3px rgba(45,140,240,0.12);} .btn{padding:11px 14px;border:none;border-radius:8px;background:#2d8cf0;color:#fff;font-weight:700;font-size:14px;cursor:pointer;transition:transform .1s ease,box-shadow .2s;} .btn:hover{box-shadow:0 8px 20px rgba(45,140,240,0.25);transform:translateY(-1px);} .alert{padding:10px 12px;border-radius:8px;background:#ffecec;color:#d93025;font-weight:600;margin-bottom:4px;border:1px solid #ffd2d2;}</style></head><body><div class="card"><h2>请输入设备编码</h2><p>用于连接云端服务，格式：用户编号-设备编号</p>%s<form method="POST" action="/submit"><label for="tenant_code">tenant_code</label><input id="tenant_code" name="tenant_code" placeholder="例：tenant-001" autofocus required><button class="btn" type="submit">保存</button></form></div></body></html>`, errBlock)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(formHTML("")))
	})
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "参数错误", http.StatusBadRequest)
			return
		}
		tenant := strings.TrimSpace(r.Form.Get("tenant_code"))
		if tenant == "" {
			_, _ = w.Write([]byte(formHTML("设备编码不能为空")))
			return
		}

		cfg.TenantCode = tenant
		if err := SaveAppConfigToJSON(*cfg, configPath); err != nil {
			http.Error(w, "保存失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write([]byte(successHTML))
		done <- nil
	})

	srv := &http.Server{Handler: mux}
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			done <- err
		}
	}()

	if err := openBrowser(addr); err != nil {
		return err
	}

	select {
	case err := <-done:
		_ = srv.Shutdown(context.Background())
		return err
	case <-time.After(10 * time.Minute):
		_ = srv.Shutdown(context.Background())
		return errors.New("等待租户编码输入超时")
	}
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		return fmt.Errorf("不支持的系统: %s", runtime.GOOS)
	}
}
