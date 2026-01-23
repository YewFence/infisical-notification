# Frontend 代码审查待办事项

> 审查日期: 2026-01-23

## P1 中优先级

### 1. 删除回滚位置错误
**文件**: `App.tsx:136-138`

删除失败时，任务被添加到列表末尾而不是原来的位置，导致用户界面混乱。

```typescript
// 当前代码
setTasks(prev => [...prev, taskToDelete]);

// 建议修复：保存原始索引，回滚时恢复到正确位置
const taskIndex = tasks.findIndex(t => t.id === id);
// ... 在 catch 中
setTasks(prev => {
  const newTasks = [...prev];
  newTasks.splice(taskIndex, 0, taskToDelete);
  return newTasks;
});
```

---

### 2. 使用 alert() 显示错误
**文件**: `App.tsx:102`, `App.tsx:139`, `App.tsx:158`

使用原生 `alert()` 阻塞用户界面，体验不佳，与应用的现代化 UI 风格不一致。

**建议**: 实现 Toast 通知组件或使用现有的通知库。

---

### 3. API 响应缺乏运行时验证
**文件**: `services/api.ts:30-38`

使用 `as` 类型断言而没有运行时验证，如果后端返回格式不符合预期，可能导致运行时错误。

**建议**: 添加运行时类型检查或使用 `zod` 等库进行验证。

---

## P2 低优先级

### 4. 语言混用
界面文本中英文混用，不一致：
- `CreateTaskModal.tsx:51` - 中文："创建新的任务"
- `CreateTaskModal.tsx:64` - 中文："Secret 路径"
- `App.tsx:207` - 英文："Project Tasks"
- `App.tsx:235` - 英文："Search by task..."
- `TaskTable.tsx:85` - 英文："No tasks found"

**建议**: 统一使用一种语言，或实现 i18n 国际化。

---

### 5. 硬编码的最小宽度
**文件**: `TaskTable.tsx:68`

```typescript
<div className="min-w-[900px] border-t border-border">
```

硬编码 900px 最小宽度，在小屏幕设备上可能导致水平滚动。

---

### 6. 日期格式化函数重复定义
**文件**: `services/adapter.ts:5`

`formatDate` 函数定义在 `todoItemToTask` 内部，每次调用都会重新创建。

**建议**: 将其提取为模块级函数。

---

### 7. CDN 依赖无完整性校验
**文件**: `index.html:8`

```html
<script src="https://cdn.tailwindcss.com"></script>
```

使用 CDN 加载 Tailwind CSS 没有 `integrity` 属性进行 SRI 校验。

**建议**:
- 生产环境应使用本地构建的 Tailwind CSS
- 或添加 `integrity` 和 `crossorigin` 属性

---

## P3 建议改进

### 8. 添加单元测试
当前项目没有测试文件，建议为关键业务逻辑添加测试。

### 9. 添加 ESLint 配置
自动检测未使用的导入和变量。

### 10. 添加错误边界
使用 React Error Boundary 捕获渲染错误。

### 11. 实现加载骨架屏
替代简单的 "加载中..." 文本，提升用户体验。

---

## 已完成 ✅

- [x] 移除前端 API Key 暴露 (`vite.config.ts`)
- [x] 移除未使用的 `@google/genai` 依赖 (`index.html`)
- [x] 移除未使用的 `Filter` 导入 (`App.tsx`)
- [x] 清理所有注释掉的死代码
- [x] 移除无功能的 ArrowRight 按钮
- [x] 移除未实现的分页组件
