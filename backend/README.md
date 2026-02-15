# aibuddy-backend (Go + Gin + PostgreSQL)

## 功能
- `POST /api/v1/homework/analyze`：接收图片上传（multipart）并调用 OpenAI 兼容 Chat Completions（vision）
- 输出结构化 JSON：
  - 解题思路
  - 讲给孩子听
  - 家长引导话术（3条）
  - 孩子可能卡点（2条）
  - 知识点
  - 建议年级
- 保存历史记录，支持列表和详情
- 统一响应格式：`{ code, message, data }`
- 基础限流：按 `X-Device-Id`（或 `device_id` query）令牌桶
- CORS 允许本地联调

## 目录
- `cmd/server/main.go`: 服务入口
- `cmd/migrate/main.go`: 迁移入口（默认顺序执行 `migrations/*.sql`）
- `internal/httpapi`: Gin 路由与处理器
- `internal/openai`: OpenAI 兼容 Chat Completions 客户端
- `internal/store`: 数据访问

## 启动前准备
1. 创建 PostgreSQL 数据库，例如 `aibuddy`
2. 配置环境变量：复制 `.env.example` 为 `.env`
3. 执行迁移：
   ```bash
   cd backend
   go run ./cmd/migrate
   ```

## 运行
```bash
cd backend
go run ./cmd/server
```
默认监听 `:8080`

## OpenAI 调用说明
- 默认使用 `OPENAI_BASE_URL/chat/completions`
- 请求包含图片 `data URL`，无需单独 OCR
- 推荐生产配置：
  - `ANALYZE_MOCK=false`
  - `OPENAI_API_KEY` 填真实 key

## API
- `GET /health`
- `POST /api/v1/homework/analyze`
  - Header: `X-Device-Id: xxx`
  - Form: `image=<file>`, `mode=guided|detailed|noanswer|quick`
- `POST /api/v1/homework/:id/regenerate`
  - Header: `X-Device-Id: xxx`
  - Query/Form: `mode=...`
- `GET /api/v1/history`
  - Header: `X-Device-Id: xxx`
- `GET /api/v1/history/:id`
  - Header: `X-Device-Id: xxx`
