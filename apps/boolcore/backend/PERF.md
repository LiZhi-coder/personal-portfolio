# BoolCore 后端性能对比与测量

本文档说明如何测量纯后端函数的分析耗时，并与 SageMath 进行可复现实验的对比。

## 一、快速测量（纯后端，无 HTTP 开销）

两种方式：

- 基准测试（建议用于稳定/持续对比）：

```bash
cd backend/pkg/booleancore
go test -bench=. -benchmem -run ^$
```

- 小型 CLI（适合像 `start := time.Now()` 的一次性测量与分步耗时）：

```bash
cd backend
go run ./cmd/perf -type int -n 6 -int 123456 -repeat 1 -format text
```

参数：
- `-type`：`int|hex|anf`
- `-n`：变量个数
- `-int` / `-hex` / `-anf`：输入的值或表达式
- `-repeat`：重复次数（>1 时可热缓存）
- `-format`：`text|json`

CLI 输出包含：
- 总耗时（ms）
- 分步骤耗时（ANF、WHT、自相关、非线性度等）
- 关键属性（非线性度、代数次数、是否 Bent 等）

> 代码参考：`pkg/booleancore/AnalyzeAll{,Timed}` 与 `cmd/perf/main.go`

## 二、在你代码里直接测量（与 `time.Now()` 风格一致）

```go
start := time.Now()
bf, _ := booleancore.NewFromInt(123456, 6)
res, stepDur, total := booleancore.AnalyzeAllTimed(bf)
fmt.Println("total:", total)
fmt.Println("nonlinearity:", res.Nonlinearity)
fmt.Println("step walsh:", stepDur["walsh_hadamard"]) // 更多键见代码
```

## 三、与 SageMath 的对比思路

确保输入与索引规则一致：本项目真值表索引与 Sage 的 BooleanFunction 保持一致（x0 为最低位）。

SageMath 示例（保存为 `compare.sage`）：

```python
# 在 Sage 环境下运行：sage -python compare.sage
from sage.all import *

def boolean_from_int(n, val):
    length = 1 << n
    tt = [(val >> i) & 1 for i in range(length)]
    return BooleanFunction(tt)

n = 6
val = 123456

start = time.time()
f = boolean_from_int(n, val)
# 常见指标
anf = f.algebraic_normal_form()
deg = f.algebraic_degree()
wht = f.walsh_spectrum()
auto = f.autocorrelation_spectrum()
nonlin = f.nonlinearity()
# 相关免疫、弹性需要额外实现或使用现有包函数（若有）
elapsed = time.time() - start

print({
    'n': n,
    'elapsed_s': elapsed,
    'deg': deg,
    'nonlinearity': nonlin,
    'walsh_len': len(wht),
    'autocorr_len': len(auto),
})
```

对齐方式：
- 用相同的输入构造：在 Go 与 Sage 都用 `n=6, int=123456`。
- 采集时间：
  - Go：`cmd/perf` 的 `total_ms` 与分步 `steps`；或 `AnalyzeAllTimed` 的 `total`。
  - Sage：`elapsed_s`（注意单位差异）。
- 关键值校验：代数次数、非线性度、谱长度等应一致（数值本身一致，时间可能因平台与实现差异不同）。

## 四、注意事项与已知问题

- API 层（HTTP）基准会包含框架与 JSON 序列化开销，若仅对算法性能感兴趣，请使用本节工具进行“纯后端”测量。
- 当前仓库中 `backend/api_test.go` 含有部分失败用例（与 Bent 判断/输入检查相关），这与本文档新增的纯后端性能测量无直接关系；如需 CI 绿色，请先修复对应测试或调整期望。
- n 较大时（如 8+），谱计算会更耗时；本库内部对 WHT/自相关做了结果缓存，复用同一实例可减少重复成本。

## 五、扩展

- 需要更精细的阶段定义或新增指标时，请在 `AnalyzeAll{,Timed}` 中集中维护，CLI 会自动受益。
- 若要并行 WHT，可使用 `WalshHadamardTransformParallel(workers)` 自行替换并测量效果。
