# FloraWiki（植物养护知识百科平台）

为园艺爱好者提供全面的植物养护指南：品种库、养护文章、病虫害诊断、季节养护日历、我的花园、问答社区与养护小测验，支持图文展示与个人花园/提醒管理。

## Docker Compose 一键启动（推荐）

```bash
cp .env.example .env
docker compose up -d --build
```

启动后访问：

- 前端：http://localhost:8102
- 后端 API：http://localhost:3102
- 健康检查：http://localhost:3102/healthz

默认种子账号：`admin / admin123`（管理员）、`gardener / user123`（普通用户）。

关闭并清理数据：

```bash
docker compose down -v --remove-orphans
```

## 本地开发（备选）

后端（Go 1.22 + Gin + GORM）：

```bash
cd backend
go mod tidy
go run ./cmd/server
go build ./...
go test ./...
```

前端（Vue 3 + TypeScript + Element Plus + Vite）：

```bash
cd frontend
npm install
npm run dev     # 开发服务器，/api 代理到 http://localhost:3102
npm run build   # 生产构建
```

## 技术栈

| 分层 | 技术 |
| --- | --- |
| 前端 | Vue 3 + TypeScript + Element Plus + Vite + Pinia + Vue Router |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 认证 | JWT（github.com/golang-jwt/jwt/v5）+ RBAC |
| 其他 | validator/v10、log/slog、Nginx |

## 项目目录结构

```
gb-61/
├── docker-compose.yml
├── .env.example
├── README.md
├── database/
│   └── init.sql                 # 建表 + 种子数据（首次启动自动执行）
├── backend/
│   ├── cmd/server/              # main.go + migrate/seed
│   └── internal/
│       ├── config/              # 环境变量解析
│       ├── model/               # 9 个实体，按实体分文件
│       ├── repository/          # 按实体分文件，哨兵错误
│       ├── service/             # 按实体分文件，构造器注入
│       ├── handler/             # 按实体分文件 + upload/home
│       ├── router/              # router.go + 按实体分文件
│       ├── middleware/          # auth/rbac/rate_limiter/error_handler/logger/cors
│       ├── dto/                 # 请求/响应结构体 + 统一响应包装
│       ├── constants/           # plant/article/favorite/error_codes/log_templates/messages
│       └── util/                # jwt/logger/formatters/app_error/file/season
└── frontend/
    ├── nginx.conf               # /api 反代 backend + SPA
    └── src/
        ├── api/                 # user/plant/article/pest/reminder/favorite/garden/question
        ├── stores/              # authStore/userStore/plantStore/articleStore/reminderStore
        ├── components/common/   # PlantCard/CareArticleCard/FavoriteButton/SearchFilter/...
        ├── hooks/               # useAuth/useFavorite/useReminderStats/useQuiz
        ├── pages/               # Home/PlantLibrary/PlantDetail/ArticleList/.../Login
        ├── router/              # index.ts + guards.ts
        ├── utils/               # request/dateFormat/season
        └── constants/           # plant/article/favorite/errorCodes
```

## 环境变量

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| COMPOSE_PROJECT_NAME | gbplantwiki | Compose 项目名/容器前缀 |
| DB_NAME | gbplantwiki_db | MySQL 库名 |
| DB_USER | gbplantwiki_user | MySQL 用户 |
| DB_PASSWORD | gbplantwiki_pwd | MySQL 密码 |
| DB_ROOT_PASSWORD | gbplantwiki_root | MySQL root 密码 |
| JWT_SECRET | change_me_to_a_long_random_string | JWT 签名密钥（生产必改） |
| FRONTEND_PORT | 8102 | 前端端口 |
| BACKEND_PORT | 3102 | 后端端口 |
| DB_PORT | 3502 | 数据库端口 |

## Docker 部署说明

- 端口映射：前端 `8102:80`，后端 `${BACKEND_PORT:-3102}:8080`，数据库 `${DB_PORT:-3502}:3306`
- 数据卷：`db_data`（MySQL 数据）、`uploads`（上传图片）
- 依赖顺序：db healthcheck → backend `depends_on: db: service_healthy` → frontend `depends_on: backend`
- 常见问题：
  - 端口冲突：修改 `.env` 中对应端口后 `docker compose up -d`
  - 数据重置：`docker compose down -v` 后重新 `up`
  - 中文目录名：Compose 使用命名卷与容器名，不依赖目录路径，任意目录名可启动

## API 接口清单

> 后端统一前缀 `/api/v1`，响应统一为 `{ "code": 0, "message": "ok", "data": ... }`。标注「登录」的接口需携带 `Authorization: Bearer <JWT>`。

