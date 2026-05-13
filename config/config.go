package config

// MaaConfig 存放 MaaFramework 初始化所需的配置项
type MaaConfig struct {
	// LibDir 原生动态库（DLL/SO）所在目录路径
	LibDir string
	// ResourceDir 资源目录路径，包含 pipeline JSON、图片、OCR 模型等
	ResourceDir string
	// AdbPath adb 可执行文件路径
	AdbPath string
	// AdbAddress 设备地址，如 "127.0.0.1:5555"
	AdbAddress string
	// AgentPath MaaAgentBinary 目录路径
	AgentPath string
	// Screencap 截图方式，如 "RawByNc"、"Encode" 等
	Screencap string
	// Input 输入方式，如 "MaaTouch"、"Maatouch"、"AdbShell" 等
	Input string
	// AdbConfig ADB 配置 JSON 字符串
	AdbConfig string
}

// DefaultConfig 返回默认配置
// 请根据实际环境修改 LibDir、AdbPath、AdbAddress、AgentPath 等字段
func DefaultConfig() *MaaConfig {
	return &MaaConfig{
		LibDir:      "./MaaFramework/bin",       // 原生库目录
		ResourceDir: "./resource",               // 资源目录
		AdbPath:     "adb",                       // adb 路径
		AdbAddress:  "127.0.0.1:5555",           // 默认模拟器地址
		AgentPath:   "./MaaAgentBinary",         // Agent 二进制文件路径
		Screencap:   "RawByNc",                  // 截图方式
		Input:       "MaaTouch",                 // 输入方式
		AdbConfig:   "{}",                        // ADB 额外配置
	}
}
