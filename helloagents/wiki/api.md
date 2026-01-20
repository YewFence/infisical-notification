# API 手册

## 概述
当前暂无统一对外 API 文档。notification 模块通过 Appwrite Function 接收 Infisical webhook。

## 认证方式
- **notification:** 使用 `x-infisical-signature` 进行签名校验
- **backend:** webhook 使用 `x-infisical-signature` 校验，CRUD 当前未加入鉴权（默认内网/本地使用）

---

## 接口列表

### notification

#### [POST] /notification
**描述:** 接收 Infisical webhook（Appwrite Function 入口）。

### backend

#### [POST] /api/todos/webhook
**描述:** 接收 Infisical webhook 并写入/更新 TODO。
**请求头:** `x-infisical-signature`
**响应:**
```json
{ "data": { "id": 1, "secretPath": "/app", "isCompleted": false, "createdAt": "2026-01-20T20:00:00Z", "completedAt": null } }
```

#### [GET] /api/todos
**描述:** 获取 TODO 列表。
**响应:**
```json
{ "data": [ { "id": 1, "secretPath": "/app", "isCompleted": false, "createdAt": "2026-01-20T20:00:00Z", "completedAt": null } ] }
```

#### [POST] /api/todos
**描述:** 创建 TODO。
**请求:**
```json
{ "secretPath": "/app" }
```
**响应:**
```json
{ "data": { "id": 1, "secretPath": "/app", "isCompleted": false, "createdAt": "2026-01-20T20:00:00Z", "completedAt": null } }
```

#### [PATCH] /api/todos/{id}
**描述:** 更新 TODO。
**请求:**
```json
{ "secretPath": "/new-path" }
```
**响应:**
```json
{ "data": { "id": 1, "secretPath": "/new-path", "isCompleted": false, "createdAt": "2026-01-20T20:00:00Z", "completedAt": null } }
```

#### [DELETE] /api/todos/{id}
**描述:** 删除 TODO。
**响应:**
```json
{ "data": true }
```

#### [POST] /api/todos/{id}/complete
**描述:** 完成 TODO。
**响应:**
```json
{ "data": { "id": 1, "secretPath": "/app", "isCompleted": true, "createdAt": "2026-01-20T20:00:00Z", "completedAt": "2026-01-20T20:30:00Z" } }
```
