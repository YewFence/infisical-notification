# Notification - Infisical 通知转发器

一个轻量的 Appwrite Function：接收 Infisical 的 Webhook，校验签名后将消息转发到 Apprise。

## 主要流程

1. 接收 Infisical Webhook 请求
2. 校验 `x-infisical-signature` 签名与时间窗
3. 解析事件与路径，渲染消息模板
4. 调用 Apprise API 推送通知

## 环境变量

| 变量名 | 说明 | 必需 |
|--------|------|------|
| `INFISICAL_WEBHOOK_SECRET` | Infisical Webhook Secret，用于签名校验 | ✅ |
| `APPRISE_URL` | Apprise 服务地址（接收 JSON 的 API 端点） | ✅ |
| `NOTIFICATION_URLS` | Apprise 目标 URL 列表（按 Apprise 规范填写） | ✅ |

## 消息模板

- 模板文件：`message.md`
- 会注入 `SecretPath`，用于展示本次变更的路径

## 支持的事件

- `secrets.modified` - 密钥变更事件
- `test` - 测试事件

非法签名或无效负载会直接返回错误信息。

## 部署
1. 在 Appwrite 控制台创建一个新的 Function
2. 选择 Go 运行时环境
3. 上传该目录下的代码，设置 entrypoint 为 `main.go`
4. 配置上述环境变量
5. 在 Infisical 控制台设置 Webhook，指向该 Function 的 URL