Infisical 通知转发器（Appwrite）

这是一个很小巧的 Appwrite Function：接收 Infisical 的 webhook，校验签名后只处理 `secrets.modified` 事件，并把消息转发到你自定义的 Apprise。

主要流程
- 接收 Infisical webhook 请求
- 校验 `x-infisical-signature` 签名与时间窗
- 解析事件与路径，渲染消息模板
- 调用 Apprise API 推送通知

环境变量（必须）
- `INFISICAL_WEBHOOK_SECRET`: Infisical Webhook Secret，用于签名校验
- `APPRISE_URL`: Apprise 服务地址（接收 JSON 的 API 端点）
- `NOTIFICATION_URLS`: Apprise 目标 URL 列表（按 Apprise 规范填写）

消息模板
- 模板文件：`notification/message.md`
- 会注入 `SecretPath`，用于展示本次变更的路径

备注
- 只处理事件：`secrets.modified`
- 非法签名或无效负载会直接返回错误信息
