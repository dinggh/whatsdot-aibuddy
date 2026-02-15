package openai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL string
	APIKey  string
	Model   string
	HTTP    *http.Client
}

type AnalyzeResult struct {
	QuestionText     string   `json:"question_text"`
	SolutionThoughts string   `json:"solution_thoughts"`
	ExplainToChild   string   `json:"explain_to_child"`
	ParentGuidance   []string `json:"parent_guidance"`
	ChildStuckPoints []string `json:"child_stuck_points"`
	KnowledgePoints  []string `json:"knowledge_points"`
	SuggestedGrade   string   `json:"suggested_grade"`
}

func New(baseURL, apiKey, model string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  apiKey,
		Model:   model,
		HTTP:    &http.Client{Timeout: 45 * time.Second},
	}
}

func (c *Client) AnalyzeHomework(ctx context.Context, imageBytes []byte, contentType string, mode string) (AnalyzeResult, error) {
	if c.APIKey == "" {
		return AnalyzeResult{}, errors.New("OPENAI_API_KEY is empty")
	}
	mediaType := normalizeContentType(contentType)
	if mediaType == "" {
		mediaType = "image/jpeg"
	}
	imageDataURL := "data:" + mediaType + ";base64," + base64.StdEncoding.EncodeToString(imageBytes)

	instructions := modePrompt(mode)
	payload := map[string]any{
		"model": c.Model,
		"messages": []map[string]any{
			{
				"role":    "system",
				"content": "你是一名小学家庭教育辅导助手。严格输出 JSON，不要输出 markdown。",
			},
			{
				"role": "user",
				"content": []map[string]any{
					{"type": "text", "text": instructions},
					{"type": "image_url", "image_url": map[string]any{"url": imageDataURL}},
				},
			},
		},
		"response_format": map[string]any{
			"type": "json_schema",
			"json_schema": map[string]any{
				"name":   "homework_analysis",
				"strict": true,
				"schema": analysisSchema(),
			},
		},
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return AnalyzeResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/chat/completions", bytes.NewReader(buf))
	if err != nil {
		return AnalyzeResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return AnalyzeResult{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return AnalyzeResult{}, fmt.Errorf("openai status %d: %s", resp.StatusCode, string(body))
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return AnalyzeResult{}, err
	}
	if len(parsed.Choices) == 0 {
		return AnalyzeResult{}, errors.New("empty choices")
	}

	content := strings.TrimSpace(parsed.Choices[0].Message.Content)
	if content == "" {
		return AnalyzeResult{}, errors.New("empty completion content")
	}

	var out AnalyzeResult
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		return AnalyzeResult{}, fmt.Errorf("invalid completion json: %w", err)
	}
	return normalize(out), nil
}

func normalize(input AnalyzeResult) AnalyzeResult {
	out := input
	out.QuestionText = strings.TrimSpace(out.QuestionText)
	out.SolutionThoughts = strings.TrimSpace(out.SolutionThoughts)
	out.ExplainToChild = strings.TrimSpace(out.ExplainToChild)
	out.SuggestedGrade = strings.TrimSpace(out.SuggestedGrade)
	if len(out.ParentGuidance) > 3 {
		out.ParentGuidance = out.ParentGuidance[:3]
	}
	if len(out.ChildStuckPoints) > 2 {
		out.ChildStuckPoints = out.ChildStuckPoints[:2]
	}
	return out
}

func normalizeContentType(contentType string) string {
	contentType = strings.TrimSpace(contentType)
	if contentType == "" {
		return ""
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return contentType
	}
	return mediaType
}

func modePrompt(mode string) string {
	mode = strings.TrimSpace(mode)
	if mode == "" {
		mode = "guided"
	}
	guide := map[string]string{
		"guided":   "强调提问式引导，不直接给最终答案。",
		"detailed": "给出详细步骤和原理，语言清晰。",
		"noanswer": "不给最终答案，只给思路和启发问题。",
		"quick":    "给简短高效提示，适合快速辅导。",
	}[mode]
	if guide == "" {
		guide = "强调提问式引导，不直接给最终答案。"
	}
	return "请识别图片中的题目并输出结构化 JSON。要求：\n" +
		"1) question_text: 题干原文，尽量完整。\n" +
		"2) solution_thoughts: 给家长看的解题思路。\n" +
		"3) explain_to_child: 讲给孩子听的版本。\n" +
		"4) parent_guidance: 恰好 3 条家长引导话术。\n" +
		"5) child_stuck_points: 恰好 2 条孩子可能卡点。\n" +
		"6) knowledge_points: 知识点列表。\n" +
		"7) suggested_grade: 建议年级。\n" +
		"模式要求：" + guide
}

func analysisSchema() map[string]any {
	return map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required": []string{
			"question_text",
			"solution_thoughts",
			"explain_to_child",
			"parent_guidance",
			"child_stuck_points",
			"knowledge_points",
			"suggested_grade",
		},
		"properties": map[string]any{
			"question_text":     map[string]any{"type": "string"},
			"solution_thoughts": map[string]any{"type": "string"},
			"explain_to_child":  map[string]any{"type": "string"},
			"parent_guidance": map[string]any{
				"type": "array", "minItems": 3, "maxItems": 3,
				"items": map[string]any{"type": "string"},
			},
			"child_stuck_points": map[string]any{
				"type": "array", "minItems": 2, "maxItems": 2,
				"items": map[string]any{"type": "string"},
			},
			"knowledge_points": map[string]any{
				"type": "array", "items": map[string]any{"type": "string"},
			},
			"suggested_grade": map[string]any{"type": "string"},
		},
	}
}
