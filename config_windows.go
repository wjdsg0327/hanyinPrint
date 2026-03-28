//go:build windows

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadAppConfigFromJSON 从 json 配置文件读取服务配置。
// filePath 支持相对路径（相对当前工作目录）。
func LoadAppConfigFromJSON(filePath string) (AppConfig, error) {
	if strings.TrimSpace(filePath) == "" {
		return AppConfig{}, errors.New("配置文件路径为空")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return AppConfig{}, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AppConfig{}, fmt.Errorf("解析配置文件失败: %w", err)
	}

	if strings.TrimSpace(cfg.HTTPAddr) == "" {
		cfg.HTTPAddr = ":8080"
	}

	cfg.Printer = normalizePrintConfig(cfg.Printer)

	if strings.TrimSpace(cfg.Printer.SDKPath) != "" && !filepath.IsAbs(cfg.Printer.SDKPath) {
		baseDir := filepath.Dir(filePath)
		cfg.Printer.SDKPath = filepath.Clean(filepath.Join(baseDir, cfg.Printer.SDKPath))
	}

	return cfg, nil
}
