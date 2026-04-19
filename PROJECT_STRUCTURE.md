# 项目结构详解

本文档说明当前仓库的目录结构、运行方式，以及内容数据之间的绑定关系。

## 📁 完整目录树

```
personal-portfolio/
├── client/                          # 主站前端（React + Vite）
│   ├── public/                      # 静态资源（构建时直接复制到 dist/public）
│   │   ├── blog.json                # 博客列表数据（运行时 fetch）
│   │   ├── photos.json              # 照片墙数据（运行时 fetch）
│   │   ├── projects.json            # 项目列表数据（运行时 fetch）
│   │   ├── blog/                    # 博客 Markdown
│   │   ├── images/photos/           # 照片墙图片
│   │   ├── images/projects/         # 项目正文图片
│   │   └── projects/                # 项目 Markdown
│   │
│   ├── src/                         # 源代码目录
│   │   ├── components/              # 可复用组件
│   │   │   ├── Navigation.tsx       # 导航栏组件
│   │   │   ├── ProjectCard.tsx      # 项目卡片组件
│   │   │   ├── Footer.tsx           # 页脚组件
│   │   │   ├── ErrorBoundary.tsx    # 错误边界（模板提供）
│   │   │   └── ui/                  # shadcn/ui 组件库
│   │   │
│   │   ├── pages/                   # 页面级组件
│   │   │   ├── Home.tsx             # 首页
│   │   │   ├── Projects.tsx         # 项目集页面
│   │   │   ├── ProjectDetail.tsx    # 项目详情页
│   │   │   ├── Resume.tsx           # 简历页面
│   │   │   ├── Blog.tsx             # 博客/碎碎念页面
│   │   │   ├── BlogDetail.tsx       # 博客详情页
│   │   │   ├── Photos.tsx           # 照片墙页面
│   │   │   └── NotFound.tsx         # 404 页面
│   │   │
│   │   ├── lib/                     # 工具函数和库
│   │   │   ├── markdown.ts          # Markdown 解析器
│   │   │   └── utils.ts             # 通用工具函数
│   │   │
│   │   ├── contexts/                # React Context（状态管理）
│   │   │   └── ThemeContext.tsx     # 主题上下文
│   │   │
│   │   ├── hooks/                   # 自定义 React Hooks
│   │   │   └── useFetchData.ts      # 通用 fetch 数据 Hook
│   │   │
│   │   ├── App.tsx                  # 应用主组件（路由定义）
│   │   ├── main.tsx                 # React 入口文件
│   │   └── index.css                # 全局样式和设计系统
│   │
│   ├── index.html                   # HTML 模板
│   └── vite.config.ts               # Vite 配置文件
│
├── server/                          # 服务器目录（占位符）
│   └── index.ts                     # 服务器入口（静态项目不需要）
│
├── shared/                          # 共享代码目录（占位符）
│   └── const.ts                     # 共享常量
│
├── DEPLOYMENT_GUIDE.md              # 部署手册
├── PROJECT_STRUCTURE.md             # 本文件
├── package.json                     # 项目配置和依赖
├── pnpm-lock.yaml                   # 锁文件
├── docker-compose.yml               # 线上部署的 Compose 编排
├── infra/                           # Dockerfile 与 Nginx 配置
├── apps/boolcore/                   # BoolCore 前后端源码
├── tsconfig.json                    # TypeScript 配置
├── vite.config.ts                   # Vite 配置
└── .github/workflows/deploy.yml     # GitHub Actions 自动部署
```

---

## 📄 关键文件说明

### 前端源代码

#### `client/src/App.tsx`
**用途**：应用主组件，定义路由和全局布局

**主要内容**：
- 路由定义（使用 Wouter）
- 主题提供者（ThemeProvider）
- 错误边界（ErrorBoundary）
- 全局 Toast 组件（Toaster）

**修改场景**：添加新路由、修改全局主题

#### `client/src/index.css`
**用途**：全局样式和设计系统定义

**主要内容**：
- Tailwind CSS 导入
- CSS 变量定义（颜色、间距、圆角）
- 排版系统（字体、字号、行高）
- 全局样式规则

**修改场景**：修改颜色、字体、间距等设计元素

#### `client/src/pages/Home.tsx`
**用途**：首页入口页

**主要内容**：
- 顶部个人信息
- 站内入口卡片
- 导航和页脚

**修改场景**：修改首页入口、欢迎信息、个人展示

