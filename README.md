# Web Tracing Admin

基于 web-tracing-sdk 的错误监控后台，用于收集、分析和展示前端错误信息。

## 功能特点

- 错误监控：收集和展示前端应用的错误信息
- 性能监控：分析页面加载性能和资源加载情况
- 请求监控：跟踪 HTTP 请求的成功率和响应时间
- 用户行为：记录用户点击和页面浏览行为
- 资源监控：监控静态资源加载情况
- 路由监控：跟踪页面路由变化
- 曝光监控：记录元素曝光情况

## 技术栈

### 后端

- Go
- Gin 框架
- GORM ORM
- MySQL 数据库
- Air (热重载)
- Swagger (API 文档)

### 前端

- Vue 3
- TypeScript
- Vite
- Tailwind CSS v4
- Shadcn Vue
- Chart.js (图表)
- Pinia (状态管理)

## 快速开始

### 开发环境

1. 克隆仓库

```bash
git clone https://github.com/yourusername/web-tracing-admin.git
cd web-tracing-admin
```

2. 启动后端

```bash
cd backend
go mod tidy
go run cmd/main.go
```

或者使用 Air 热重载：

```bash
air
```

3. 启动前端

```bash
cd frontend
pnpm install
pnpm dev
```

4. 访问应用

打开浏览器访问 http://localhost:5173

### 使用 Docker 部署

1. 构建并启动容器

```bash
docker-compose up -d
```

2. 访问应用

打开浏览器访问 http://localhost

## 项目结构

```
web-tracing-admin/
├── backend/               # 后端代码
│   ├── cmd/               # 入口文件
│   ├── config/            # 配置文件
│   ├── internal/          # 内部包
│   │   ├── api/           # API 处理器
│   │   ├── middleware/    # 中间件
│   │   ├── model/         # 数据模型
│   │   ├── repository/    # 数据访问层
│   │   └── service/       # 业务逻辑层
│   ├── migrations/        # 数据库迁移
│   └── pkg/               # 公共包
├── frontend/              # 前端代码
│   ├── public/            # 静态资源
│   └── src/               # 源代码
│       ├── assets/        # 资源文件
│       ├── components/    # 组件
│       ├── composables/   # 组合式函数
│       ├── layouts/       # 布局组件
│       ├── pages/         # 页面组件
│       ├── services/      # 服务
│       ├── stores/        # 状态管理
│       └── utils/         # 工具函数
├── docker-compose.yml     # Docker Compose 配置
├── Dockerfile.backend     # 后端 Dockerfile
├── Dockerfile.frontend    # 前端 Dockerfile
└── nginx.conf             # Nginx 配置
```

## SDK 集成

要在您的前端项目中集成 web-tracing-sdk，请参考以下步骤：

### Vue 3 项目

```javascript
import { createApp } from 'vue'
import App from './App.vue'
import WebTracing from '@web-tracing/vue3'

const app = createApp(App)

app.use(WebTracing, {
  dsn: 'http://your-server-url/trackweb',
  appName: 'your-app-name',
  debug: true,
  pv: true,
  performance: true,
  error: true,
  event: true
})

app.mount('#app')
```

## 许可证

MIT
