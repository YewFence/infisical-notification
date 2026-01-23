# Infisical Notification

一个轻量级的 Todo 清单应用，可接收 Infisical Webhook 自动创建待办事项。当密钥变更时，自动提醒你需要同步更新哪些服务，确保不会遗漏任何操作。

## 项目结构

```
infisical-notification/
├── backend/          # Go 后端服务
├── frontend/         # React 前端应用
├── notification/     # Appwrite Function 通知转发器（计划整合）
└── docs/             # API 文档
```

## 组件说明

### Backend

基于 Go + Gin + GORM 的 RESTful API 服务，负责：
- 接收 Infisical Webhook 并验证签名
- 管理待办事项（Todo）的 CRUD 操作
- 提供 Swagger API 文档

👉 详细信息请查看 [backend/README.md](./backend/README.md)

### Frontend

基于 React 19 + TypeScript + Tailwind CSS v4 的 Web 应用，功能包括：
- 任务列表展示与搜索
- 任务状态管理
- 乐观更新与 Toast 通知

👉 详细信息请查看 [frontend/README.md](./frontend/README.md)

### Notification（计划整合）

基于 Appwrite Function 的通知转发器，负责：
- 接收 Infisical Webhook 并校验签名
- 将密钥变更事件转发到 Apprise 通知服务

👉 详细信息请查看 [notification/README.md](./notification/README.md)

## 安全注意事项

> ⚠️ **重要提示**：本项目的后端 API 不包含身份验证机制，任何能访问服务的人都可以增删改查待办事项。

**推荐的部署方式：**

1. **Tailscale / ZeroTier 等 VPN**：将服务部署在私有网络中，仅允许 VPN 内的设备访问
2. **Cloudflare Access / Authentik**：在反向代理层添加身份验证，保护整个应用
3. **内网部署**：仅在本地局域网中运行，不暴露到公网

**关于 Webhook 端点：**

`/api/webhook` 端点需要公开暴露以接收 Infisical 的回调。该端点已实现签名验证，只接受携带正确签名的请求，因此即使暴露也是安全的。

如需公网部署，请确保：
- 配置正确的 `INFISICAL_WEBHOOK_SECRET` 环境变量
- 使用 HTTPS（可通过 Cloudflare 等服务实现）

## 快速开始

### 本地开发

1. **启动后端**
   ```bash
   cd backend
   go run main.go
   ```

2. **启动前端**
   ```bash
   cd frontend
   pnpm install && pnpm dev
   ```

详细配置和开发说明请参阅各组件的 README。

### Docker 部署

#### 本地构建测试

```bash
docker compose -f compose.dev.yaml up --build
```

访问 http://localhost:${TODO_BIND_ADDR} 即可。

#### 生产部署

1. 创建 `compose.yaml`：

可以参考 [compose.yaml](./compose.prod.example) 文件。

2. 创建数据目录：

```bash
mkdir data
```

3. 配置环境变量：

创建 `.env` 文件，参考 [.env.prod.example](./.env.prod.example)。

4. 拉取并启动：

```bash
docker compose pull && docker compose up -d
```

#### 镜像版本

| 镜像 | 说明 |
|------|------|
| `ghcr.io/yewfence/infisical-notification-frontend:latest` | 前端最新版 |
| `ghcr.io/yewfence/infisical-notification-backend:latest` | 后端最新版 |
| `...:1.0.0` | 指定版本号 |
| `...:main` | main 分支构建 |
| `...:feature-web` | feature/web 分支构建 |

## CI/CD

项目使用 GitHub Actions 自动构建 Docker 镜像：

| 工作流 | 触发条件 | 说明 |
|-------|---------|------|
| Build Frontend | push main (frontend/**) / 手动 | 构建前端镜像 |
| Build Backend | push main (backend/**) / 手动 | 构建后端镜像 |
| Build All | 手动 | 同时构建前后端 |
| Release | 推送 v* 标签 | 构建镜像 + 创建 GitHub Release |

**发布新版本：**

```bash
git tag v1.0.0
git push origin v1.0.0
```

## 许可证

MIT