| 方法 | 路径 | 权限 | 说明 |
| --- | --- | --- | --- |
| GET | /healthz | 公开 | 健康检查 |
| GET | /api/v1/home/overview | 公开 | 首页聚合：热门品种+最新文章+当季任务 |
| POST | /api/v1/users/register | 公开（限流） | 注册并返回 JWT |
| POST | /api/v1/users/login | 公开（限流） | 登录并返回 JWT |
| GET | /api/v1/users/me | 登录 | 获取当前用户资料 |
| PUT | /api/v1/users/me | 登录 | 更新当前用户资料 |
| GET | /api/v1/plants | 公开 | 品种分页列表/筛选 |
| GET | /api/v1/plants/:id | 公开 | 品种详情 |
| POST | /api/v1/plants | 管理员（限流） | 新增品种 |
| PUT | /api/v1/plants/:id | 管理员 | 更新品种 |
| DELETE | /api/v1/plants/:id | 管理员 | 删除品种 |
| GET | /api/v1/articles | 公开 | 养护文章分页列表/筛选 |
| GET | /api/v1/articles/:id | 公开 | 文章详情并自增阅读数 |
| POST | /api/v1/articles | 登录（限流） | 发布文章 |
| PUT | /api/v1/articles/:id | 登录 | 编辑自己的文章 |
| DELETE | /api/v1/articles/:id | 登录 | 删除自己的文章 |
| GET | /api/v1/pests | 公开 | 病虫害手册搜索 |
| GET | /api/v1/pests/:id | 公开 | 病虫害详情 |
| POST | /api/v1/pests | 管理员（限流） | 新增病虫害条目 |
| PUT | /api/v1/pests/:id | 管理员 | 更新病虫害条目 |
| DELETE | /api/v1/pests/:id | 管理员 | 删除病虫害条目 |
| GET | /api/v1/reminders | 登录 | 当前用户提醒列表（自动标记逾期） |
| GET | /api/v1/reminders/calendar | 登录 | 按月查询提醒 |
| POST | /api/v1/reminders | 登录（限流） | 创建养护提醒 |
| PUT | /api/v1/reminders/:id/status | 登录 | 状态流转 pending/done |
| DELETE | /api/v1/reminders/:id | 登录 | 删除提醒 |
| GET | /api/v1/favorites | 登录 | 收藏列表 |
| POST | /api/v1/favorites | 登录（限流） | 添加收藏 |
| DELETE | /api/v1/favorites/:targetType/:targetId | 登录 | 取消收藏 |
| GET | /api/v1/gardens | 登录 | 我的花园列表 |
| POST | /api/v1/gardens | 登录（限流） | 加入我的花园 |
| PUT | /api/v1/gardens/:id/reminder | 登录 | 关联养护提醒 |
| DELETE | /api/v1/gardens/:id | 登录 | 移除花园条目 |
| GET | /api/v1/questions | 公开 | 问答列表 |
| GET | /api/v1/questions/:id | 公开 | 问题详情 |
| GET | /api/v1/questions/:id/answers | 公开 | 问题回答列表 |
| POST | /api/v1/questions | 登录（限流） | 发布问题 |
| POST | /api/v1/questions/:id/answers | 登录 | 回答问题 |
| PUT | /api/v1/questions/:id/adopt | 登录 | 采纳最佳回答（事务：清旧最佳+标最佳+关闭问题） |
| PUT | /api/v1/answers/:id/like | 登录 | 回答点赞 |
| POST | /api/v1/uploads | 登录（限流） | 上传图片 |

## 枚举出现位置清单

### PlantType（植物类型：flower/foliage/succulent/aquatic）

- 后端：`backend/internal/constants/plant.go`（定义）、`backend/internal/model/plant_species.go`（GORM 模型字段）、`backend/internal/service/plant_species_service.go`（校验/日志）、`backend/internal/util/formatters.go`（PlantTypeText）、`backend/internal/constants/log_templates.go`（日志模板）、`backend/internal/constants/error_codes.go`（错误码）、`database/init.sql`（种子数据）
- 前端：`frontend/src/constants/plant.ts`（定义）、`frontend/src/components/common/PlantCard.vue`（类型标签）、`frontend/src/pages/PlantLibrary.vue`（筛选器）、`frontend/src/pages/PlantDetail.vue`（详情）、`frontend/src/utils/season.ts`（plantTypeOptions）

### CareTopicTag（养护主题：fertilizing/pruning/repotting/pest_control/propagation）

- 后端：`backend/internal/constants/article.go`（定义）、`backend/internal/model/care_article.go`（模型）、`backend/internal/service/care_article_service.go`（校验）、`backend/internal/util/formatters.go`（CareTopicText）、`backend/internal/constants/log_templates.go`、`database/init.sql`
- 前端：`frontend/src/constants/article.ts`（定义）、`frontend/src/components/common/CareArticleCard.vue`（标签）、`frontend/src/pages/ArticleList.vue`（筛选）、`frontend/src/pages/ArticleDetail.vue`（详情）

### FavoriteTargetType（收藏目标：plant/article）

- 后端：`backend/internal/constants/favorite.go`（定义）、`backend/internal/model/favorite.go`（模型）、`backend/internal/service/favorite_service.go`（校验）、`backend/internal/constants/log_templates.go`、`database/init.sql`
- 前端：`frontend/src/constants/favorite.ts`（定义）、`frontend/src/components/common/FavoriteButton.vue`（交互）、`frontend/src/pages/Garden.vue` 与 `frontend/src/pages/Profile.vue`（收藏夹列表）

## 横切关注点

- 认证授权（JWT + RBAC）：`internal/middleware/auth.go`、`rbac.go`、`util/jwt.go`、`router/*.go`、前端 `router/index.ts` 守卫、`components/common/RoleGuard.vue`
- 全局错误处理：`internal/middleware/error_handler.go`、`util/app_error.go`、`constants/error_codes.go`、前端 `utils/request.ts` 拦截器
- 接口限流：`internal/middleware/rate_limiter.go`、登录/发布/上传接口启用
- 文件上传：`handler/upload_handler.go`、`service` 本地存储、`util/file.go`、Nginx `/uploads/` 代理

## License

MIT