#### `client/src/pages/Projects.tsx`
**用途**：项目集页面

**主要内容**：
- 读取 `projects.json`
- ProjectCard 组件的使用

**修改场景**：添加/修改项目信息

#### `client/src/pages/Resume.tsx`
**用途**：简历页面

**主要内容**：
- 工作经历
- 教育背景
- 技能列表

**修改场景**：更新简历信息

#### `client/src/pages/Blog.tsx`
**用途**：博客/碎碎念页面

**主要内容**：
- 读取 `blog.json`
- 渲染博客文章列表

**修改场景**：发布新文章、修改现有文章

#### `client/src/pages/Photos.tsx`
**用途**：照片墙页面

**主要内容**：
- 读取 `photos.json`
- 渲染照片网格
- 点击后使用 Dialog 查看大图与文案

**修改场景**：新增照片、调整排序、修改展示文案

### 组件

#### `client/src/components/Navigation.tsx`
**用途**：导航栏组件

**特点**：
- 响应式设计（移动端菜单）
- 粘性定位
- 菜单项配置化

**修改场景**：修改菜单项、修改导航样式

#### `client/src/components/ProjectCard.tsx`
**用途**：项目卡片组件

**特点**：
- 显示项目标题、描述、标签
- 支持外部链接
- 极简设计风格

**修改场景**：修改卡片样式、添加新字段

#### `client/src/components/Footer.tsx`
**用途**：页脚组件

**特点**：
- 社交媒体链接
- 版权信息
- 自动年份更新

**修改场景**：修改社交链接、修改页脚信息

### 工具和库

#### `client/src/lib/markdown.ts`
**用途**：Markdown 解析工具

**功能**：
- 将 Markdown 转换为 HTML
- 提取 Frontmatter 元数据
- 生成摘要
- 计算阅读时间

**使用场景**：处理 Markdown 格式的内容

### 内容数据

#### `client/public/projects.json`
**用途**：项目列表数据（页面运行时通过 `fetch("/projects.json")` 加载）

**格式**：
```json
{
  "id": "项目唯一标识",
  "title": "项目标题",
  "description": "项目描述",
  "tags": ["标签1", "标签2"],
  "link": "项目链接",
  "detailsFile": "project-1.md",
  "date": "年份"
}
```

**修改场景**：添加/修改项目信息

#### `client/public/photos.json`
**用途**：照片墙数据（页面运行时通过 `fetch("/photos.json")` 加载）

**格式**：
```json
{
  "id": "photo-1",
  "title": "照片标题",
  "date": "2026-04-19",
  "description": "一句简洁说明",
  "image": "/images/photos/example.jpg",
  "alt": "图片替代文本"
}
```

**修改场景**：新增照片、改排序、改文案

#### `client/public/blog.json`
**用途**：博客列表数据（页面运行时通过 `fetch("/blog.json")` 加载）

#### `client/public/blog/*.md` / `client/public/projects/*.md`
**用途**：文章/项目详情正文（详情页运行时 fetch 对应 Markdown 并渲染）

#### 内容绑定关系

**项目详情绑定**
- 列表页先读取 `client/public/projects.json`
- 用户点击 `/projects/:id` 后，详情页会根据 `id` 找到对应项目
- 再读取该项里的 `detailsFile`
- 最终请求 `/projects/${detailsFile}` 渲染 Markdown

例子：
```json
{
  "id": "project-1",
  "detailsFile": "project-1.md"
}
```

表示访问 `/projects/project-1` 时，会加载 `client/public/projects/project-1.md`

**博客详情绑定**
- 逻辑与项目一致：`blog.json` 中的 `id` 决定路由，`detailsFile` 决定正文 Markdown

**照片墙绑定**
- `photos.json` 只负责顺序、标题、日期、说明和图片路径
- 图片文件统一放在 `client/public/images/photos/`
- 页面会严格按 `photos.json` 的数组顺序展示，不做自动分类和排序

### 配置文件

#### `package.json`
**用途**：项目配置和依赖管理

**主要内容**：
- 项目元数据（名称、版本）
- 脚本命令（dev、build、preview）
- 依赖列表
- 开发依赖列表

**常见修改**：添加新依赖

#### `vite.config.ts`
**用途**：Vite 构建工具配置

**主要内容**：
- React 插件配置
- 路径别名（@）
- 构建选项

**修改场景**：修改构建配置、添加新插件

