# 技术设计: TODO 后端

## 技术方案
### 核心技术
- Go 1.23
- Gin
- GORM
- SQLite

### 实现要点
- 新增 `backend` Go 模块，独立运行 Gin 服务
- 统一读取环境变量:
  - `INFISICAL_WEBHOOK_SECRET`（签名校验）
  - `TODO_DB_PATH`（可选，默认 `backend/data/todos.db`）
- 签名校验逻辑复用 `notification/main.go` 的实现（时间窗 5 分钟）
- webhook 仅处理 `secrets.modified`，`test` 事件直接返回成功但不入库
- 数据模型：
  - `id` 自增主键
  - `secret_path` 唯一索引
  - `is_completed` 默认 false
  - `created_at` 自动写入
  - `completed_at` 可空
- webhook 幂等策略：以 `secret_path` 为唯一键，存在则复位完成状态，不改动 `created_at`

## API设计
### [POST] /api/todos/webhook
- **请求:** Infisical webhook payload
- **响应:** `{ "data": { "id": 1, "secretPath": "/xxx", "isCompleted": false, "createdAt": "...", "completedAt": null } }`

### [GET] /api/todos
- **响应:** `{ "data": [ ... ] }`

### [POST] /api/todos
- **请求:** `{ "secretPath": "/path" }`
- **响应:** `{ "data": { ... } }`

### [PATCH] /api/todos/{id}
- **请求:** `{ "secretPath": "/path" }`
- **响应:** `{ "data": { ... } }`

### [DELETE] /api/todos/{id}
- **响应:** `{ "data": true }`

### [POST] /api/todos/{id}/complete
- **响应:** `{ "data": { ... } }`

## 数据模型
```sql
CREATE TABLE todo_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  secret_path TEXT NOT NULL UNIQUE,
  is_completed BOOLEAN NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  completed_at DATETIME NULL
);
CREATE INDEX idx_todo_items_secret_path ON todo_items(secret_path);
```

## 安全与性能
- **安全:** 校验 `x-infisical-signature` 与时间戳，拒绝非法请求
- **性能:** SQLite 连接池设置合理上限，避免高并发写锁

## 测试与部署
- **测试:** 手动 curl 验证 webhook 与 CRUD（若用户要求再补充）
- **部署:** `go run ./backend` 本地启动，后续可接入服务管理器
