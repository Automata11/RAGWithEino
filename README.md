# RAG with Eino

一个基于Go后端和Vue 3前端的RAG (检索增强生成) 应用程序，集成了Eino向量数据库。

## 技术架构

### 后端技术栈
- **Go** - 主要编程语言
- **Gin** - Web框架
- **GORM** - ORM框架
- **MySQL** - 关系型数据库
- **Redis** - 缓存和JWT令牌存储
- **Milvus** - 向量数据库
- **Eino** - RAG查询引擎
- **JWT** - 身份验证

### 前端技术栈
- **Vue 3** - JavaScript框架
- **TypeScript** - 类型系统
- **Vite** - 构建工具
- **Element Plus** - UI组件库
- **Pinia** - 状态管理
- **Vue Router** - 路由管理
- **Axios** - HTTP客户端

## 功能特性

### 用户功能
- 🔐 用户注册和登录
- 📁 文档上传和管理
- 💬 基于文档的智能问答
- 📊 文档处理进度查询
- 👤 个人资料管理

### 管理员功能
- 👥 用户管理
- 📈 系统监控面板
- 🔧 系统配置管理

## 项目结构

```
RAGWithEino/
├── backend/                 # Go后端
│   ├── cmd/server/         # 应用程序入口
│   ├── config/            # 配置管理
│   ├── internal/          # 内部包
│   │   ├── auth/         # 认证服务
│   │   ├── database/     # 数据库连接
│   │   ├── handlers/     # HTTP处理器
│   │   ├── middleware/   # 中间件
│   │   ├── models/       # 数据模型
│   │   └── services/     # 业务逻辑
│   ├── pkg/              # 公共包
│   ├── uploads/          # 文件上传目录
│   └── .env.example      # 环境变量示例
└── frontend/              # Vue 3前端
    ├── src/
    │   ├── components/   # Vue组件
    │   ├── views/        # 页面视图
    │   ├── stores/       # Pinia状态管理
    │   ├── router/       # 路由配置
    │   ├── types/        # TypeScript类型
    │   └── utils/        # 工具函数
    ├── public/           # 静态资源
    └── index.html        # 入口HTML
```

## 快速开始

### 环境要求

- Go 1.19+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+
- Milvus 2.0+

### 后端设置

1. **进入后端目录**
```bash
cd backend
```

2. **配置环境变量**
```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库连接等信息
```

3. **安装依赖**
```bash
go mod download
```

4. **启动数据库服务**
```bash
# 启动MySQL
systemctl start mysql

# 启动Redis
systemctl start redis

# 启动Milvus (使用Docker)
docker run -p 19530:19530 -p 9091:9091 milvusdb/milvus:latest standalone
```

5. **运行应用**
```bash
go run cmd/server/main.go
```

后端服务将在 http://localhost:8080 启动

### 前端设置

1. **进入前端目录**
```bash
cd frontend
```

2. **安装依赖**
```bash
npm install
```

3. **启动开发服务器**
```bash
npm run dev
```

前端应用将在 http://localhost:5173 启动

## 环境变量配置

在 `backend/.env` 文件中配置以下环境变量：

```bash
# 服务器配置
PORT=8080

# 数据库配置
MYSQL_DSN=root:password@tcp(localhost:3306)/ragdb?charset=utf8mb4&parseTime=True&loc=Local
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
MILVUS_ADDR=localhost:19530

# JWT配置
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY_HOURS=24

# Eino配置
MODEL_NAME=gpt-3.5-turbo
EMBEDDING_DIM=1536
MAX_TOKENS=4096

# OpenAI API (用于Eino)
OPENAI_API_KEY=your-openai-api-key
```

## API文档

### 认证接口
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/logout` - 用户登出
- `GET /api/v1/auth/profile` - 获取用户信息

### 文档管理
- `GET /api/v1/documents` - 获取文档列表
- `POST /api/v1/documents` - 上传文档
- `GET /api/v1/documents/:id` - 获取文档详情
- `DELETE /api/v1/documents/:id` - 删除文档
- `GET /api/v1/documents/:id/status` - 获取处理状态

### RAG查询
- `POST /api/v1/rag/query` - 发送查询请求
- `GET /api/v1/rag/conversations` - 获取对话列表
- `GET /api/v1/rag/conversations/:id` - 获取对话详情

### 用户管理 (管理员)
- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `PUT /api/v1/users/:id` - 更新用户信息
- `DELETE /api/v1/users/:id` - 删除用户

## 部署

### Docker部署

1. **构建镜像**
```bash
# 构建后端镜像
cd backend
docker build -t rag-backend .

# 构建前端镜像
cd ../frontend
docker build -t rag-frontend .
```

2. **使用Docker Compose**
```bash
docker-compose up -d
```

### 生产环境部署

1. **构建前端**
```bash
cd frontend
npm run build
```

2. **构建后端**
```bash
cd backend
go build -o main cmd/server/main.go
```

3. **配置Nginx**
```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /path/to/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # API代理
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 开发指南

### 添加新功能

1. **后端开发**
   - 在 `internal/models/` 中定义数据模型
   - 在 `internal/handlers/` 中实现HTTP处理器
   - 在 `cmd/server/main.go` 中注册路由

2. **前端开发**
   - 在 `src/types/` 中定义TypeScript类型
   - 在 `src/stores/` 中添加状态管理
   - 在 `src/views/` 中创建页面组件

### 代码规范

- 后端使用 `go fmt` 格式化代码
- 前端使用 ESLint 和 Prettier
- 提交前运行测试

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 联系方式

如有问题或建议，请创建 Issue 或联系项目维护者。

---

## English Version

# RAG with Eino

A RAG (Retrieval Augmented Generation) application with Go backend and Vue 3 frontend, integrated with Eino vector database.

## Quick Start

### Prerequisites
- Go 1.19+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+
- Milvus 2.0+

### Backend Setup
```bash
cd backend
cp .env.example .env
go mod download
go run cmd/server/main.go
```

### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

The application will be available at http://localhost:5173

## Features

- 🔐 User authentication and authorization
- 📁 Document upload and management
- 💬 AI-powered document Q&A
- 📊 Processing progress tracking
- 👤 User profile management
- 👥 Admin user management

## License

MIT License - see the [LICENSE](LICENSE) file for details.