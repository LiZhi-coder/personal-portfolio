# 项目结构详解

本文档详细说明项目的目录结构和各文件的用途。

## 📁 完整目录树

```
personal-portfolio/
├── client/                          # 前端应用目录
│   ├── public/                      # 静态资源（直接复制到 dist）
│   │   ├── blog.json                # 博客列表数据（运行时 fetch）
│   │   ├── projects.json            # 项目列表数据（运行时 fetch）
│   │   ├── blog/                    # 博客 Markdown
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
│   │   │   ├── Home.tsx             # 首页（主页面）
│   │   │   ├── Projects.tsx         # 项目集页面
│   │   │   ├── Resume.tsx           # 简历页面
│   │   │   ├── Blog.tsx             # 博客/碎碎念页面
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
│   │   │   └── useTheme.ts          # 主题 Hook
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
├── docs/                            # 文档目录（可选）
│   └── ...
│
├── DEPLOYMENT_GUIDE.md              # 部署手册
├── QUICK_START.md                   # 快速开始指南
├── PROJECT_STRUCTURE.md             # 本文件
├── README.md                        # 项目说明
├── nginx.conf.template              # Nginx 配置模板
├── deploy.sh                        # 自动化部署脚本
├── package.json                     # 项目配置和依赖
├── tsconfig.json                    # TypeScript 配置
├── vite.config.ts                   # Vite 配置
└── .gitignore                       # Git 忽略文件
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
**用途**：首页，包含所有主要内容

**主要内容**：
- Hero 介绍部分
- 项目集部分（导入 Projects 组件）
- 简历部分（导入 Resume 组件）
- 博客部分（导入 Blog 组件）
- 联系表单

**修改场景**：修改个人信息、添加新的页面部分

#### `client/src/pages/Projects.tsx`
**用途**：项目集页面

**主要内容**：
- 项目列表数据
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
- 博客文章列表
- 文章元数据（日期、阅读时间）

**修改场景**：发布新文章、修改现有文章

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
  "date": "年份"
}
```

**修改场景**：添加/修改项目信息

#### `client/public/blog.json`
**用途**：博客列表数据（页面运行时通过 `fetch("/blog.json")` 加载）

#### `client/public/blog/*.md` / `client/public/projects/*.md`
**用途**：文章/项目详情正文（详情页运行时 fetch 对应 Markdown 并渲染）

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

#### `README.md`
**用途**：项目说明文档

**内容**：
- 项目特性
- 快速开始
- 项目结构
- 设计系统
- 自定义指南
- 部署说明

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

#### `QUICK_START.md`
**用途**：快速开始指南

**内容**：
- 本地开发步骤
- 自定义内容指南
- 常见问题解答

#### `PROJECT_STRUCTURE.md`
**用途**：本文件，项目结构详解

### 脚本和配置

#### `deploy.sh`
**用途**：自动化部署脚本

**功能**：
- 检查前置条件
- 拉取最新代码
- 安装依赖
- 构建项目
- 设置权限
- 重新加载 Nginx

**使用**：`./deploy.sh`

#### `nginx.conf.template`
**用途**：Nginx 配置模板

**内容**：
- 虚拟主机配置
- 静态资源缓存
- 路由处理
- 安全头配置
- SSL 配置示例

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
对应的页面组件 (Home/Projects/Resume/Blog)
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
│       │   │   ├── Navigation
│       │   │   ├── Projects
│       │   │   ├── Resume
│       │   │   ├── Blog
│       │   │   └── Footer
│       │   └── NotFound
│       └── Toaster
```

---

## 📝 文件修改指南

### 常见修改场景

| 需求 | 修改文件 | 说明 |
| :--- | :--- | :--- |
| 修改个人名字 | `Home.tsx` | 修改 H1 标题 |
| 修改导航菜单 | `Navigation.tsx` | 修改 `navItems` 数组 |
| 添加项目 | `projects.json` | 添加新的项目对象 |
| 修改颜色 | `index.css` | 修改 CSS 变量 |
| 修改字体 | `index.html` + `index.css` | 导入新字体并更新 CSS |
| 添加新页面 | 创建新 `.tsx` 文件 | 在 `pages/` 目录中创建 |
| 添加新组件 | 创建新 `.tsx` 文件 | 在 `components/` 目录中创建 |
| 修改 Nginx 配置 | `nginx.conf.template` | 复制到服务器后修改 |

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
   # 上传到服务器并运行 deploy.sh
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

- **快速开始**：`QUICK_START.md`
- **部署指南**：`DEPLOYMENT_GUIDE.md`
- **项目说明**：`README.md`

祝您开发愉快！🎉
