package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config 存储所有配置信息
type Config struct {
	InputFile         string  // 输入文件路径
	OutputFile        string  // 输出文件路径
	SiliconFlowAPIKey string  `json:"silicon_flow_api_key"`
	Model             string  `json:"model"`
	MaxTokens         int     `json:"max_tokens"`
	Temperature       float64 `json:"temperature"`
	TopP              float64 `json:"top_p"`
	TopK              int     `json:"top_k"`
	FrequencyPenalty  float64 `json:"frequency_penalty"`
}

const configFileName = "config.json"
const appConfigDirName = ".ai-js-anti-obfuscation"

// New 创建默认配置
func New() *Config {
	return &Config{
		// 设置默认值
		Model:            "Qwen/Qwen2.5-Coder-32B-Instruct",
		MaxTokens:        4096,
		Temperature:      0.2,
		TopP:             0.9,
		TopK:             50,
		FrequencyPenalty: 0.0,
	}
}

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return configFileName
	}
	return filepath.Join(homeDir, appConfigDirName, configFileName)
}

func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return &Config{}, nil
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func SaveConfig(config *Config) error {
	configPath := GetConfigPath()

	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
