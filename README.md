# GoBox

GoBox 是一个基于 Go + React 的在线工具箱平台 MVP。当前实现包含 `Gin + GORM + MySQL + Redis + React + Vite` 的单仓结构，支持工具执行、登录注册、邮箱验证码注册、统计分析、用户偏好以及 Docker / Nginx / GitHub Actions 配套。

## 项目结构

```text
GoBox/
├── backend/                  # Gin + GORM API
├── frontend/                 # React + Vite 前端
├── docker-compose.yml        # 本地容器编排
├── docker-compose.prod.yml   # 生产容器编排
├── nginx.conf                # 反向代理
├── .github/workflows/ci.yml  # CI
└── Makefile                  # 常用命令
```

## 已实现能力

- 用户登录、JWT 鉴权
- 邮箱验证码注册
- 用户中心与偏好保存
- 工具目录查询、工具运行、工具调用统计
- 管理台汇总数据展示
- 内置 12 个常用开发工具
- MySQL 持久化、Redis 预留、Nginx 代理

## 本地启动

### Docker 方式

```bash
docker compose up --build
```

启动后访问：

- 前端首页：`http://localhost/`
- 工具页：`http://localhost/tools`
- 用户中心：`http://localhost/user`
- 管理页：`http://localhost/admin`

默认管理员：

- 邮箱：`admin@gobox.local`
- 密码：`admin123456`

本地编排默认启动：

- `frontend`
- `backend`
- `nginx`
- `mysql`
- `redis`

MySQL 默认参数：

- 数据库：`gobox`
- 用户名：`gobox`
- 密码：`gobox123`
- Root 密码：`root123456`

首次启动时，MySQL 会自动执行这些初始化脚本：

- [backend/scripts/mysql/01-schema.sql](/d:/godemo/GoBox/backend/scripts/mysql/01-schema.sql:1)
- [backend/scripts/mysql/02-seed.sql](/d:/godemo/GoBox/backend/scripts/mysql/02-seed.sql:1)

### 手动方式

先准备一个 MySQL 8 数据库并创建 `gobox` 库，然后配置 [backend/.env.example](/d:/godemo/GoBox/backend/.env.example:1)：

```env
DATABASE_DRIVER=mysql
DATABASE_DSN=gobox:gobox123@tcp(127.0.0.1:3306)/gobox?charset=utf8mb4&parseTime=True&loc=Local
```

说明：

- 直接在本机运行 `backend` 时，数据库主机应使用 `127.0.0.1`
- 通过 `docker compose` 运行时，backend 会被编排文件覆盖为 `mysql:3306`

启动后端：

```bash
cd backend
cp .env.example .env
go mod tidy
go run ./cmd/server
```

启动前端：

```bash
cd frontend
npm install
npm run dev
```

## 邮箱验证码注册

注册流程已经改成：

1. 输入邮箱
2. 点击“发送验证码”
3. 邮箱收到验证码后输入
4. 校验通过才会真正创建用户

相关接口：

- `POST /api/v1/auth/register/send-code`
- `POST /api/v1/auth/register`

开发环境如果没有配置 SMTP，不会真实发邮件，后端会返回 `previewCode` 方便本地联调。生产环境下如果未配置 SMTP，发送验证码会直接失败。

### SMTP 配置

见 [backend/.env.example](/d:/godemo/GoBox/backend/.env.example:1) 和 [backend/.env.production](/d:/godemo/GoBox/backend/.env.production:1)：

```env
SMTP_ENABLED=true
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your-smtp-username
SMTP_PASSWORD=your-smtp-password
SMTP_FROM=GoBox <noreply@example.com>
```

## 云服务器部署

生产部署前先修改：

- [backend/.env.production](/d:/godemo/GoBox/backend/.env.production:1)

至少检查这些值：

- `DATABASE_DSN`
- `AUTH_JWT_SECRET`
- `SMTP_ENABLED`
- `SMTP_HOST`
- `SMTP_PORT`
- `SMTP_USERNAME`
- `SMTP_PASSWORD`
- `SMTP_FROM`

然后启动：

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

放行服务器 `80` 端口后，直接访问服务器 IP 或域名即可。

## 说明

- 默认实现已经切到 `MySQL`
- MySQL 首次建库建表不再只依赖 `GORM AutoMigrate`，已经补了显式 SQL 初始化文件
- 代码里仍保留 `sqlite` fallback，仅用于特殊本地调试场景
- 前端容器已支持 SPA 路由回退，直接访问 `/tools`、`/user`、`/admin` 不会 404
