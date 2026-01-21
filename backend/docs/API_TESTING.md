# API 测试指南

本文档提供了在 Apifox 等 API 测试工具中测试本项目接口的方法。

## Webhook 接口测试

### 背景说明

Infisical Webhook 接口使用 HMAC-SHA256 签名验证机制，签名需要基于以下信息动态计算：

- **签名格式**: `t=<timestamp>,v1=<signature>`
- **计算方式**: `signature = HMAC-SHA256(secret, timestamp + "." + requestBody)`
- **密钥来源**: 环境变量 `INFISICAL_WEBHOOK_SECRET`

由于签名包含时间戳且每次请求都不同，因此无法使用固定的 Header 值进行测试。

### 在 Apifox 中测试

#### 步骤 1: 配置环境变量

在 Apifox 的环境设置中添加以下变量：

- `INFISICAL_WEBHOOK_SECRET`: 你的 Webhook 密钥（与后端环境变量中的值一致）

#### 步骤 2: 添加前置脚本

在 Webhook 接口的"前置操作"标签中，粘贴以下脚本：

```javascript
// 1. 获取 webhook secret
const secret = pm.environment.get("WEBHOOK_SECRET") || "your_secret_here";

// 2. 获取当前时间戳（秒）
const timestamp = Math.floor(Date.now() / 1000);

// 3. 获取请求体
const requestBody = pm.request.body.raw || "";

// 4. 构造签名字符串: timestamp.payload
const signedPayload = `${timestamp}.${requestBody}`;

// 5. 计算 HMAC-SHA256 签名
const signature = CryptoJS.HmacSHA256(signedPayload, secret);

// 6. 转换为 hex 格式
const signatureHex = signature.toString(CryptoJS.enc.Hex);

// 7. 构造签名头: t=<timestamp>;sha256=<signature>
const signatureHeader = `t=${timestamp};sha256=${signatureHex}`;

// 8. 设置到请求头
pm.request.headers.upsert({
    key: "x-infisical-signature",
    value: signatureHeader
});

console.log("Timestamp:", timestamp);
console.log("Signed Payload:", signedPayload);
console.log("Signature:", signatureHex);
console.log("Header:", signatureHeader);

```

#### 步骤 3: 发送请求

配置完成后，直接发送请求即可。前置脚本会自动：

1. 读取当前的请求 Body
2. 生成当前时间戳
3. 使用 HMAC-SHA256 计算签名
4. 自动设置 `X-Infisical-Signature` Header

### 测试用例示例

**请求示例:**

```json
POST /api/todos/webhook
Content-Type: application/json

{
  "event": "secret.created",
  "project": {
    "environment": "dev",
    "projectId": "test-project-123",
    "projectName": "Test Project",
    "secretName": "DATABASE_URL",
    "secretPath": "/config/database",
    "reminderNote": "测试提醒"
  },
  "timestamp": 1705912345
}
```

**预期响应:**

```json
{
    "data": "ok"
}
```

### 常见问题

#### Q: 签名验证失败 (401 错误)

**可能原因:**

1. `INFISICAL_WEBHOOK_SECRET` 环境变量未设置或与后端不一致
2. 前置脚本执行失败
3. 请求 Body 被意外修改（空格、换行等）

**解决方法:**

1. 检查 Apifox 环境变量配置
2. 查看前置脚本执行日志
3. 确保前置脚本在最终发送前执行

#### Q: 如何验证签名是否正确？

在前置脚本中添加 `console.log` 输出生成的签名和时间戳，然后在 Apifox 的控制台中查看。

## 其他接口测试

其他 Todo 相关接口（GET/POST/PATCH/DELETE）无需特殊配置，可直接测试。

### 示例：创建 Todo

```json
POST /api/todos
Content-Type: application/json

{
  "secretPath": "/config/database"
}
```

### 示例：获取所有 Todo

```
GET /api/todos
```

### 示例：标记完成

```
POST /api/todos/1/complete
```
