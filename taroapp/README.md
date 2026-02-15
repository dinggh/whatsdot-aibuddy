# Taro 统一工程（React 一套代码：H5 + 微信小程序）

目录：`/Users/dinggh/projects/whatsdot-aibuddy/taroapp`

## 已实现页面
- `src/pages/home`
- `src/pages/loading`
- `src/pages/result`
- `src/pages/mode`
- `src/pages/history`
- `src/pages/profile`

## 安装依赖
```bash
cd /Users/dinggh/projects/whatsdot-aibuddy/taroapp
npm install
```

## 启动 H5（本地网页）
```bash
npm run dev:h5
```
默认会启动/监听 H5 构建（端口参数已固定为 `10087`），按终端提示访问本地地址。

## 启动微信小程序构建
```bash
npm run dev:weapp
```
构建输出目录：`/Users/dinggh/projects/whatsdot-aibuddy/taroapp/dist`

## 微信开发者工具测试
1. 打开微信开发者工具
2. 导入项目目录：`/Users/dinggh/projects/whatsdot-aibuddy/taroapp/dist`
3. AppID 选择测试号（或你自己的小程序 AppID）
4. 编译后即可调试 6 个页面

## 常用构建
```bash
npm run build:h5
npm run build:weapp
```

## 注意
- 脚本已默认加 `--no-check`，用于绕过 Taro `doctor` 在部分 Node/macOS 组合下的启动崩溃。
