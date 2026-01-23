# infisical-notification

> 本文件包含项目级别的核心信息。详细的模块文档见 modules/ 目录。

---

## 1. 项目概述

### 目标与背景
接收 Infisical webhook，校验签名并转发通知；新增 TODO 后端用于前端清单展示。

### 范围
- **范围内:** webhook 接收与校验、消息转发、TODO 数据落库与查询
- **范围外:** 复杂权限体系、分页与全文检索

### 干系人
- **负责人:** 叶云枫

---

## 2. 模块索引

| 模块名称 | 职责 | 状态 | 文档 |
|---------|------|------|------|
| notification | 处理 Infisical webhook 与通知转发 | ✅稳定 | [notification](modules/notification.md) |
| backend | TODO 后端服务与 API | 🚧开发中 | [backend](modules/backend.md) |
| frontend | TODO 列表展示界面 | 📝规划中 | [frontend](modules/frontend.md) |

---

## 3. 快速链接
- [技术约定](../project.md)
- [架构设计](arch.md)
- [API 手册](api.md)
- [数据模型](data.md)
- [变更历史](../history/index.md)


