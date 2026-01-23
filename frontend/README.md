# YewTally 前端

一个简洁的任务管理 Web 应用，基于 React 19 + TypeScript + Tailwind CSS v4 构建。

## 技术栈

- **React 19** - UI 框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Tailwind CSS v4** - 样式方案
- **Lucide React** - 图标库
- **Zod** - 运行时数据校验

## 功能特性

- 任务列表展示与搜索
- 创建/删除任务
- 任务状态切换（todo ↔ done）
- 自动轮询刷新
- 乐观更新 + 错误回滚
- Toast 通知提示
- 骨架屏加载状态

## 本地开发

**前置条件:** Node.js 18+，推荐使用 pnpm

1. 安装依赖：
   ```bash
   pnpm install
   ```

2. 配置环境变量，复制 `.env.example` 到 `.env.local`：
   ```bash
   cp .env.example .env.local
   ```

3. 启动开发服务器：
   ```bash
   pnpm dev
   ```

## 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `VITE_API_BASE_URL` | 后端 API 地址 | `http://localhost:8080/api/todos` |
| `VITE_POLL_INTERVAL_SECONDS` | 轮询间隔（秒） | `30` |
| `NGINX_PORT` | Nginx 监听端口（Docker 部署） | `5473` |

## 构建部署

```bash
pnpm build
```

构建产物位于 `dist/` 目录，可直接部署到静态服务器或通过 nginx 代理。
