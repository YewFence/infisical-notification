# notification

## 目的
接收 Infisical webhook，校验签名并向 Apprise 转发通知。

## 模块概述
- **职责:** 解析 webhook、校验签名、过滤事件、渲染消息、发送通知
- **状态:** ✅稳定
- **最后更新:** 2026-01-20

## 规范

### 需求: 处理 Infisical webhook
**模块:** notification
处理 secrets.modified 与 	est 事件，校验 x-infisical-signature。

#### 场景: 合法请求
签名正确且时间戳在窗口内。
- 预期结果: 事件被解析并推送到 Apprise

#### 场景: 非法请求
签名错误或 payload 无效。
- 预期结果: 返回错误文本并记录日志

## API接口
### [POST] /notification
**描述:** Appwrite Function webhook 入口。
**输入:** Infisical webhook payload
**输出:** 文本响应

## 数据模型
无持久化数据。

## 依赖
- Apprise 服务
- Infisical webhook

## 变更历史
- 初始化记录待补充
