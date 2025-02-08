package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config 存储所有配置信息
type Config struct {
	InputFile         string // 输入文件路径
	OutputFile        string // 输出文件路径
	SiliconFlowAPIKey string `json:"silicon_flow_api_key"`
}

const configFileName = "config.json"
const appConfigDirName = ".ai-js-anti-obfuscation"

// New 创建默认配置
func New() *Config {
	return &Config{}
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
