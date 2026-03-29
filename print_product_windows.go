//go:build windows

package main

import (
	"errors"
	"strings"

	"go.uber.org/zap"
)

// PrintProductLabel 打印常见“商品标签”（固定字段版）。
//
// 如果你的前端传的是动态字段数组（ZhyhField），请使用 PrintFields。
func PrintProductLabel(cfg PrintConfig, p ProductLabel, copies int32) (PrinterInfo, error) {
	cfg = normalizePrintConfig(cfg)
	if copies <= 0 {
		copies = 1
	}

	L().Info("starting product label print",
		zap.String("name", p.Name),
		zap.String("sku", p.SKU),
		zap.Int32("copies", copies),
		zap.String("sdk_path", cfg.SDKPath),
		zap.String("model", cfg.Model),
		zap.String("port", cfg.Port),
	)

	sdk, err := newTSPLSDK(cfg.SDKPath)
	if err != nil {
		L().Error("load sdk failed for product label", zap.Error(err), zap.String("sdk_path", cfg.SDKPath))
		return PrinterInfo{}, err
	}
	sdk.sdkInit()
	defer sdk.sdkDeInit()

	handle, usedModel, err := tryCreatePrinter(sdk, cfg.Model)
	if err != nil {
		L().Error("create printer handle failed for product label", zap.Error(err), zap.String("model", cfg.Model))
		return PrinterInfo{}, err
	}
	defer func() {
		if err := sdk.printerDestroy(handle); err != nil {
			L().Warn("destroy printer handle failed", zap.Error(err), zap.String("used_model", usedModel))
		}
	}()

	if err := sdk.portOpen(handle, cfg.Port); err != nil {
		L().Error("open printer port failed for product label", zap.Error(err), zap.String("port", cfg.Port), zap.String("used_model", usedModel))
		return PrinterInfo{}, err
	}
	defer func() {
		if err := sdk.portClose(handle); err != nil {
			L().Warn("close printer port failed", zap.Error(err), zap.String("port", cfg.Port), zap.String("used_model", usedModel))
		}
	}()

	info, err := readPrinterInfo(sdk, handle, usedModel)
	if err != nil {
		L().Error("read printer info failed for product label", zap.Error(err), zap.String("used_model", usedModel))
		return PrinterInfo{}, err
	}

	if err := printProductLabel(sdk, handle, cfg.Options, p, copies); err != nil {
		L().Error("render product label failed", zap.Error(err), zap.String("used_model", usedModel), zap.String("sku", p.SKU))
		return info, err
	}

	L().Info("product label print prepared", zap.String("used_model", usedModel), zap.String("sku", p.SKU), zap.Int32("copies", copies))
	return info, nil
}

func printProductLabel(sdk *tsplSDK, handle uintptr, opt LabelOptions, p ProductLabel, copies int32) error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("商品名称不能为空")
	}
	if strings.TrimSpace(p.SKU) == "" {
		return errors.New("SKU/货号不能为空")
	}
	if strings.TrimSpace(p.Barcode) == "" {
		p.Barcode = p.SKU
	}
	if strings.TrimSpace(p.QR) == "" {
		p.QR = p.Barcode
	}

	if err := sdk.tsplClearBuffer(handle); err != nil {
		return err
	}
	if err := sdk.tsplSetup(handle, opt.LabelWidthMM, opt.LabelHeightMM, opt.Speed, opt.Density, opt.Type, opt.GapMM, opt.OffsetMM); err != nil {
		return err
	}

	xLeft := MmToDots(opt.MarginLeftMM, opt.DPI)
	yTop := MmToDots(opt.MarginTopMM, opt.DPI)

	nameScale := int32(2)
	if len([]rune(p.Name)) > 12 {
		nameScale = 1
	}
	if err := sdk.tsplTextCompatibleGBK(handle, xLeft, yTop, 9, 0, nameScale, nameScale, p.Name); err != nil {
		return err
	}

	y := yTop + MmToDots(8, opt.DPI)
	if err := sdk.tsplTextCompatibleGBK(handle, xLeft, y, 9, 0, 1, 1, "SKU:"+p.SKU); err != nil {
		return err
	}

	if strings.TrimSpace(p.Spec) != "" {
		y += MmToDots(5, opt.DPI)
		if err := sdk.tsplTextCompatibleGBK(handle, xLeft, y, 9, 0, 1, 1, "规格:"+p.Spec); err != nil {
			return err
		}
	}

	if strings.TrimSpace(p.Price) != "" {
		y += MmToDots(5, opt.DPI)
		if err := sdk.tsplTextCompatibleGBK(handle, xLeft, y, 9, 0, 2, 2, "¥"+p.Price); err != nil {
			return err
		}
	}

	barcodeY := MmToDots(opt.LabelHeightMM-18, opt.DPI)
	if barcodeY < y+MmToDots(8, opt.DPI) {
		barcodeY = y + MmToDots(8, opt.DPI)
	}
	if strings.TrimSpace(p.Barcode) != "" {
		if err := sdk.tsplBarCode(handle, xLeft, barcodeY, 0, MmToDots(10, opt.DPI), 1, 0, 2, 2, p.Barcode); err != nil {
			return err
		}
	}

	if strings.TrimSpace(p.QR) != "" {
		qrX := MmToDots(opt.LabelWidthMM-18, opt.DPI)
		qrY := MmToDots(12, opt.DPI)
		if qrX < xLeft+MmToDots(25, opt.DPI) {
			qrX = xLeft + MmToDots(25, opt.DPI)
		}
		if err := sdk.tsplQrCode(handle, qrX, qrY, 3, 4, 0, 0, 0, 7, "\""+p.QR+"\""); err != nil {
			return err
		}
	}

	if err := sdk.tsplPrint(handle, 1, copies); err != nil {
		return err
	}
	return nil
}
