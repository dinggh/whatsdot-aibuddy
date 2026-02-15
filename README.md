# whatsdot-aibuddy

微信小程序「拍作业」端到端 Demo：
- 小程序选图/拍照上传
- Go(Gin) 服务端调用 OpenAI 兼容 Chat Completions（vision）
- 返回结构化 JSON 并展示重点引导话术
- 历史记录列表 + 详情复用结果页

## 1. 启动后端
```bash
cd /Users/dinggh/projects/whatsdot-aibuddy/backend
cp .env.example .env
# 编辑 .env：DATABASE_URL、OPENAI_API_KEY（或 ANALYZE_MOCK=true）
go run ./cmd/migrate
go run ./cmd/server
```

## 2. 启动小程序构建
```bash
cd /Users/dinggh/projects/whatsdot-aibuddy/taroapp
npm install
npm run dev:weapp
```
微信开发者工具导入：`/Users/dinggh/projects/whatsdot-aibuddy/taroapp/dist`

## 3. 自测流程（端到端）
1. 打开首页，点击「拍作业」或「从相册上传」
2. 选择图片后进入 Loading，自动上传到后端
3. 自动跳转 Result，检查以下结构：
   - 题目预览图 + 题干文本
   - 解题思路
   - 讲给孩子听
   - 家长可以这样引导（3条，重点卡片）
   - 孩子可能卡点（2条）/知识点/建议年级
4. 点击「再生成一次」验证 `mode` 可再次生成
5. 点击右上角模式切换（引导思考/详细讲解/不给答案/快速提示）并验证结果更新
6. 进入「历史记录」页，查看列表并点开任一条，确认详情可复用 Result 渲染

## 4. 截图占位
把你的联调截图放到：
- `/Users/dinggh/projects/whatsdot-aibuddy/docs/screenshots/home.png`
- `/Users/dinggh/projects/whatsdot-aibuddy/docs/screenshots/result.png`
- `/Users/dinggh/projects/whatsdot-aibuddy/docs/screenshots/history.png`
- `/Users/dinggh/projects/whatsdot-aibuddy/docs/screenshots/mode.png`

## 5. 接口摘要
- `POST /api/v1/homework/analyze`
- `POST /api/v1/homework/:id/regenerate`
- `GET /api/v1/history`
- `GET /api/v1/history/:id`

请求需带 `X-Device-Id`，返回统一格式：
```json
{ "code": 0, "message": "ok", "data": {} }
```
