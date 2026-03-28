//go:build windows

package main

import "strings"

func normalizePrintConfig(cfg PrintConfig) PrintConfig {
	if strings.TrimSpace(cfg.SDKPath) == "" {
		cfg.SDKPath = DefaultDLLPath()
	}
	if strings.TrimSpace(cfg.Model) == "" {
		cfg.Model = "N41"
	}
	if strings.TrimSpace(cfg.Port) == "" {
		cfg.Port = "USB"
	}
	if cfg.Options.LabelWidthMM == 0 {
		cfg.Options.LabelWidthMM = 60
	}
	if cfg.Options.LabelHeightMM == 0 {
		cfg.Options.LabelHeightMM = 40
	}
	if cfg.Options.Speed == 0 {
		cfg.Options.Speed = 2
	}
	if cfg.Options.Density == 0 {
		cfg.Options.Density = 6
	}
	if cfg.Options.Type == 0 {
		cfg.Options.Type = 1
	}
	if cfg.Options.GapMM == 0 {
		cfg.Options.GapMM = 2
	}
	if cfg.Options.DPI == 0 {
		cfg.Options.DPI = 203
	}
	if cfg.Options.MarginLeftMM == 0 {
		cfg.Options.MarginLeftMM = 2
	}
	if cfg.Options.MarginTopMM == 0 {
		cfg.Options.MarginTopMM = 2
	}
	return cfg
}

func readPrinterInfo(sdk *tsplSDK, handle uintptr, usedModel string) (PrinterInfo, error) {
	status, err := sdk.getPrinterStatus(handle)
	if err != nil {
		return PrinterInfo{}, err
	}

	info := PrinterInfo{
		UsedModel: usedModel,
		Status:    status,
	}
	if fw, err := sdk.getFirmwareVersion(handle); err == nil {
		info.Firmware = strings.TrimSpace(fw)
	}
	if sn, err := sdk.getSN(handle); err == nil {
		info.SN = strings.TrimSpace(sn)
	}
	return info, nil
}
