//go:build windows

package main

import (
	"math"
	"os"
	"path/filepath"
)

// DefaultDLLPath 返回默认的 TSPL_SDK.dll 位置（相对项目根目录）。
func DefaultDLLPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Join(wd, "Resources", "Libs", "x64", "Ansi", "TSPL_SDK.dll")
}

// MmToDots 将 mm 换算为打印点（dots），dpi 常见为 203/300。
func MmToDots(mm int32, dpi int32) int32 {
	if mm <= 0 || dpi <= 0 {
		return 0
	}
	return int32(math.Round(float64(mm) * float64(dpi) / 25.4))
}
