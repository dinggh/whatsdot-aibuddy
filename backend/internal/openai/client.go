package openai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"strings"
	"time"

	oosdk "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

type Client struct {
	BaseURL string
	APIKey  string
	Model   string
	SDK     oosdk.Client
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
	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
		option.WithBaseURL(strings.TrimRight(baseURL, "/")),
		option.WithRequestTimeout(45 * time.Second),
	}
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  apiKey,
		Model:   model,
		SDK:     oosdk.NewClient(opts...),
	}
}

func (c *Client) AnalyzeHomework(ctx context.Context, imageBytes []byte, contentType string, mode string) (AnalyzeResult, error) {
	if strings.TrimSpace(c.APIKey) == "" {
		return AnalyzeResult{}, errors.New("OPENAI_API_KEY is empty")
	}

	mediaType := normalizeContentType(contentType)
	if mediaType == "" {
		mediaType = "image/jpeg"
	}
	imageDataURL := "data:" + mediaType + ";base64," + base64.StdEncoding.EncodeToString(imageBytes)

	messages := []oosdk.ChatCompletionMessageParamUnion{
		oosdk.SystemMessage(systemPrompt()),
		oosdk.UserMessage([]oosdk.ChatCompletionContentPartUnionParam{
			oosdk.TextContentPart(modePrompt(mode)),
			oosdk.ImageContentPart(oosdk.ChatCompletionContentPartImageImageURLParam{URL: imageDataURL, Detail: "high"}),
		}),
	}

	resp, err := c.SDK.Chat.Completions.New(ctx, oosdk.ChatCompletionNewParams{
		Model:    shared.ChatModel(c.Model),
		Messages: messages,
		ResponseFormat: oosdk.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:        "homework_analysis",
					Description: oosdk.String("Homework analysis JSON for parent guidance in Chinese"),
					Strict:      oosdk.Bool(true),
					Schema:      analysisSchema(),
				},
			},
		},
		Temperature: oosdk.Float(0.2),
	})
	if err != nil {
		return AnalyzeResult{}, fmt.Errorf("chat completion failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return AnalyzeResult{}, errors.New("empty choices")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
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

func systemPrompt() string {
	return "你是一名有耐心的小学家庭学习教练和家长沟通顾问。" +
		"你的目标不是替孩子做题，而是帮助家长通过提问让孩子自己思考。" +
		"输出必须是严格 JSON，不能输出 markdown、不能输出解释文字。"
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
				"type": "array", "minItems": 2, "maxItems": 5,
				"items": map[string]any{"type": "string"},
			},
			"suggested_grade": map[string]any{"type": "string"},
		},
	}
}
