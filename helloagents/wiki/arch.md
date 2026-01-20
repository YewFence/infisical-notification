# 架构设计

## 总体架构
```mermaid
flowchart TD
    A[Infisical Webhook] --> B[notification 函数]
    B --> C[Apprise 服务]
    A --> D[TODO 后端服务]
    D --> E[SQLite]
    D --> F[前端 TODO 列表]
```

## 技术栈
- **后端:** Go 1.23 / Gin + GORM
- **前端:** TBD
- **数据:** SQLite

## 核心流程
```mermaid
sequenceDiagram
    participant Infisical
    participant Notification
    participant Apprise
    Infisical->>Notification: webhook
    Notification->>Notification: 校验签名/解析事件
    Notification->>Apprise: 推送通知
```

## 重大架构决策
完整的ADR存储在各变更的how.md中，本章节提供索引。

| adr_id | title | date | status | affected_modules | details |
|--------|-------|------|--------|------------------|---------|
| ADR-001 | 引入 Gin + SQLite 作为 TODO 后端 | 2026-01-20 | ✅已采纳 | backend | helloagents/history/2026-01/202601202042_todo-backend/how.md |

