# 个人网站 + BoolCore 工具（前后端分离）企业级 Docker + CI/CD 部署手册

适用场景：
- 你已有 DO Ubuntu 24 服务器，并按旧版 `DEPLOYMENT_GUIDE.md` 通过 Nginx 静态托管部署过
- 现在要升级为：Docker + GitHub Actions 自动部署 + 单域名多应用
- BoolCore 为独立工具服务，挂载在主站 `/tools/boolcore/`

本手册目标：
- 一步一步走完整“企业运维”流程（含留痕）
- 给出可复现的部署命令、配置文件、回滚方法

---

## 1. 最终架构（蓝图）

### 访问路径设计
- 主站：`https://www.hui-nexus.me/`
- 工具前端：`https://www.hui-nexus.me/tools/boolcore/`
- 工具 API：`https://www.hui-nexus.me/api/boolcore/`

### 网络拓扑
```
浏览器
  |
 HTTPS
  v
Edge Nginx (Docker)
  |-- /                     -> portfolio-web (静态)
  |-- /tools/boolcore/       -> boolcore-web (静态)
  '-- /api/boolcore/         -> boolcore-api (Go 服务)
```

### CI/CD 流程图
```
本地提交 -> GitHub Actions
                 |-> 构建镜像(3个)
                 |-> 推送 GHCR
                 |-> 服务器拉取镜像
                 |-> docker compose up -d
                 '--> 留痕：Action 记录 + 镜像 Tag + 服务器日志
```

---

## 2. 迁移前留痕（重要）

你的服务器已经按旧方案部署过，建议先保留旧配置以便回退：

```bash
# 1) 备份旧 Nginx 配置
sudo cp /etc/nginx/sites-available/personal-portfolio /root/nginx-personal-portfolio.bak

# 2) 记录当前线上资源
ls -la /var/www/personal-portfolio/dist/public

# 3) 记录旧版服务状态
systemctl status nginx
```

留痕点：
- 旧 Nginx 配置备份
- 旧网站构建产物目录清单
- 旧服务运行状态

---

## 3. 本地代码准备（已完成的改动点）

这些改动已经在仓库中完成：
- BoolCore 前端 API 改为相对路径 `/api/boolcore`
- BoolCore 前端 `base` 设为 `/tools/boolcore/`
- Go 后端监听 `0.0.0.0:8080` 并支持 `PORT`/`ADDR`/`CORS_ALLOW_ORIGINS`
- 新增 Dockerfile、Nginx 配置、`docker-compose.yml`
- 新增 GitHub Actions 工作流

---

## 4. GitHub Actions 准备（CI/CD 留痕入口）

### 4.1 开启 GHCR（镜像仓库）
1. GitHub 项目页 -> Settings -> Packages
2. 确保有权限使用 GHCR

### 4.2 添加 Secrets（必需）
在 GitHub 仓库 `Settings -> Secrets and variables -> Actions` 中添加：
- `SSH_HOST`：服务器 IP
- `SSH_USER`：服务器用户名（建议 `root` 或有 Docker 权限的用户）
- `SSH_KEY`：私钥（与服务器 `/root/.ssh/authorized_keys` 对应）

留痕点：
- GitHub Actions 中每次部署会生成完整日志
- GHCR 中每次构建会生成镜像版本（tag = commit sha）

---

## 5. 服务器准备（Ubuntu 24）

### 5.1 安装 Docker & Compose
```bash
sudo apt update
sudo apt install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo $VERSION_CODENAME) stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

### 5.2 迁移目录（推荐）
```bash
sudo mkdir -p /opt/personal-portfolio
sudo chown -R $USER:$USER /opt/personal-portfolio
```

### 5.3 开放 80/443（安全组 + UFW）
确认 DO 防火墙已放行 80/443（若使用 UFW）：
```bash
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw status
```

### 5.4 登录 GHCR（镜像为私有时必需）
如果 GHCR 镜像是私有的，服务器需要登录：
```bash
docker login ghcr.io -u <你的GitHub用户名>
```

### 5.5 拉取代码
```bash
cd /opt/personal-portfolio
# 建议使用 Git 同步
git clone git@github.com:<你的用户名>/<仓库名>.git .
```

### 5.6 停掉旧 Nginx（释放 80/443）
```bash
sudo systemctl stop nginx
sudo systemctl disable nginx
```

留痕点：
- `docker --version` / `docker compose version`
- 服务器 `/opt/personal-portfolio` 目录结构

---

## 6. 证书签发（首次部署必做）

当前 Docker 版本的 Nginx 会读取证书文件，因此需要先签发证书：

```bash
cd /opt/personal-portfolio
mkdir -p infra/certbot/www infra/certbot/conf

