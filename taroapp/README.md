# Taro 小程序端

目录：`/Users/dinggh/projects/whatsdot-aibuddy/taroapp`

## 页面
- `Home`：拍作业/从相册上传
- `Loading`：上传与分析中
- `Result`：题目+三段结构+模式切换+再生成
- `History`：历史列表，点击查看详情
- `Profile`：占位页

## 安装依赖
```bash
cd /Users/dinggh/projects/whatsdot-aibuddy/taroapp
npm install
```

## 启动微信小程序构建
```bash
npm run dev:weapp
```
构建输出：`/Users/dinggh/projects/whatsdot-aibuddy/taroapp/dist`

## 构建
```bash
npm run build:weapp
```

## 微信开发者工具联调
1. 导入目录：`/Users/dinggh/projects/whatsdot-aibuddy/taroapp/dist`
2. 打开“不校验合法域名/HTTPS”用于本地调试
3. 确保后端运行在 `http://127.0.0.1:8080`

## API 基地址
- 通过 `process.env.TARO_APP_API_BASE` 注入
- 默认值：`http://127.0.0.1:8080`
