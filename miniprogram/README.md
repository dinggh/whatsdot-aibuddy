# 微信小程序（原生）

目录：`/Users/dinggh/projects/whatsdot-aibuddy/miniprogram`

## 当前状态
- 页面已接入后端 API，不再是纯静态假数据。
- 首次进入会调用 `wx.login` 完成服务端登录。
- 「我的」页面支持：
  - 同步昵称/头像（`wx.getUserProfile`）
  - 绑定手机号（`getPhoneNumber` -> 服务端换取）

## 依赖后端
后端目录：`/Users/dinggh/projects/whatsdot-aibuddy/backend`

请先启动后端（默认 `http://127.0.0.1:8080`），再打开小程序。

## 本地测试（微信开发者工具）
1. 打开微信开发者工具。
2. 选择「导入项目」。
3. 项目目录：`/Users/dinggh/projects/whatsdot-aibuddy/miniprogram`。
4. 开发调试阶段可关闭「校验合法域名」。
5. 打开首页后会自动登录并拉取用户/历史数据。

## 重要说明
- `touristappid` 无法用于真实微信登录、手机号能力。
- 要测试真实昵称/手机号，请使用你自己的小程序 `AppID`，并配置对应 `WECHAT_APP_ID / WECHAT_APP_SECRET` 到后端环境变量。
- 真机与线上必须 HTTPS 且在小程序后台配置合法请求域名。
