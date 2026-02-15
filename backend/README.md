# aibuddy-backend (Go + PostgreSQL)

## 功能
- 微信小程序登录：`wx.login` code -> `code2session`
- 用户建档与登录态：JWT（无 Redis）
- 获取并绑定手机号：`getPhoneNumber` code -> 微信服务端接口
- 更新昵称/头像：客户端 `wx.getUserProfile` 后上报
- 查询用户信息与历史记录

## 目录
- `cmd/server/main.go`: 服务入口
- `internal/...`: 业务代码
- `migrations/001_init.sql`: PostgreSQL 初始化脚本

## 启动前准备
1. PostgreSQL 创建数据库，例如 `aibuddy`
2. 执行初始化 SQL:
   ```bash
   psql "$DATABASE_URL" -f migrations/001_init.sql
   ```
3. 配置环境变量（参考 `.env.example`）

## 运行
```bash
cd backend
go mod tidy
go run ./cmd/server
```

默认监听 `:8080`。

## API
- `GET /health`
- `POST /api/v1/auth/wechat/login`
  - body: `{ "code": "wx.login返回code" }`
  - resp: `{ "token": "...", "user": {...} }`
- `POST /api/v1/auth/wechat/profile` (Bearer token)
  - body: `{ "nickName": "张妈妈", "avatarUrl": "https://..." }`
- `POST /api/v1/auth/wechat/phone` (Bearer token)
  - body: `{ "code": "getPhoneNumber回调里的code" }`
- `GET /api/v1/me` (Bearer token)
- `GET /api/v1/history` (Bearer token)

## 小程序联调
- 开发阶段可以在微信开发者工具中关闭“校验合法域名”。
- 真机与线上环境必须使用 HTTPS 且在小程序后台配置合法请求域名。

## 注意
- 昵称/头像不是通过 `wx.login` 获取，而是客户端授权后上报到服务端。
- 手机号必须通过 `getPhoneNumber` 的 `code` 在服务端换取，不能前端直接解密。