# 使用 standalone 模式签发（会占用 80 端口）
docker run --rm -p 80:80 \
  -v $(pwd)/infra/certbot/conf:/etc/letsencrypt \
  -v $(pwd)/infra/certbot/www:/var/www/certbot \
  certbot/certbot certonly --standalone \
  -d hui-nexus.me -d www.hui-nexus.me \
  --email 你的邮箱 --agree-tos --no-eff-email
```

成功后应看到：
```
/etc/letsencrypt/live/hui-nexus.me/
```

留痕点：
- 证书签发日志
- `infra/certbot/conf` 下的证书文件

---

## 7. 启动 Docker Compose（首次上线）

### 7.1 配置 `.env`
```bash
cd /opt/personal-portfolio
cp .env.example .env
nano .env
```

推荐填写：
```
REGISTRY=ghcr.io/你的GitHub用户名
IMAGE_TAG=latest
CORS_ALLOW_ORIGINS=http://localhost:5173,http://localhost:4173,https://hui-nexus.me,https://www.hui-nexus.me
```

### 7.2 启动容器
```bash
docker compose pull
docker compose up -d --no-build
```

### 7.3 验证
```bash
docker compose ps
curl -I https://www.hui-nexus.me/
curl -I https://www.hui-nexus.me/tools/boolcore/
curl -I https://www.hui-nexus.me/api/boolcore/ping
```

留痕点：
- `docker compose ps` 输出
- curl 状态码

---

## 8. CI/CD 全流程（自动化运维）

### 触发条件
- `main` 分支 push 自动部署
- 手动 `workflow_dispatch` 可手动触发

### CI/CD 在做什么
1) 构建 3 个镜像
   - `portfolio-web`
   - `boolcore-web`
   - `boolcore-api`
2) 推送到 GHCR
3) 服务器自动拉取镜像并重启服务

### 关键留痕点
- GitHub Actions 运行日志
- GHCR 镜像 tag（tag = commit sha）
- 服务器 `docker compose ps` 状态

---

## 9. 运维留痕（企业级习惯）

建议每次发布记录：
- 发布人、时间、commit sha
- GitHub Actions 运行链接
- 变更说明
- 回滚策略

常用命令：
```bash
# 查看容器状态
docker compose ps

# 查看日志
docker compose logs -n 200 --tail=200

# 查看镜像版本
docker image ls | head -n 20

# 查看部署 tag（最新 commit）
git log --oneline -n 5
```

---

## 10. 回滚方案（秒级回退）

1) 找到可用镜像 tag（commit sha）
2) 替换 `.env` 中 `IMAGE_TAG`
3) 重新拉取并启动

```bash
nano .env
# IMAGE_TAG=旧 commit sha

docker compose pull
docker compose up -d --no-build
```

---

## 11. 日常更新流程（企业实践版）

1) 本地改代码并提交
```bash
git add .
git commit -m "feat: update boolcore"
git push origin main
```

2) GitHub Actions 自动完成部署

3) 线上验证
```bash
curl -I https://www.hui-nexus.me/
```

---

## 12. 常见问题排查

- 端口被占用：确认 `systemctl stop nginx` 已生效
- 证书失效：`docker compose run --rm certbot renew`
- BoolCore 前端 404：确认 `/tools/boolcore/` 是否设置为 base
- API 404：确认 Nginx `location /api/boolcore/` 代理规则

---

## 13. 对应文件一览（你可以逐一核对）

- `docker-compose.yml`
- `infra/nginx/edge.conf`
- `infra/nginx/portfolio.conf`
- `infra/nginx/boolcore.conf`
- `infra/docker/Dockerfile.portfolio`
- `infra/docker/Dockerfile.boolcore-frontend`
- `infra/docker/Dockerfile.boolcore-backend`
- `.github/workflows/deploy.yml`
- `.env.example`

---

如果你准备好了，我们可以开始跑完整流程。你也可以让我按照这一份文档，直接在你的服务器上逐步执行并验证。
