---
title: Boolean Function Core
date: 2025-12-28 14:30:00
tags: [Go, Backend, 并行计算, Algorithm Optimization]
link: https://github.com/hui-cyber/BoolCore
description: 现如今关于Boolean Function的搜索的研究做的人很多，但是关于在搜索的过程中肯定是要计算其性质的，目前的网络资源中我目前只看到sagemath是有这些性质计算的官方库，但其效率一般而不全，而且没有可视化。
---



# 项目一：BoolCore - 高性能布尔函数密码学性质分析平台

项目关键词：Golang, 高性能计算, 密码学, 并行计算, Vue3  
源码仓库：https://github.com/hui-cyber/BoolCore

## 项目背景

在分组密码 Block Cipher 的设计与安全性评估中，布尔函数 Boolean Function 的密码学性质如非线性度、代数免疫度起着决定性作用。

然而在面对高元 n > 16 的布尔函数时，学术界主流工具 SageMath 基于 Python/C 存在明显性能瓶颈。一次针对 9 元布尔函数的代数免疫度计算中，SageMath 耗时约 89ms，难以满足大规模函数筛选的工程需求。

## 项目概述

 BoolCore，一个基于 Go 的高性能分析引擎，将核心计算从数据库与脚本层剥离出来。通过位打包和并行计算优化，实现了比 SageMath 快两个数量级的计算性能。

## 主要功能

- 计算 16 类性质：非线性度、代数次数、代数免疫度、相关免疫、Bent 判定、谱分布等
- 多输入形式：真值表、十六进制 Hex、整数、ANF 表达式
- 频谱可视化：Walsh 谱与自相关谱展示
- 多格式导出：JSON、Matlab 矩阵与 SageMath 代码

## 技术栈

- 前端：Vue 3 + Vite + Element Plus，负责频谱可视化与输入解析
- 后端：Go Gin Framework，无状态计算服务
- 核心计算层：纯 Go 算法库，包含并行版 FWHT 与位操作库


![ANF图](/images/projects/project-1/image.png)


## 开发过程与挑战

### 1. 位打包 Bit-packing

传统真值表通常使用 []bool 或 []byte 存储，空间利用率极低。我基于 uint64 实现位打包存储，每个 uint64 存储 64 个真值点。

- 内存优化：占用降低 98.4%，约为 1/64  
- 计算加速：利用 CPU popcount 指令进行汉明重量并行计算

```go
// 位打包存储实现
func uint64SliceFromTruthTable(tt []byte) []uint64 {
    n := len(tt)
    words := (n + 63) / 64
    out := make([]uint64, words)
    for i := 0; i < n; i++ {
        if tt[i] != 0 {
            out[i>>6] |= 1 << uint(i&63)
        }
    }
    return out
}
```

### 2. 并行 FWHT

代数免疫度等高阶性质计算瓶颈在快速沃尔什变换 FWHT。我用 Go 的 GMP 调度模型将递归算法重构为分治并行算法。

- 根据 CPU 核心数切分任务  
- sync.WaitGroup 控制并发，显著减少大数组遍历开销

### 3. 跨语言对齐与数学验证

实现代数正规式 ANF 转换时，发现 Go 输出与 SageMath 不一致。

- 问题排查：定位到变量序定义差异  
- 解决方案：实现对偶验证机制，将 ANF 转回真值表比对，并通过 Benchmark 脚本确保 18 种性质计算与 SageMath 一致

## 项目成果

性能对比，9 元布尔函数 Hex 输入实测：

| 计算指标 | SageMath Python/C | BoolCore Go | 性能提升 |
| --- | --- | --- | --- |
| 代数免疫度 AI | 89,716 μs | **531 μs** | **约 168 倍** |
| 代数正规式 ANF | 16,752 μs | **166 μs** | **约 100 倍** |
| 非线性度 NL | 29 μs | **0.3 μs** | **约 96 倍** |
| 总耗时 | > 100 ms | **1.38 ms** | **极速响应** |

注：测试环境为同一台机器，单次运行重复 100000 次取平均值。
