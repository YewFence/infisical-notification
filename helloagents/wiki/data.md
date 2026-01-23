# 数据模型

## 概述
TODO 后端已引入 SQLite 存储待办数据。

---

## 数据表/集合

### todo_items

**描述:** TODO 清单条目。

| 字段名 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | INTEGER | 主键/自增 | 递增主键 |
| secret_path | TEXT | 非空 | 对应 webhook 的 secretPath |
| is_completed | BOOLEAN | 非空 | 是否完成 |
| created_at | DATETIME | 非空 | 创建时间 |
| completed_at | DATETIME | 可空 | 完成时间 |

**索引:**
- idx_todo_items_secret_path: secret_path

