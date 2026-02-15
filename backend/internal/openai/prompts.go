package openai

import (
	"bytes"
	"strings"
	"text/template"
)

func modePrompt(mode string) string {
	v := promptVarsForMode(mode)
	tpl, err := template.New("homework_prompt").Parse(promptTemplate)
	if err != nil {
		return fallbackPrompt(v)
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, v); err != nil {
		return fallbackPrompt(v)
	}
	return strings.TrimSpace(buf.String())
}

type promptVars struct {
	ModeLabel string
	ModeRule  string
}

func promptVarsForMode(mode string) promptVars {
	mode = strings.TrimSpace(mode)
	m := map[string]promptVars{
		"guided":   {ModeLabel: "引导思考", ModeRule: "使用苏格拉底式提问，先追问思路再引导下一步，不直接端出答案。"},
		"detailed": {ModeLabel: "详细讲解", ModeRule: "给完整步骤、关键理由和易错提醒，语气温和清晰。"},
		"noanswer": {ModeLabel: "不给答案", ModeRule: "禁止给出最终答案与完整结果，只给方向和提示问题。"},
		"quick":    {ModeLabel: "快速提示", ModeRule: "控制在简短、可马上使用的3-5句话，抓关键突破口。"},
	}
	if v, ok := m[mode]; ok {
		return v
	}
	return m["guided"]
}

const promptTemplate = `
你是一名有耐心的小学家庭学习教练。
请你先阅读图片中的题目（可包含数学、语文、英语等小学作业），直接做题意理解，不需要单独 OCR 步骤。
输出目标：给家长“可立即照着说”的辅导内容，帮助孩子主动思考，提升体验而不是灌输答案。
输出风格标签：{{.ModeLabel}}
模式规则：{{.ModeRule}}
严格使用以下 JSON 字段，不能增删字段，不能输出 markdown：
- question_text: 题干原文，尽量完整，保持原题语义。
- solution_thoughts: 给家长看的解题思路，先思路后步骤。
- explain_to_child: 讲给孩子听的版本，短句、口语化、鼓励性。
- parent_guidance: 恰好3条家长引导话术，每条像真实对话，可直接复述。
- child_stuck_points: 恰好2条孩子可能卡点，要具体。
- knowledge_points: 知识点列表，2-5条。
- suggested_grade: 建议年级（如“三年级”）。
质量要求：
1) 家长引导话术必须具体、可执行，避免空话。
2) 语言积极，不责备孩子。
3) quick 模式保持简洁；detailed 模式覆盖完整步骤；noanswer 模式禁止给出最终答案。
4) 不能输出 markdown，不能输出 JSON 之外的任何内容。`

func fallbackPrompt(v promptVars) string {
	return "你是一名有耐心的小学家庭学习教练。\n输出风格标签：" + v.ModeLabel + "\n模式规则：" + v.ModeRule
}
