//go:build windows

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// tsplSDK 封装了 TSPL_SDK.dll 的关键函数。
// 这里使用 Windows DLL 调用方式，handle 为 SDK 创建的打印机句柄。
type tsplSDK struct {
	dll *windows.LazyDLL

	procSDKInit            *windows.LazyProc
	procSDKDeInit          *windows.LazyProc
	procSetConfigDir       *windows.LazyProc
	procFormatError        *windows.LazyProc
	procPrinterCreator     *windows.LazyProc
	procPrinterDestroy     *windows.LazyProc
	procPortOpen           *windows.LazyProc
	procPortClose          *windows.LazyProc
	procGetPrinterStatus   *windows.LazyProc
	procGetFirmwareVersion *windows.LazyProc
	procTSPLClearBuffer    *windows.LazyProc
	procTSPLSetup          *windows.LazyProc
	procTSPLTextCompatible *windows.LazyProc
	procTSPLBarCode        *windows.LazyProc
	procTSPLQrCode         *windows.LazyProc
	procTSPLBitMap         *windows.LazyProc
	procTSPLPrint          *windows.LazyProc
	procTSPLSelfTest       *windows.LazyProc
	procTSPLGetSN          *windows.LazyProc
}

func newTSPLSDK(dllPath string) (*tsplSDK, error) {
	if strings.TrimSpace(dllPath) == "" {
		return nil, errors.New("TSPL_SDK.dll 路径为空")
	}
	if _, err := os.Stat(dllPath); err != nil {
		return nil, fmt.Errorf("找不到 TSPL_SDK.dll: %w", err)
	}

	dll := windows.NewLazyDLL(dllPath)
	if err := dll.Load(); err != nil {
		return nil, fmt.Errorf("加载 TSPL_SDK.dll 失败: %w", err)
	}

	sdk := &tsplSDK{
		dll: dll,

		procSDKInit:            dll.NewProc("SDKInit"),
		procSDKDeInit:          dll.NewProc("SDKDeInit"),
		procSetConfigDir:       dll.NewProc("SetConfigDir"),
		procFormatError:        dll.NewProc("FormatError"),
		procPrinterCreator:     dll.NewProc("PrinterCreator"),
		procPrinterDestroy:     dll.NewProc("PrinterDestroy"),
		procPortOpen:           dll.NewProc("PortOpen"),
		procPortClose:          dll.NewProc("PortClose"),
		procGetPrinterStatus:   dll.NewProc("TSPL_GetPrinterStatus"),
		procGetFirmwareVersion: dll.NewProc("TSPL_GetFirmwareVersion"),
		procTSPLClearBuffer:    dll.NewProc("TSPL_ClearBuffer"),
		procTSPLSetup:          dll.NewProc("TSPL_Setup"),
		procTSPLTextCompatible: dll.NewProc("TSPL_TextCompatible"),
		procTSPLBarCode:        dll.NewProc("TSPL_BarCode"),
		procTSPLQrCode:         dll.NewProc("TSPL_QrCode"),
		procTSPLBitMap:         dll.NewProc("TSPL_BitMap"),
		procTSPLPrint:          dll.NewProc("TSPL_Print"),
		procTSPLSelfTest:       dll.NewProc("TSPL_SelfTest"),
		procTSPLGetSN:          dll.NewProc("TSPL_GetSN"),
	}

	if err := sdk.procFormatError.Find(); err != nil {
		return nil, fmt.Errorf("在 DLL 中找不到必要函数 FormatError: %w", err)
	}
	if err := sdk.procPrinterCreator.Find(); err != nil {
		return nil, fmt.Errorf("在 DLL 中找不到必要函数 PrinterCreator: %w", err)
	}
	if err := sdk.procPortOpen.Find(); err != nil {
		return nil, fmt.Errorf("在 DLL 中找不到必要函数 PortOpen: %w", err)
	}
	if err := sdk.procTSPLPrint.Find(); err != nil {
		return nil, fmt.Errorf("在 DLL 中找不到必要函数 TSPL_Print: %w", err)
	}

	if err := sdk.procSetConfigDir.Find(); err == nil {
		if err := sdk.setConfigDir(filepath.Dir(dllPath)); err != nil {
			return nil, err
		}
	}

	return sdk, nil
}

