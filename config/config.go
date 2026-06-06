package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const defaultConfigPath = "./config/app.json"

// MaaConfig 存放 MaaFramework 初始化所需的配置项。
type MaaConfig struct {
	// LibDir 原生动态库（DLL/SO）所在目录路径。
	LibDir string
	// ResourceDir 资源目录路径，包含 pipeline JSON、图片、OCR 模型等。
	ResourceDir string
	// LogDir MaaFramework 日志和标记图保存目录。
	LogDir string
	// ConfigDir MaaFramework 工具包配置目录。
	ConfigDir string
	// TargetWindowTitle 可选的自动连接窗口标题。
	TargetWindowTitle string
	// AutoConnect 启动后是否自动连接 TargetWindowTitle。
	AutoConnect bool
	// SaveDraw 是否保存识别标记图。
	SaveDraw bool
	// DebugMode 是否启用 MaaFramework debug 模式。
	DebugMode bool
	// StdoutLevel MaaFramework stdout 日志级别。
	StdoutLevel int
	// ServerAddr Web 服务监听地址。
	ServerAddr string
	// AdbPath adb 可执行文件路径。
	AdbPath string
	// AdbAddress 设备地址，如 "127.0.0.1:5555"。
	AdbAddress string
	// AgentPath MaaAgentBinary 目录路径。
	AgentPath string
	// Screencap 截图方式，如 "RawByNc"、"Encode" 等。
	Screencap string
	// Input 输入方式，如 "MaaTouch"、"Maatouch"、"AdbShell" 等。
	Input string
	// AdbConfig ADB 配置 JSON 字符串。
	AdbConfig string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *MaaConfig {
	return &MaaConfig{
		LibDir:      "./MaaFramework/bin",
		ResourceDir: "./resource",
		LogDir:      "./log",
		ConfigDir:   ".",
		SaveDraw:    true,
		DebugMode:   true,
		StdoutLevel: 2,
		ServerAddr:  ":8080",
		AdbPath:     "adb",
		AdbAddress:  "127.0.0.1:5555",
		AgentPath:   "./MaaAgentBinary",
		Screencap:   "RawByNc",
		Input:       "MaaTouch",
		AdbConfig:   "{}",
	}
}

// LoadConfig loads defaults, then optional JSON, then environment overrides.
func LoadConfig() (*MaaConfig, error) {
	cfg := DefaultConfig()
	configPath := envString("MAA_CONFIG", defaultConfigPath)
	if configPath != "" {
		if data, err := os.ReadFile(configPath); err == nil {
			if err := json.Unmarshal(data, cfg); err != nil {
				return nil, err
			}
		} else if !os.IsNotExist(err) {
			return nil, err
		}
	}

	cfg.LibDir = envString("MAA_LIB_DIR", cfg.LibDir)
	cfg.ResourceDir = envString("MAA_RESOURCE_DIR", cfg.ResourceDir)
	cfg.LogDir = envString("MAA_LOG_DIR", cfg.LogDir)
	cfg.ConfigDir = envString("MAA_CONFIG_DIR", cfg.ConfigDir)
	cfg.TargetWindowTitle = envString("MAA_TARGET_WINDOW_TITLE", cfg.TargetWindowTitle)
	cfg.ServerAddr = envString("MAA_SERVER_ADDR", cfg.ServerAddr)
	cfg.AdbPath = envString("MAA_ADB_PATH", cfg.AdbPath)
	cfg.AdbAddress = envString("MAA_ADB_ADDRESS", cfg.AdbAddress)
	cfg.AgentPath = envString("MAA_AGENT_PATH", cfg.AgentPath)
	cfg.Screencap = envString("MAA_SCREENCAP", cfg.Screencap)
	cfg.Input = envString("MAA_INPUT", cfg.Input)
	cfg.AdbConfig = envString("MAA_ADB_CONFIG", cfg.AdbConfig)
	cfg.AutoConnect = envBool("MAA_AUTO_CONNECT", cfg.AutoConnect)
	cfg.SaveDraw = envBool("MAA_SAVE_DRAW", cfg.SaveDraw)
	cfg.DebugMode = envBool("MAA_DEBUG", cfg.DebugMode)
	cfg.StdoutLevel = envInt("MAA_STDOUT_LEVEL", cfg.StdoutLevel)

	return cfg, nil
}

func (c *MaaConfig) ResourcePath(parts ...string) string {
	items := append([]string{c.ResourceDir}, parts...)
	return filepath.Join(items...)
}

func envString(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}

func envBool(name string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envInt(name string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
