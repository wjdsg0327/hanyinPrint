//go:build windows

package main

import (
	"errors"
	"strings"
)

// PrintFields 打印动态字段数组（由前端传入）。
//
// 字段规则（可按需扩展）：
// - zhyh_type = text:   打印 "key:value"
// - zhyh_type = barcode:打印 Code128 条码（value 作为内容）
// - zhyh_type = qrcode: 打印二维码（value 作为内容）
func PrintFields(cfg PrintConfig, fields []ZhyhField, copies int32) (PrinterInfo, error) {
	cfg = normalizePrintConfig(cfg)
	if copies <= 0 {
		copies = 1
	}

	sdk, err := newTSPLSDK(cfg.SDKPath)
	if err != nil {
		return PrinterInfo{}, err
	}
	sdk.sdkInit()
	defer sdk.sdkDeInit()

	handle, usedModel, err := tryCreatePrinter(sdk, cfg.Model)
	if err != nil {
		return PrinterInfo{}, err
	}
	defer func() {
		_ = sdk.printerDestroy(handle)
	}()

	if err := sdk.portOpen(handle, cfg.Port); err != nil {
		return PrinterInfo{}, err
	}
	defer func() {
		_ = sdk.portClose(handle)
	}()

	info, err := readPrinterInfo(sdk, handle, usedModel)
	if err != nil {
		return PrinterInfo{}, err
	}

	if err := printFields(sdk, handle, cfg.Options, fields, copies); err != nil {
		return info, err
	}

	return info, nil
}

func printFields(sdk *tsplSDK, handle uintptr, opt LabelOptions, fields []ZhyhField, copies int32) error {
	if len(fields) == 0 {
		return errors.New("打印内容为空")
	}

	if err := sdk.tsplClearBuffer(handle); err != nil {
		return err
	}
	if err := sdk.tsplSetup(handle, opt.LabelWidthMM, opt.LabelHeightMM, opt.Speed, opt.Density, opt.Type, opt.GapMM, opt.OffsetMM); err != nil {
		return err
	}

	xLeft := MmToDots(opt.MarginLeftMM, opt.DPI)
	y := MmToDots(opt.MarginTopMM, opt.DPI)
	lineGap := MmToDots(5, opt.DPI)

	for _, f := range fields {
		typ := strings.ToLower(strings.TrimSpace(f.ZhyhType))
		key := strings.TrimSpace(f.ZhyhKey)
		val := strings.TrimSpace(f.ZhyhValue)

		switch typ {
		case "barcode", "bar", "code128":
			if val == "" {
				continue
			}
			if err := sdk.tsplBarCode(handle, xLeft, y, 0, MmToDots(10, opt.DPI), 1, 0, 2, 2, val); err != nil {
				return err
			}
			y += MmToDots(16, opt.DPI)
		case "qrcode", "qr":
			if val == "" {
				continue
			}
			qrX := MmToDots(opt.LabelWidthMM-18, opt.DPI)
			if qrX < xLeft {
				qrX = xLeft
			}
			if err := sdk.tsplQrCode(handle, qrX, y, 3, 4, 0, 0, 0, 7, "\""+val+"\""); err != nil {
				return err
			}
			y += MmToDots(18, opt.DPI)
		default:
			if key == "" && val == "" {
				continue
			}
			text := val
			if key != "" && val != "" {
				text = key + ":" + val
			} else if key != "" {
				text = key
			}
			if err := sdk.tsplTextCompatibleGBK(handle, xLeft, y, 9, 0, 1, 1, text); err != nil {
				return err
			}
			y += lineGap
		}
	}

	if err := sdk.tsplPrint(handle, 1, copies); err != nil {
		return err
	}
	return nil
}
