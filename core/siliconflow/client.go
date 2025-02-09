package siliconflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BaseURL = "https://api.siliconflow.cn/v1"

const deobfuscatePrompt = `你是一个 JavaScript 代码专家。请分析下面的混淆代码，在保持功能完全相同的情况下，将混淆的变量名改为符合其功能的、易读的变量名。
要求：
1. 只修改变量名，不要改变代码逻辑和结构
2. 新的变量名要符合代码中的实际用途
3. 保持代码格式不变
4. 直接返回完整的修改后代码，不要包含任何解释、代码块标记或其他额外内容

代码：
%s`

// APIError represents the error response from SiliconFlow API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("SiliconFlow API error (code=%d): %s", e.Code, e.Message)
}

// APIResponse represents the success response from SiliconFlow API
type APIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Client struct {
	apiKey string
	config *DeobfuscateConfig
}

type DeobfuscateConfig struct {
	Model            string
	MaxTokens        int
	Temperature      float64
	TopP             float64
	TopK             int
	FrequencyPenalty float64
}

func NewClient(apiKey string, config *DeobfuscateConfig) *Client {
	return &Client{
		apiKey: apiKey,
		config: config,
	}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model            string        `json:"model"`
	Messages         []ChatMessage `json:"messages"`
	Stream           bool          `json:"stream"`
	MaxTokens        int           `json:"max_tokens"`
	Stop             []string      `json:"stop,omitempty"`
	Temperature      float64       `json:"temperature"`
	TopP             float64       `json:"top_p"`
	TopK             int           `json:"top_k"`
	FrequencyPenalty float64       `json:"frequency_penalty"`
	N                int           `json:"n"`
}

func (c *Client) Chat(req *ChatRequest) (string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check if response is an error
	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return "", fmt.Errorf("API error (status=%d): %s", resp.StatusCode, string(body))
		}
		return "", &apiErr
	}

	// Parse success response
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	return apiResp.Choices[0].Message.Content, nil
}

func (c *Client) Deobfuscate(code string) (string, error) {
	req := &ChatRequest{
		Model: c.config.Model,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: fmt.Sprintf(deobfuscatePrompt, code),
			},
		},
		Stream:           false,
		MaxTokens:        c.config.MaxTokens,
		Temperature:      c.config.Temperature,
		TopP:             c.config.TopP,
		TopK:             c.config.TopK,
		FrequencyPenalty: c.config.FrequencyPenalty,
		N:                1,
	}

	return c.Chat(req)
}
