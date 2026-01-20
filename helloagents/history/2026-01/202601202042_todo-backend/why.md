# 变更提案: TODO 后端

## 需求背景
当前仅有 `notification` Appwrite Function 处理 Infisical webhook 并推送通知，前端还无法获取可展示的 TODO 清单。需要新增一个 TODO 后端服务：接收 webhook、落库、提供 CRUD 与完成操作接口，支撑前端展示与管理。

## 产品分析

### 目标用户与场景
- **用户群体:** 叶云枫及团队成员
- **使用场景:** Infisical 触发 webhook 后，前端可查看与处理待办清单
- **核心痛点:** 没有可查询的 TODO 数据源，无法追踪变更事项

### 价值主张与成功指标
- **价值主张:** 将 webhook 事件转为可管理的 TODO
- **成功指标:** webhook 可入库、前端可拉取 TODO 列表、完成状态可更新

### 人文关怀
仅保存必要字段（不存储密钥内容），避免敏感信息泄露；日志避免输出完整 payload。

## 变更内容
1. 新增 `backend` Gin 服务，提供 webhook 与 CRUD 接口
2. 引入 SQLite + GORM 持久化 TODO 数据
3. 复用 `notification/main.go` 的签名校验逻辑，实现一致的安全校验

## 影响范围
- **模块:** backend
- **文件:** `backend/go.mod`, `backend/main.go`, `backend/internal/*`
- **API:** `/api/todos/*`
- **数据:** `todo_items` 表（SQLite）

## 核心场景

### 需求: Webhook 入库
**模块:** backend
接收 Infisical webhook，校验签名并写入/更新 TODO。

#### 场景: 签名合法
payload 为 `secrets.modified`，且签名与时间窗校验通过。
- 预期结果: 新增或更新 `secretPath` 对应 TODO，未完成状态，`completed_at` 为空

#### 场景: 签名非法
签名错误或时间戳过期。
- 预期结果: 返回错误响应，不写入数据库

### 需求: TODO 管理
**模块:** backend
提供 CRUD 与完成操作接口。

#### 场景: 查看与维护列表
前端发起列表/新增/更新/删除请求。
- 预期结果: 返回 JSON 数据，状态与字段正确

#### 场景: 完成 TODO
用户标记 TODO 为完成。
- 预期结果: `is_completed=true` 且写入 `completed_at`

## 风险评估
- **风险:** webhook 签名校验不一致导致误拒绝
  - **缓解:** 复用 `notification` 的时间窗与算法逻辑
- **风险:** SQLite 并发写入引发锁等待
  - **缓解:** 使用单连接池配置与合理超时
- **风险:** CRUD 接口缺少鉴权
  - **缓解:** 先限定本地/内网使用，后续补充鉴权
