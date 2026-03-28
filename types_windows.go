//go:build windows

package main

// ZhyhField 是前端传入的打印字段。
//
// - zhyh_type: 字段类型，例如 text / barcode / qrcode
// - zhyh_key:  字段名称（用于 text 的 key:value）
// - zhyh_value:字段内容
type ZhyhField struct {
	ZhyhType  string `json:"zhyh_type"`
	ZhyhKey   string `json:"zhyh_key"`
	ZhyhValue string `json:"zhyh_value"`
}

// ProductLabel 是常见的“商品标签”结构（如果前端固定字段，可用这个结构）。
type ProductLabel struct {
	Name    string
	SKU     string
	Spec    string
	Price   string
	Barcode string
	QR      string
}

// LabelOptions 是 TSPL_Setup 相关配置。
type LabelOptions struct {
	LabelWidthMM  int32 `json:"label_width_mm"`
	LabelHeightMM int32 `json:"label_height_mm"`
	Speed         int32 `json:"speed"`
	Density       int32 `json:"density"`
	Type          int32 `json:"type"`
	GapMM         int32 `json:"gap_mm"`
	OffsetMM      int32 `json:"offset_mm"`
	DPI           int32 `json:"dpi"`
	MarginLeftMM  int32 `json:"margin_left_mm"`
	MarginTopMM   int32 `json:"margin_top_mm"`
}

// PrintConfig 描述一次打印需要的基础信息。
type PrintConfig struct {
	SDKPath string       `json:"sdk_path"`
	Model   string       `json:"model"`
	Port    string       `json:"port"`
	Options LabelOptions `json:"options"`
}

// AppConfig 是服务端配置文件的结构。
type AppConfig struct {
	HTTPAddr string      `json:"http_addr"`
	Printer  PrintConfig `json:"printer"`
}

// PrinterInfo 返回打印机基础信息，便于前端展示/排障。
type PrinterInfo struct {
	UsedModel string
	Status    int32
	Firmware  string
	SN        string
}
