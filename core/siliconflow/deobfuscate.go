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
func RunDeobfuscate(inputFile, outputFile string) error {
	// 加载配置
	c, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// 检查 API key
	if c.SiliconFlowAPIKey == "" {
		fmt.Print("Please enter your SiliconFlow API key: ")
		var apiKey string
		fmt.Scanln(&apiKey)
		c.SiliconFlowAPIKey = apiKey

		if err := config.SaveConfig(c); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}
	}

	// 读取输入文件
	inputCode, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	// 创建客户端并执行反混淆
	client := NewClient(c.SiliconFlowAPIKey)
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