func (s *tsplSDK) sdkInit() {
	_, _, _ = s.procSDKInit.Call()
}

func (s *tsplSDK) sdkDeInit() {
	_, _, _ = s.procSDKDeInit.Call()
}

func (s *tsplSDK) setConfigDir(dir string) error {
	if strings.TrimSpace(dir) == "" {
		return nil
	}
	dirPtr, err := windows.BytePtrFromString(dir)
	if err != nil {
		return err
	}
	r, _, _ := s.procSetConfigDir.Call(uintptr(unsafe.Pointer(dirPtr)), uintptr(len(dir)+1))
	ret := int32(r)
	if ret != 0 {
		return fmt.Errorf("SetConfigDir 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return nil
}

func (s *tsplSDK) formatError(errorNo int32) string {
	buf := make([]byte, 512)
	r, _, _ := s.procFormatError.Call(
		uintptr(errorNo),
		uintptr(1),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(0),
		uintptr(len(buf)),
	)
	_ = int32(r)
	n := bytesCStringLen(buf)
	return string(buf[:n])
}

func (s *tsplSDK) printerCreator(model string) (uintptr, error) {
	var handle uintptr
	modelPtr, err := windows.BytePtrFromString(model)
	if err != nil {
		return 0, err
	}
	r, _, _ := s.procPrinterCreator.Call(
		uintptr(unsafe.Pointer(&handle)),
		uintptr(unsafe.Pointer(modelPtr)),
	)
	ret := int32(r)
	if ret != 0 {
		return 0, fmt.Errorf("PrinterCreator 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return handle, nil
}

func (s *tsplSDK) printerDestroy(handle uintptr) error {
	r, _, _ := s.procPrinterDestroy.Call(handle)
	ret := int32(r)
	if ret != 0 {
		return fmt.Errorf("PrinterDestroy 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return nil
}

func (s *tsplSDK) portOpen(handle uintptr, ioSettings string) error {
	ioPtr, err := windows.BytePtrFromString(ioSettings)
	if err != nil {
		return err
	}
	r, _, _ := s.procPortOpen.Call(handle, uintptr(unsafe.Pointer(ioPtr)))
	ret := int32(r)
	if ret != 0 {
		return fmt.Errorf("PortOpen 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return nil
}

func (s *tsplSDK) portClose(handle uintptr) error {
	r, _, _ := s.procPortClose.Call(handle)
	ret := int32(r)
	if ret != 0 {
		return fmt.Errorf("PortClose 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return nil
}

func (s *tsplSDK) getPrinterStatus(handle uintptr) (int32, error) {
	var status int32
	r, _, _ := s.procGetPrinterStatus.Call(handle, uintptr(unsafe.Pointer(&status)))
	ret := int32(r)
	if ret != 0 {
		return 0, fmt.Errorf("TSPL_GetPrinterStatus 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return status, nil
}

func (s *tsplSDK) getFirmwareVersion(handle uintptr) (string, error) {
	buf := make([]byte, 64)
	r, _, _ := s.procGetFirmwareVersion.Call(handle, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	ret := int32(r)
	if ret != 0 {
		return "", fmt.Errorf("TSPL_GetFirmwareVersion 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return strings.TrimRight(string(buf), "\x00"), nil
}

func (s *tsplSDK) getSN(handle uintptr) (string, error) {
	buf := make([]byte, 64)
	r, _, _ := s.procTSPLGetSN.Call(handle, uintptr(unsafe.Pointer(&buf[0])))
	ret := int32(r)
	if ret != 0 {
		return "", fmt.Errorf("TSPL_GetSN 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return strings.TrimRight(string(buf), "\x00"), nil
}

func (s *tsplSDK) tsplSelfTest(handle uintptr) error {
	r, _, _ := s.procTSPLSelfTest.Call(handle)
	ret := int32(r)
	if ret != 0 {
		return fmt.Errorf("TSPL_SelfTest 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return nil
}

func (s *tsplSDK) tsplClearBuffer(handle uintptr) error {
	r, _, _ := s.procTSPLClearBuffer.Call(handle)
	ret := int32(r)
	if ret != 0 {
		return fmt.Errorf("TSPL_ClearBuffer 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return nil
}

func (s *tsplSDK) tsplSetup(handle uintptr, labelWidth, labelHeight, speed, density, typ, gap, offset int32) error {
	r, _, _ := s.procTSPLSetup.Call(
		handle,
		uintptr(labelWidth),
		uintptr(labelHeight),
		uintptr(speed),
		uintptr(density),
		uintptr(typ),
		uintptr(gap),
		uintptr(offset),
	)
	ret := int32(r)
	if ret != 0 {
		return fmt.Errorf("TSPL_Setup 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return nil
}

func (s *tsplSDK) tsplTextCompatibleGBK(handle uintptr, xPos, yPos, font, rotation, xMul, yMul int32, text string) error {
	data, err := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(text))
	if err != nil {
		return err
	}
	data = append(data, 0)
	r, _, _ := s.procTSPLTextCompatible.Call(
		handle,
		uintptr(xPos),
		uintptr(yPos),
		uintptr(font),
		uintptr(rotation),
		uintptr(xMul),
		uintptr(yMul),
		uintptr(unsafe.Pointer(&data[0])),
	)
	runtime.KeepAlive(data)
	ret := int32(r)
	if ret != 0 {
		return fmt.Errorf("TSPL_TextCompatible 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return nil
}

func (s *tsplSDK) tsplBarCode(handle uintptr, xPos, yPos, codeType, height, readable, rotation, narrow, wide int32, data string) error {
	dataPtr, err := windows.BytePtrFromString(data)
	if err != nil {
		return err
	}
	r, _, _ := s.procTSPLBarCode.Call(
		handle,
		uintptr(xPos),
		uintptr(yPos),
		uintptr(codeType),
		uintptr(height),
		uintptr(readable),
		uintptr(rotation),
		uintptr(narrow),
		uintptr(wide),
		uintptr(unsafe.Pointer(dataPtr)),
	)
	ret := int32(r)
	if ret != 0 {
		return fmt.Errorf("TSPL_BarCode 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return nil
}

func (s *tsplSDK) tsplQrCode(handle uintptr, xPos, yPos, eccLevel, width, mode, rotation, model, mask int32, data string) error {
	dataPtr, err := windows.BytePtrFromString(data)
	if err != nil {
		return err
	}
	r, _, _ := s.procTSPLQrCode.Call(
		handle,
		uintptr(xPos),
		uintptr(yPos),
		uintptr(eccLevel),
		uintptr(width),
		uintptr(mode),
		uintptr(rotation),
		uintptr(model),
		uintptr(mask),
		uintptr(unsafe.Pointer(dataPtr)),
	)
	ret := int32(r)
	if ret != 0 {
		return fmt.Errorf("TSPL_QrCode 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return nil
}

func (s *tsplSDK) tsplPrint(handle uintptr, num, copies int32) error {
	r, _, _ := s.procTSPLPrint.Call(handle, uintptr(num), uintptr(copies))
	ret := int32(r)
	if ret != 0 {
		return fmt.Errorf("TSPL_Print 失败: code=%d msg=%s", ret, s.formatError(ret))
	}
	return nil
}

func bytesCStringLen(b []byte) int {
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			return i
		}
	}
	return len(b)
}