#### `tsconfig.json`
**用途**：TypeScript 配置

**主要内容**：
- 编译选项
- 路径映射
- 库配置

**修改场景**：修改 TypeScript 编译选项

#### `client/index.html`
**用途**：HTML 模板

**主要内容**：
- 页面元数据（标题、描述）
- Google Fonts 导入
- 分析脚本

**修改场景**：修改页面标题、添加字体、修改分析脚本

### 文档文件

#### `DEPLOYMENT_GUIDE.md`
**用途**：详细的部署手册

**内容**：
- 系统要求
- 前置准备
- 构建项目
- Nginx 配置
- SSL 证书配置
- 自动化部署
- 故障排查

#### `PROJECT_STRUCTURE.md`
**用途**：本文件，项目结构详解

### 脚本和配置

#### `.github/workflows/deploy.yml`
**用途**：GitHub Actions 自动部署工作流

**功能**：
- 构建并推送三个镜像到 GHCR
- SSH 到服务器同步仓库到 `origin/main`
- 更新业务容器
- 重启 `edge` 容器，避免 upstream 解析串位

---

## 🔄 数据流

### 页面加载流程

```
浏览器请求
    ↓
Nginx 服务器
    ↓
index.html
    ↓
client/src/main.tsx (React 入口)
    ↓
client/src/App.tsx (路由和主题)
    ↓
对应的页面组件 (Home/Projects/ProjectDetail/Blog/BlogDetail/Photos/Resume)
    ↓
渲染到浏览器
```

### 组件树结构

```
App
├── ThemeProvider
│   └── TooltipProvider
│       ├── Router
│       │   ├── Home
│       │   ├── Projects
│       │   ├── ProjectDetail
│       │   ├── Blog
│       │   ├── BlogDetail
│       │   ├── Photos
│       │   ├── Resume
│       │   └── NotFound
│       └── Toaster
```

---

## 📝 文件修改指南

### 常见修改场景

| 需求 | 修改文件 | 说明 |
| :--- | :--- | :--- |
| 修改个人名字 | `client/src/config/site.ts` | 修改站点基础信息 |
| 修改导航菜单 | `Navigation.tsx` | 修改 `navItems` 数组 |
| 添加项目 | `client/public/projects.json` + `client/public/projects/*.md` | 列表和正文要同时补 |
| 添加照片 | `client/public/images/photos/` + `client/public/photos.json` | 图片和元数据要同时补 |
| 修改照片墙顺序 | `client/public/photos.json` | 页面按 JSON 数组顺序展示 |
| 修改颜色 | `index.css` | 修改 CSS 变量 |
| 修改字体 | `client/src/index.css` | 调整全局字体变量和排版 |
| 添加新页面 | 创建新 `.tsx` 文件 | 在 `pages/` 目录中创建 |
| 添加新组件 | 创建新 `.tsx` 文件 | 在 `components/` 目录中创建 |
| 修改 Nginx 配置 | `infra/nginx/*.conf` | 由 GitHub Actions 同步到服务器 |

---

## 🚀 开发工作流

### 添加新功能的步骤

1. **创建组件或页面**
   ```bash
   # 在 client/src/components/ 或 client/src/pages/ 中创建新文件
   ```

2. **编写代码**
   ```tsx
   // 使用 TypeScript 和 React Hooks
   // 遵循现有的代码风格
   ```

3. **测试**
   ```bash
   pnpm dev
   # 在浏览器中测试
   ```

4. **构建**
   ```bash
   pnpm build
   ```

5. **部署**
   ```bash
   git push origin main
   ```

---

## 💡 最佳实践

### 代码组织

- ✅ 将可复用组件放在 `components/` 中
- ✅ 将页面级组件放在 `pages/` 中
- ✅ 将工具函数放在 `lib/` 中
- ✅ 将运行时内容数据放在 `client/public/` 中

### 命名规范

- ✅ 组件文件使用 PascalCase（如 `Navigation.tsx`）
- ✅ 工具文件使用 camelCase（如 `markdown.ts`）
- ✅ 使用描述性的文件名

### 样式管理

- ✅ 优先使用 Tailwind CSS 类名
- ✅ 在 `index.css` 中定义全局样式
- ✅ 使用 CSS 变量保持一致性

---

## 📚 相关文档

- **部署指南**：`DEPLOYMENT_GUIDE.md`
- **项目结构**：`PROJECT_STRUCTURE.md`
