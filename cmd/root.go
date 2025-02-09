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
		// 加载配置文件
		fileConfig, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		// 只有当命令行参数没有设置时，才使用配置文件中的值
		if cmd.Flags().Changed("model") == false && fileConfig.Model != "" {
			cfg.Model = fileConfig.Model
		}
		if cmd.Flags().Changed("max-tokens") == false && fileConfig.MaxTokens != 0 {
			cfg.MaxTokens = fileConfig.MaxTokens
		}
		if cmd.Flags().Changed("temperature") == false && fileConfig.Temperature != 0 {
			cfg.Temperature = fileConfig.Temperature
		}
		if cmd.Flags().Changed("top-p") == false && fileConfig.TopP != 0 {
			cfg.TopP = fileConfig.TopP
		}
		if cmd.Flags().Changed("top-k") == false && fileConfig.TopK != 0 {
			cfg.TopK = fileConfig.TopK
		}
		if cmd.Flags().Changed("frequency-penalty") == false && fileConfig.FrequencyPenalty != 0 {
			cfg.FrequencyPenalty = fileConfig.FrequencyPenalty
		}

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
		return siliconflow.RunDeobfuscate(cfg.InputFile, cfg.OutputFile, cfg)
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

	// 模型参数
	rootCmd.Flags().StringVar(&cfg.Model, "model", cfg.Model, "Model name")
	rootCmd.Flags().IntVar(&cfg.MaxTokens, "max-tokens", cfg.MaxTokens, "Maximum number of tokens to generate")
	rootCmd.Flags().Float64Var(&cfg.Temperature, "temperature", cfg.Temperature, "Sampling temperature (0.0-1.0)")
	rootCmd.Flags().Float64Var(&cfg.TopP, "top-p", cfg.TopP, "Top-p sampling parameter")
	rootCmd.Flags().IntVar(&cfg.TopK, "top-k", cfg.TopK, "Top-k sampling parameter")
	rootCmd.Flags().Float64Var(&cfg.FrequencyPenalty, "frequency-penalty", cfg.FrequencyPenalty, "Frequency penalty parameter")

	// 只标记输入文件为必需参数
	rootCmd.MarkFlagRequired("input")
}
