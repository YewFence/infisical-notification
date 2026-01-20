# 任务清单: TODO 后端

目录: `helloagents/history/2026-01/202601202042_todo-backend/`

---

## 1. 基础服务
- [√] 1.1 在 `backend/go.mod`、`backend/main.go`、`backend/internal/router/router.go` 初始化 Gin 服务与路由骨架，验证 why.md#需求-webhook-入库-场景-签名合法
- [√] 1.2 在 `backend/internal/config/config.go` 读取环境变量并在 `backend/main.go` 注入配置，验证 why.md#需求-todo-管理-场景-查看与维护列表

## 2. 数据持久化
- [√] 2.1 在 `backend/internal/models/todo.go` 定义模型，在 `backend/internal/db/sqlite.go` 初始化 SQLite，并在 `backend/main.go` 完成自动迁移，验证 why.md#需求-webhook-入库-场景-签名合法
- [√] 2.2 在 `backend/internal/repo/todo_repo.go` 实现按 `secret_path` 的查询与 upsert，验证 why.md#需求-webhook-入库-场景-签名合法

## 3. Webhook 接入
- [√] 3.1 在 `backend/internal/signature/verify.go` 复刻签名校验逻辑，在 `backend/internal/handlers/webhook.go` 处理 payload 并入库，验证 why.md#需求-webhook-入库-场景-签名合法
- [√] 3.2 在 `backend/internal/router/router.go` 注册 `/api/todos/webhook`，验证 why.md#需求-webhook-入库-场景-签名合法

## 4. CRUD 与完成操作
- [√] 4.1 在 `backend/internal/handlers/todos.go` 实现列表/新增/更新/删除/完成接口，验证 why.md#需求-todo-管理-场景-查看与维护列表
- [√] 4.2 在 `backend/internal/router/router.go` 注册 CRUD 路由，验证 why.md#需求-todo-管理-场景-查看与维护列表
- [√] 4.3 在 `backend/internal/repo/todo_repo.go` 支持完成/删除操作，验证 why.md#需求-todo-管理-场景-完成-todo

## 5. 安全检查
- [√] 5.1 执行安全检查（按G9: 输入验证、敏感信息处理、权限控制、EHRB风险规避）

## 6. 文档更新
- [√] 6.1 更新 `helloagents/wiki/api.md`
- [√] 6.2 更新 `helloagents/wiki/data.md`
- [√] 6.3 更新 `helloagents/wiki/modules/backend.md`
- [√] 6.4 更新 `helloagents/CHANGELOG.md`

## 7. 测试
- [-] 7.1 若用户要求，补充手动验证步骤（curl）
> 备注: 用户未要求验证。








