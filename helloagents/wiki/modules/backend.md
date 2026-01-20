# backend

## 目的
提供 TODO 后端服务与 API，支撑前端清单展示与 webhook 入库。

## 模块概述
- **职责:** CRUD + 完成状态管理，接收 webhook 写入/更新
- **状态:** 🚧开发中
- **最后更新:** 2026-01-20

## 规范

### 需求: TODO 后端
**模块:** backend
提供 webhook 入库、CRUD、完成状态接口。

#### 场景: webhook 入库
收到合法 webhook。
- 预期结果: 新增/更新 TODO 条目

#### 场景: 完成 TODO
用户执行完成操作。
- 预期结果: 标记完成并记录完成时间

## API接口
### [POST] /api/todos/webhook
**描述:** 接收 Infisical webhook 并写入/更新 TODO。
**请求头:** `x-infisical-signature`

### [GET] /api/todos
**描述:** 获取 TODO 列表。

### [POST] /api/todos
**描述:** 创建 TODO。

### [PATCH] /api/todos/{id}
**描述:** 更新 TODO。

### [DELETE] /api/todos/{id}
**描述:** 删除 TODO。

### [POST] /api/todos/{id}/complete
**描述:** 完成 TODO。

## 数据模型
### todo_items
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 自增主键 |
| secret_path | TEXT | secretPath |
| is_completed | BOOLEAN | 是否完成 |
| created_at | DATETIME | 创建时间 |
| completed_at | DATETIME | 完成时间 |

## 依赖
- SQLite
- Infisical webhook（签名校验规则与 notification 保持一致）

## 变更历史
- [202601202042_todo-backend](../../history/2026-01/202601202042_todo-backend/) - 新增 TODO 后端



