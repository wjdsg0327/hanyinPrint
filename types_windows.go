//go:build windows

package main

// ZhyhField 是前端传入的打印字段。
//
// - zhyh_type: 字段类型，例如 text / barcode / qrcode
// - zhyh_key:  字段名称（用于 text 的 key:value）
// - zhyh_value:字段内容
type ZhyhField struct {
	ZhyhType  string      `json:"zhyh_type"`
	ZhyhKey   string      `json:"zhyh_key"`
	ZhyhValue interface{} `json:"zhyh_value"`
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

//- label_width_mm ：标签宽度（mm）
//- label_height_mm ：标签高度（mm）
//- dpi ：打印机分辨率（常见 203 或 300），用于把 mm 转成坐标点（dots）
//- speed ：打印速度（SDK 参数，整数）
//- density ：打印浓度（SDK 参数，整数）
//- type ：纸张类型（SDK 参数，整数，示例里用 1）
//- gap_mm ：标签间隙（mm）
//- offset_mm ：偏移（mm）
//和“内容位置”有关：
//
//- margin_left_mm ：左边距（mm），越大内容越往右
//- margin_top_mm ：上边距（mm），越大内容越往下

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

// LogConfig 描述日志输出配置。
type LogConfig struct {
	FilePath string `json:"file_path"`
	Level    string `json:"level"`
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
	HTTPAddr   string      `json:"http_addr"`
	TenantCode string      `json:"tenant_code"`
	Key        string      `json:"key"`
	Printer    PrintConfig `json:"printer"`
	Log        LogConfig   `json:"log"`
}

// PrinterInfo 返回打印机基础信息，便于前端展示/排障。
type PrinterInfo struct {
	UsedModel string
	Status    int32
	Firmware  string
	SN        string
}
