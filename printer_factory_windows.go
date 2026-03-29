//go:build windows

package main

import (
	"errors"
	"strings"

	"go.uber.org/zap"
)

// tryCreatePrinter 用多个候选 model 尝试创建打印机句柄。
// 有些 SDK 版本/配置下，必须用特定字符串才能识别型号。
func tryCreatePrinter(sdk *tsplSDK, model string) (uintptr, string, error) {
	candidates := []string{
		strings.TrimSpace(model),
		strings.ToUpper(strings.TrimSpace(model)),
		"N41BT",
		"N41",
		"ANY",
	}
	seen := map[string]struct{}{}
	var lastErr error
	for _, m := range candidates {
		if strings.TrimSpace(m) == "" {
			continue
		}
		if _, ok := seen[m]; ok {
			continue
		}
		seen[m] = struct{}{}
		L().Info("trying printer model", zap.String("model", m))
		h, err := sdk.printerCreator(m)
		if err == nil {
			L().Info("printer model selected", zap.String("model", m))
			return h, m, nil
		}
		L().Warn("printer model failed", zap.String("model", m), zap.Error(err))
		lastErr = err
	}
	if lastErr == nil {
		lastErr = errors.New("PrinterCreator 失败：没有可用的 model")
	}
	L().Error("all printer models failed", zap.Error(lastErr))
	return 0, "", lastErr
}
