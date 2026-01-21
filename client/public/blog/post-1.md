---
title: ECS 部署云服务系统
date: 2026-01-05 10:00:00
excerpt: 使用阿里云ECS部署用户服务展示系统
readingTime: 
---

# Cloud Server 生产部署

这篇记录我在 ECS Ubuntu 22.04 上部署 Cloud Server 的真实流程，目标是稳定、可回滚、可复现，并且只把 80/443 暴露到公网。以下内容都是我线上环境的实际做法和注意点。

## 信息
- 系统：Ubuntu 22.04
- 域名：surennongmu.com / www.surennongmu.com
- 后端监听：127.0.0.1:8080
- 证书目录：/etc/letsencrypt/live/surennongmu.com
- 发布目录：/opt/cloud-server/...
- 进程托管：systemd + NGINX

## 配置内容
1. main.go 写死了 `config.LoadConfig("configs/config.yaml")`，因此配置路径固定为：
   /opt/cloud-server/current/configs/config.yaml  
2. 使用Nginx 代替 TLS， Go 配置里的 tls.enabled 必须保持 false。

## 1. 主要内容
- ECS部署 Go + Gin
- 域名访问：https://surennongmu.com、https://www.surennongmu.com
- 公网只暴露 80/443，后端只监听 127.0.0.1:8080
- systemd 开机自启、故障自动重启、日志可查
- HTTPS 使用 Let’s Encrypt，HTTP 自动跳转 HTTPS
- 静态资源由 NGINX 直出，路径 /static/*

## 2. 
### 2.1 架构
```
Internet
  |
  | 80/443
  v
NGINX 入口层，TLS 终止、反向代理、静态资源
  |
  | http://127.0.0.1:8080
  v
cloud-server Go/Gin，systemd 托管，本机监听
```

### 2.2 主要命令
- 编译产生二进制 `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o`
  出现的问题：编译后的二进制文件 cloud-server 缺少执行权限（-rw-rw-r--） 添加执行权 `chmod +x`

## 3. 服务器目录
### 3.1 发布目录
- 版本目录  
  /opt/cloud-server/releases/20251212_001
- 当前线上软链  
  /opt/cloud-server/current -> /opt/cloud-server/releases/20251212_001
- 二进制输出  
  /opt/cloud-server/current/bin/cloud-server

### 3.2 配置与密钥
- 业务配置  
  /opt/cloud-server/current/configs/config.yaml
- 生产密钥  
  /etc/cloud-server/cloud-server.env 

### 3.3 systemd
- /etc/systemd/system/cloud-server.service

### 3.4 NGINX
- /etc/nginx/sites-available/surennongmu.conf
- /etc/nginx/sites-enabled/surennongmu.conf -> /etc/nginx/sites-available/surennongmu.conf

### 3.5 证书
- /etc/letsencrypt/live/surennongmu.com/

## 4. 一次性初始化
### 4.1 创建 deploy
```
sudo adduser deploy
sudo usermod -aG sudo deploy
```

### 4.2 SSH 基线
- 禁用 root 登录
- 使用密钥登录 deploy

验证
```
ssh deploy@<ECS公网IP>
```

### 4.3 安装依赖
```
sudo apt update
sudo apt install -y nginx
```

如需编译 Go，服务器需安装 Go。

### 4.4 创建目录
```
sudo mkdir -p /opt/cloud-server/releases
sudo mkdir -p /etc/cloud-server
sudo mkdir -p /var/log/cloud-server
sudo chown -R deploy:deploy /opt/cloud-server
```

## 5. 生产配置
### 5.1 配置文件
编辑
```
vim /opt/cloud-server/current/configs/config.yaml
```

关键要求
```
server:
  port: "127.0.0.1:8080"

tls:
  enabled: false
```

### 5.2 密钥 env
创建
```
sudo vim /etc/cloud-server/cloud-server.env
```

示例
```
SESSION_SECRET=<随机长串>
SYNC_API_KEY=<随机长串>
DB_PASSWORD=<数据库密码>
SMS_PASSWORD=<短信平台密码>
```

权限
```
sudo chown root:root /etc/cloud-server/cloud-server.env
sudo chmod 600 /etc/cloud-server/cloud-server.env
```

## 6. 编译与发布
### 6.1 上传新版本到 releases
```
/opt/cloud-server/releases/<版本号>
```

### 6.2 切换 current 软链
```
sudo ln -sfn /opt/cloud-server/releases/<版本号> /opt/cloud-server/current
ls -ld /opt/cloud-server/current
```

### 6.3 编译二进制
```
cd /opt/cloud-server/current
mkdir -p bin
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/cloud-server ./cmd
chmod +x /opt/cloud-server/current/bin/cloud-server
ls -lh /opt/cloud-server/current/bin/cloud-server
```


## 7. systemd 服务
### 7.1 服务文件
编辑
```
sudo vim /etc/systemd/system/cloud-server.service
```

内容
```
[Unit]
Description=Cloud Server Application
After=network.target

[Service]
Type=simple
User=deploy
Group=deploy
WorkingDirectory=/opt/cloud-server/current

EnvironmentFile=-/etc/cloud-server/cloud-server.env

ExecStart=/opt/cloud-server/current/bin/cloud-server

Restart=on-failure
RestartSec=5s
LimitNOFILE=65536

StandardOutput=journal
StandardError=journal
SyslogIdentifier=cloud-server

[Install]
WantedBy=multi-user.target
```

### 7.2 生效与启动
```
sudo systemctl daemon-reload
sudo systemctl enable cloud-server
sudo systemctl restart cloud-server
sudo systemctl status cloud-server --no-pager -l
```

查看日志
```
journalctl -u cloud-server -n 200 --no-pager
```

### 7.3 本机健康检查
```
curl -i http://127.0.0.1:8080/health
```

## 8. NGINX 站点与 HTTPS
### 8.1 启用机制
```
ls -la /etc/nginx/sites-enabled/
```
期望看到 surennongmu.conf 软链。

### 8.2 站点配置
备份
```
sudo cp /etc/nginx/sites-available/surennongmu.conf \
  /etc/nginx/sites-available/surennongmu.conf.backup-$(date +%Y%m%d-%H%M%S)
```

编辑
```
sudo vim /etc/nginx/sites-available/surennongmu.conf
```

证书目录必须是
```
/etc/letsencrypt/live/surennongmu.com/
```

### 8.3 配置检查与热加载
```
sudo nginx -t
sudo systemctl reload nginx
```

## 9. 上线验收
### 9.1 入口与反代
```
curl -k -i https://127.0.0.1/health
curl -i https://surennongmu.com/health
```

### 9.2 静态资源
```
curl -I https://surennongmu.com/static/css/main.css
```

### 9.3 端口暴露面
```
sudo ss -lntp | egrep ':8080|:80|:443'
```

期望
- cloud-server：127.0.0.1:8080
- nginx：0.0.0.0:80、0.0.0.0:443

## 10. 回滚流程
```
sudo ln -sfn /opt/cloud-server/releases/<old> /opt/cloud-server/current
sudo systemctl restart cloud-server
curl -i https://surennongmu.com/health
```

## 11. 常见问题与定位
### 11.1 systemd 启动失败
```
sudo systemctl status cloud-server --no-pager -l
journalctl -u cloud-server -n 200 --no-pager
```
常见原因
- env 缺失导致 panic
- configs/config.yaml 不正确
- ExecStart 路径错误

### 11.2 NGINX reload 失败
```
sudo nginx -t
sudo tail -n 200 /var/log/nginx/surennongmu.error.log
```
常见原因
- 证书目录名不对
- 配置语法错误

### 11.3 /health OK 但 / 404
```
curl -i http://127.0.0.1:8080/
```
说明后端没有注册 GET /，需要补路由或在 NGINX 做跳转。

## 12. 发布检查清单
- /opt/cloud-server/current 是否为软链且指向正确版本
- bin/cloud-server 是否为最新编译时间
- cloud-server.service ExecStart 路径正确
- cloud-server 仅监听 127.0.0.1:8080
- nginx -t 通过并 reload 成功
- https://surennongmu.com/health 返回 200
- static 资源 200 且带缓存头

## 13. 修改策略
### 方案 A 推荐
新建 release 目录，编译后切换软链，再重启后端。  
好处：版本语义不被破坏，可一键回滚。

### 方案 B 不推荐但可用
直接在 current 改代码必须重新编译并重启服务。  
它等价于污染当前 release，回滚与对比会失真。

## 14. NGINX 是否需要重启
- 只改 Go 代码：不需要动 NGINX，重启 cloud-server 即可
- 改了 NGINX 配置：先 nginx -t，再 reload
- 只改静态文件：一般无需 reload，注意浏览器缓存

