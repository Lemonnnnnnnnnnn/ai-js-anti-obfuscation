package siliconflow

import (
	"fmt"
	"os"
	"strings"

	"ai-js-anti-obfuscation/core/config"
)

// cleanOutput 清理输出内容，移除可能的代码块标记
func cleanOutput(output string) string {
	// 移除开头的 ```javascript 或 ```js
	output = strings.TrimPrefix(output, "```javascript")
	output = strings.TrimPrefix(output, "```js")
	output = strings.TrimPrefix(output, "```")

	// 移除结尾的 ```
	output = strings.TrimSuffix(output, "```")

	// 清理首尾的空白字符
	return strings.TrimSpace(output)
}

// RunDeobfuscate 执行反混淆的主要逻辑
func RunDeobfuscate(inputFile, outputFile string, cliConfig *config.Config) error {
	// 加载配置文件
	fileConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// 创建默认配置
	defaultConfig := config.New()

	// 按优先级合并配置：命令行 > 配置文件 > 默认值
	finalConfig := &DeobfuscateConfig{
		// 如果命令行参数非空则使用命令行参数，否则使用配置文件参数，如果配置文件也为空则使用默认值
		Model:            getConfigValue(cliConfig.Model, fileConfig.Model, defaultConfig.Model),
		MaxTokens:        getConfigIntValue(cliConfig.MaxTokens, fileConfig.MaxTokens, defaultConfig.MaxTokens),
		Temperature:      getConfigFloatValue(cliConfig.Temperature, fileConfig.Temperature, defaultConfig.Temperature),
		TopP:             getConfigFloatValue(cliConfig.TopP, fileConfig.TopP, defaultConfig.TopP),
		TopK:             getConfigIntValue(cliConfig.TopK, fileConfig.TopK, defaultConfig.TopK),
		FrequencyPenalty: getConfigFloatValue(cliConfig.FrequencyPenalty, fileConfig.FrequencyPenalty, defaultConfig.FrequencyPenalty),
	}

	// 检查 API key
	apiKey := cliConfig.SiliconFlowAPIKey
	if apiKey == "" {
		apiKey = fileConfig.SiliconFlowAPIKey
	}
	if apiKey == "" {
		fmt.Print("Please enter your SiliconFlow API key: ")
		fmt.Scanln(&apiKey)
		fileConfig.SiliconFlowAPIKey = apiKey
		if err := config.SaveConfig(fileConfig); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}
	}

	// 读取输入文件
	inputCode, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	// 创建客户端并执行反混淆
	client := NewClient(apiKey, finalConfig)
	fmt.Println("Deobfuscating code...")

	deobfuscated, err := client.Deobfuscate(string(inputCode))
	if err != nil {
		return err
	}

	// 清理输出内容
	deobfuscated = cleanOutput(deobfuscated)

	// 写入输出文件
	if err := os.WriteFile(outputFile, []byte(deobfuscated), 0644); err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	fmt.Printf("Successfully deobfuscated code and saved to: %s\n", outputFile)
	return nil
}

// 辅助函数：获取字符串配置值
func getConfigValue(cli, file, defaultVal string) string {
	if cli != "" {
		return cli
	}
	if file != "" {
		return file
	}
	return defaultVal
}

// 辅助函数：获取整数配置值
func getConfigIntValue(cli, file, defaultVal int) int {
	if cli != 0 {
		return cli
	}
	if file != 0 {
		return file
	}
	return defaultVal
}

// 辅助函数：获取浮点数配置值
func getConfigFloatValue(cli, file, defaultVal float64) float64 {
	if cli != 0 {
		return cli
	}
	if file != 0 {
		return file
	}
	return defaultVal
}
