package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ai-js-anti-obfuscation/core/config"
	"ai-js-anti-obfuscation/core/siliconflow"

	"github.com/spf13/cobra"
)

var (
	cfg = config.New()
)

var rootCmd = &cobra.Command{
	Use:   "ai-js-anti-obfuscation",
	Short: "AI JavaScript Anti-obfuscation Tool",
	Long:  `A tool to deobfuscate JavaScript code using AI`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// 验证输入文件存在
		if _, err := os.Stat(cfg.InputFile); os.IsNotExist(err) {
			// 文件不存在属于用法错误，应该显示使用说明
			cmd.SilenceUsage = false
			return fmt.Errorf("input file does not exist: %s", cfg.InputFile)
		}

		// 如果没有指定输出文件，生成默认输出文件路径
		if cfg.OutputFile == "" {
			ext := filepath.Ext(cfg.InputFile)
			baseName := strings.TrimSuffix(cfg.InputFile, ext)
			cfg.OutputFile = baseName + "_output" + ext
		}

		// 确保输出文件的目录存在
		outputDir := filepath.Dir(cfg.OutputFile)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// 运行时错误不显示使用说明
		cmd.SilenceUsage = true
		return siliconflow.RunDeobfuscate(cfg.InputFile, cfg.OutputFile)
	},
	// 禁用错误追踪
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// 必需参数
	rootCmd.Flags().StringVarP(&cfg.InputFile, "input", "i", "", "Input file path")
	rootCmd.Flags().StringVarP(&cfg.OutputFile, "output", "o", "", "Output file path (optional, defaults to input_output.js)")

	// 只标记输入文件为必需参数
	rootCmd.MarkFlagRequired("input")
}
