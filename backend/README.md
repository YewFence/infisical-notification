# Backend - Todo 管理后端服务

这是一个基于 Go 语言开发的 RESTful API 后端服务，用于管理待办事项（Todo），并支持通过 Webhook 接收 Infisical 的密钥变更通知。

## 📚 技术栈

- **语言**: Go 1.23.0
- **Web 框架**: [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- **ORM**: [GORM](https://gorm.io/) - Go 语言的 ORM 库
- **数据库**: SQLite - 轻量级嵌入式数据库
- **架构模式**: 分层架构（Router → Handler → Repository → Database）

## 📁 项目结构

```
backend/
├── main.go                 # 程序入口，负责初始化和启动服务
├── go.mod                  # Go 模块依赖管理
├── data/                   # 数据库文件存储目录
│   └── todos.db           # SQLite 数据库文件（自动生成）
└── internal/              # 内部包，遵循 Go 项目规范
    ├── config/            # 配置管理
    │   └── config.go      # 从环境变量加载配置
    ├── db/                # 数据库连接
    │   └── sqlite.go      # SQLite 初始化逻辑
    ├── handlers/          # HTTP 处理器（Controller 层）
    │   ├── response.go    # 统一响应格式
    │   ├── todos.go       # Todo 相关接口
    │   └── webhook.go     # Webhook 接口
    ├── models/            # 数据模型（Model 层）
    │   └── todo.go        # TodoItem 结构体定义
    ├── repo/              # 数据访问层（Repository 层）
    │   └── todo_repo.go   # Todo 数据操作封装
    ├── router/            # 路由配置
    │   └── router.go      # HTTP 路由注册
    └── signature/         # 签名验证
        └── verify.go      # Infisical Webhook 签名验证
```

### 分层架构说明

```
┌─────────────┐
│   Router    │  路由层：定义 URL 路径和 HTTP 方法
└──────┬──────┘
       │
┌──────▼──────┐
│   Handler   │  处理器层：处理请求、验证参数、调用业务逻辑
└──────┬──────┘
       │
┌──────▼──────┐
│ Repository  │  数据访问层：封装数据库操作
└──────┬──────┘
       │
┌──────▼──────┐
│  Database   │  数据库：SQLite 持久化存储
└─────────────┘
```

## 🔌 API 接口文档

### 基础信息

- **Base URL**: `http://localhost:8080/api/todos`
- **数据格式**: JSON
- **时间格式**: RFC3339（例如：`2024-01-20T14:30:00Z`）

### 响应格式

**成功响应**:
```json
{
  "data": {
    // 响应数据
  }
}
```

**错误响应**:
```json
{
  "error": "错误信息"
}
```

### 接口列表

#### 1. 获取待办事项列表

**请求**:
```
GET /api/todos
```

**响应示例**:
```json
{
  "data": [
    {
      "id": 1,
      "secretPath": "/dev/db/password",
      "isCompleted": false,
      "createdAt": "2024-01-20T10:00:00Z",
      "completedAt": null
    },
    {
      "id": 2,
      "secretPath": "/prod/api/key",
      "isCompleted": true,
      "createdAt": "2024-01-19T15:30:00Z",
      "completedAt": "2024-01-20T08:45:00Z"
    }
  ]
}
```

#### 2. 创建待办事项

**请求**:
```
POST /api/todos
Content-Type: application/json

{
  "secretPath": "/dev/db/password"
}
```

**响应示例**:
```json
{
  "data": {
    "id": 3,
    "secretPath": "/dev/db/password",
    "isCompleted": false,
    "createdAt": "2024-01-20T14:30:00Z",
    "completedAt": null
  }
}
```

**错误码**:
- `400 Bad Request`: 请求参数错误
- `409 Conflict`: `secretPath` 已存在（唯一约束）

#### 3. 更新待办事项

**请求**:
```
PATCH /api/todos/:id
Content-Type: application/json

{
  "secretPath": "/dev/db/new-password"
}
```

**响应示例**:
```json
{
  "data": {
    "id": 1,
    "secretPath": "/dev/db/new-password",
    "isCompleted": false,
    "createdAt": "2024-01-20T10:00:00Z",
    "completedAt": null
  }
}
```

**错误码**:
- `400 Bad Request`: 请求参数错误
- `404 Not Found`: 待办事项不存在
- `409 Conflict`: 新的 `secretPath` 已存在

#### 4. 删除待办事项

**请求**:
```
DELETE /api/todos/:id
```

**响应示例**:
```json
{
  "data": true
}
```

**错误码**:
- `404 Not Found`: 待办事项不存在

#### 5. 标记为完成

**请求**:
```
POST /api/todos/:id/complete
```

**响应示例**:
```json
{
  "data": {
    "id": 1,
    "secretPath": "/dev/db/password",
    "isCompleted": true,
    "createdAt": "2024-01-20T10:00:00Z",
    "completedAt": "2024-01-20T16:00:00Z"
  }
}
```

**错误码**:
- `404 Not Found`: 待办事项不存在

#### 6. Webhook 接口

接收来自 Infisical 的密钥变更通知。

**请求**:
```
POST /api/todos/webhook
Content-Type: application/json
x-infisical-signature: t=1234567890;sha256=abc123...

{
  "event": "secrets.modified",
  "project": {
    "secretPath": "/dev/db/password",
    "projectId": "proj_123",
    "projectName": "MyProject",
    "environment": "dev"
  },
  "timestamp": 1234567890
}
```

**支持的事件类型**:
- `secrets.modified`: 密钥被修改
- `test`: 测试事件

**行为**:
- 如果 `secretPath` 已存在且未完成 → 标记为完成
- 如果 `secretPath` 不存在 → 创建新的待办事项

**安全验证**:
- 请求必须包含有效的 `x-infisical-signature` 签名头
- 签名使用 HMAC-SHA256 算法，基于配置的 `INFISICAL_WEBHOOK_SECRET` 进行验证

**错误码**:
- `400 Bad Request`: 请求体格式错误或缺失签名
- `401 Unauthorized`: 签名验证失败
- `500 Internal Server Error`: 服务端处理失败

## 🚀 开发说明

### 环境要求

- Go 1.23.0 或更高版本
- 操作系统：Windows / Linux / macOS

### 安装依赖

```bash
cd backend
go mod download
```

### 配置环境变量

创建 `.env` 文件或直接设置环境变量（可选，程序有默认值）：

```bash
# Webhook 签名验证密钥（必需，如果要使用 Webhook 功能）
INFISICAL_WEBHOOK_SECRET=your_secret_key_here

# 数据库文件路径（可选，默认：backend/data/todos.db）
TODO_DB_PATH=./data/todos.db

# HTTP 服务监听地址（可选，默认：:8080）
TODO_BIND_ADDR=:8080
```

**PowerShell 设置环境变量示例**:
```powershell
$env:INFISICAL_WEBHOOK_SECRET="your_secret_key_here"
$env:TODO_BIND_ADDR=":8080"
```

**Bash 设置环境变量示例**:
```bash
export INFISICAL_WEBHOOK_SECRET="your_secret_key_here"
export TODO_BIND_ADDR=":8080"
```

### 运行服务

#### 开发模式
```bash
cd backend
go run main.go
```

#### 编译并运行
```bash
cd backend
go build -o todo-server
./todo-server          # Linux/macOS
# 或
.\todo-server.exe      # Windows
```

#### 热重载开发（推荐）

安装 [air](https://github.com/cosmtrek/air)：
```bash
go install github.com/cosmtrek/air@latest
```

运行：
```bash
cd backend
air
```

### 测试接口

使用 curl 测试：

```bash
# 创建待办事项
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"secretPath": "/dev/test"}'

# 获取列表
curl http://localhost:8080/api/todos

# 标记完成
curl -X POST http://localhost:8080/api/todos/1/complete

# 删除
curl -X DELETE http://localhost:8080/api/todos/1
```

使用 PowerShell 测试：

```powershell
# 创建待办事项
Invoke-RestMethod -Uri "http://localhost:8080/api/todos" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"secretPath": "/dev/test"}'

# 获取列表
Invoke-RestMethod -Uri "http://localhost:8080/api/todos"
```

### 代码格式化

```bash
# 格式化代码
go fmt ./...

# 运行 linter
go vet ./...
```

### 项目构建

```bash
# 为当前平台构建
go build -o bin/todo-server main.go

# 跨平台构建示例
GOOS=linux GOARCH=amd64 go build -o bin/todo-server-linux main.go
GOOS=windows GOARCH=amd64 go build -o bin/todo-server.exe main.go
GOOS=darwin GOARCH=amd64 go build -o bin/todo-server-mac main.go
```

## 🔒 安全说明

### Webhook 签名验证

本服务实现了严格的 Webhook 签名验证机制，确保只有 Infisical 能够触发回调：

1. **签名算法**: HMAC-SHA256
2. **签名头格式**: `x-infisical-signature: t=<timestamp>;<signature>`
3. **签名计算**: 对 `<timestamp>.<payload>` 格式的字符串进行 HMAC-SHA256 计算
4. **验证步骤**:
   - 提取时间戳和签名
   - 检查时间戳是否在允许范围内（防止重放攻击）
   - 使用配置的密钥重新计算 `timestamp.payload` 的签名
   - 比较计算结果与请求中的签名是否一致

### 数据库安全

- SQLite 数据库文件默认存储在 `backend/data/` 目录
- 确保数据库文件权限设置正确，避免未授权访问
- 生产环境建议定期备份数据库文件

## 📝 数据模型

### TodoItem

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| `id` | `uint` | 主键 | 自增 |
| `secret_path` | `string` | 密钥路径 | 唯一、非空 |
| `is_completed` | `bool` | 是否已完成 | 非空、默认 false |
| `created_at` | `time.Time` | 创建时间 | 非空、自动填充 |
| `completed_at` | `*time.Time` | 完成时间 | 可为空 |

## 🐛 故障排查

### 问题：端口已被占用

```
Error: listen tcp :8080: bind: address already in use
```

**解决方法**:
- 修改 `TODO_BIND_ADDR` 环境变量使用其他端口
- 或关闭占用 8080 端口的进程

### 问题：数据库文件权限错误

```
Error: unable to open database file
```

**解决方法**:
- 确保 `data/` 目录存在且有写权限
- Windows: 检查文件夹的安全设置
- Linux/macOS: `chmod 755 data/`

### 问题：Webhook 签名验证失败

```
Error: invalid signature
```

**解决方法**:
- 检查 `INFISICAL_WEBHOOK_SECRET` 是否与 Infisical 配置一致
- 确保请求头中包含正确的 `x-infisical-signature`
- 检查时间戳是否在有效范围内（避免时钟偏移）

## 📖 相关资源

- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [Infisical Webhook 文档](https://infisical.com/docs/integrations/webhooks)
- [Go 语言官方文档](https://go.dev/doc/)

## 📄 许可证

本项目遵循 MIT 许可证。
