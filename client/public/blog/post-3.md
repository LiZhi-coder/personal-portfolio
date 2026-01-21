---
title: Git 配置与基础命令
date: 2025-12-05 10:00:00
excerpt: 从 Git 用户配置到日常提交推送的指南
readingTime: 5 分钟
---

# Git 配置与基础工作流

这篇文章记录了 Git 的用户配置方法以及日常开发中最常用的提交推送流程。
## 1. Git 用户配置

Git 的用户配置分为三个层级，优先级从低到高依次为：系统级 → 用户级 → 项目级。

### 1.1 系统级配置

```bash
# 配置文件位置：/etc/gitconfig
git config --system user.name "你的名字"
git config --system user.email "你的邮箱"
```

> [!NOTE]
> 系统级配置对所有用户生效，几乎不会使用。需要管理员权限。

### 1.2 用户级配置（推荐）

```bash
# 配置文件位置：~/.gitconfig
git config --global user.name "hui-cyber"
git config --global user.email "1597305757@qq.com"
```

> [!TIP]
> 这是最常用的配置方式，对当前系统用户的所有项目生效。

### 1.3 项目级配置

```bash
# 配置文件位置：项目/.git/config
# 需要先进入项目根目录
cd /path/to/your/project

git config user.name "hui-cyber"
git config user.email "1597305757@qq.com"
```

> [!IMPORTANT]
> 项目级配置只对当前项目生效，优先级最高。适用于需要使用不同身份（如公司账号 vs 个人账号）的场景。

## 2. Git 基础工作流

日常开发中，最常用的 Git 操作流程如下：

```
修改代码 → git add → git commit → git pull → git push
```

### 2.1 暂存更改：git add

```bash
# 暂存当前目录下所有已修改或新增的文件
git add .

# 暂存指定文件
git add filename

# 暂存已修改和删除的文件（不含新文件）
git add -u
```

**作用**：将工作区的更改添加到暂存区（staging area），告诉 Git "我准备把这些更改纳入下次提交"。

> [!NOTE]
> - `.` 表示当前目录及所有子目录
> - 被 `.gitignore` 忽略的文件不会被添加
> - 删除的文件需要用 `git add -u` 或 `git rm` 处理

### 2.2 提交更改：git commit

```bash
git commit -m "修复登录页面的表单验证问题"
```

**作用**：将暂存区的内容打包成一个本地提交（commit），并附上说明信息。

> [!TIP]
> **提交信息的最佳实践**
> - 清晰、简洁地描述**为什么改**，而不只是**改了什么**
> - 每个 commit 都有唯一的 hash ID，形成版本历史
> - 这个提交只存在于本地，远程仓库还不知道

### 2.3 拉取远程更新：git pull

```bash
git pull
```

**作用**：从远程仓库拉取最新更改并合并到本地。相当于 `git fetch` + `git merge`。

**为什么这一步很重要？**

- 在你写代码的同时，可能有同事已经推送了新代码
- 如果直接 push，Git 会拒绝（因为历史不一致）
- `git pull` 确保本地代码是最新的，避免冲突或覆盖他人代码

> [!WARNING]
> **合并冲突**：如果你和别人修改了同一个文件的同一部分，且修改内容不同，pull 时可能产生合并冲突（merge conflict），需要手动解决。遇到了再说，不用过于担心。

### 2.4 推送到远程：git push

```bash
git push

# 第一次推送新分支时，需要建立关联
git push -u origin 分支名
```

**作用**：将本地的新提交推送到远程仓库（如 GitHub、GitLab）。

**推送的本质**：告诉远程仓库"把你更新成我本地的样子"，最终远程和本地保持一致。

> [!CAUTION]
> 如果远程仓库有你本地没有的新提交，`git push` 会被拒绝。这是 Git 在提醒你：远程比本地新，请先 `git pull` 同步。

## 3. 完整的提交同步流程

把上面的命令串起来，就是日常开发中最常用的完整流程：

```bash
# 1. 暂存所有更改
git add .

# 2. 提交到本地仓库
git commit -m "你的提交信息"

# 3. 拉取远程最新代码
git pull

# 4. 推送到远程仓库
git push
```

## 4. 流程图

```mermaid
graph LR
    A[工作区] -->|git add| B[暂存区]
    B -->|git commit| C[本地仓库]
    C -->|git push| D[远程仓库]
    D -->|git pull| C
```

## 5. 常见问题

### Q: git add 和 git commit 有什么区别？

- `git add`：把更改放入"购物车"（暂存区）
- `git commit`：结账（真正记录到版本历史）

你可以多次 `add` 后再一次性 `commit`，这样可以把相关的更改打包成一个有意义的提交。

### Q: 为什么要先 pull 再 push？

这是团队协作的基本礼仪。先确认远程有没有新提交，如果有就先合并到本地，解决可能的冲突后再推送，避免覆盖他人的工作。

### Q: 合并冲突怎么解决？

当 Git 无法自动合并时，会在冲突文件中标记出冲突区域：

```
<<<<<<< HEAD
你的代码
=======
别人的代码
>>>>>>> branch-name
```

手动编辑保留正确的代码，删除标记符号，然后重新 `add` 和 `commit` 即可。
