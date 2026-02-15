package openai

import (
	"strings"
	"testing"
)

func TestModePromptContainsStructuredExperienceGuidance(t *testing.T) {
	p := modePrompt("guided")

	checks := []string{
		"你是一名有耐心的小学家庭学习教练",
		"家长引导话术",
		"孩子可能卡点",
		"不能输出 markdown",
		"question_text",
		"suggested_grade",
	}
	for _, c := range checks {
		if !strings.Contains(p, c) {
			t.Fatalf("expected prompt to contain %q, got: %s", c, p)
		}
	}
}

func TestModePromptDifferByMode(t *testing.T) {
	guided := modePrompt("guided")
	detailed := modePrompt("detailed")
	noanswer := modePrompt("noanswer")
	quick := modePrompt("quick")

	if !strings.Contains(guided, "苏格拉底式提问") {
		t.Fatalf("guided mode prompt should require questioning style")
	}
	if !strings.Contains(detailed, "完整步骤") {
		t.Fatalf("detailed mode prompt should require complete steps")
	}
	if !strings.Contains(noanswer, "禁止给出最终答案") {
		t.Fatalf("noanswer mode prompt should forbid final answer")
	}
	if !strings.Contains(quick, "控制在简短") {
		t.Fatalf("quick mode prompt should require concise style")
	}
}

func TestModePromptIsRenderedFromTemplateVariables(t *testing.T) {
	p := modePrompt("guided")
	if strings.Contains(p, "{{") || strings.Contains(p, "}}") {
		t.Fatalf("expected rendered prompt without template tokens, got: %s", p)
	}
	if !strings.Contains(p, "输出风格标签：引导思考") {
		t.Fatalf("expected mode label to be injected into prompt, got: %s", p)
	}
}
