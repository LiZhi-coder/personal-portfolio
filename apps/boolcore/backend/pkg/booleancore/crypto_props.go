package booleancore

import (
	"fmt"
	"math/bits"
	"strings"
)

// Nonlinearity 计算非线性度.
func (f *BooleanFunction) Nonlinearity() int64 {
	// 公式: (2^(n-1) - max(|WHT_coeffs|)) / 2
	wht := f.WalshHadamardTransform() // 调用方法自动处理缓存

	var maxAbs int64 = 0
	for _, v := range wht {
		absV := v
		if absV < 0 {
			absV = -absV
		}
		if absV > maxAbs {
			maxAbs = absV
		}
	}
	// nonlinearity = 2^(n-1) - maxAbs/2
	return (1 << (f.n - 1)) - maxAbs/2
}

// CorrelationImmunity 计算相关免疫阶数.
// 【重构】采用逐阶验证的清晰逻辑.
func (f *BooleanFunction) CorrelationImmunity() int {
	wht := f.WalshHadamardTransform()
	maxOrder := 0
	for t := 1; t <= f.n; t++ {
		isTOrderImmune := true
		for w := 1; w < len(wht); w++ {
			if bits.OnesCount(uint(w)) <= t {
				if wht[w] != 0 {
					isTOrderImmune = false
					break
				}
			}
		}
		if isTOrderImmune {
			maxOrder = t
		} else {
			break // 一旦在 t 阶不满足，就不可能在 t+1 阶满足
		}
	}
	return maxOrder
}

// ResiliencyOrder 计算弹性阶数. 是在相关免疫的基础上的，如果平衡就是弹性函数
func (f *BooleanFunction) ResiliencyOrder() int {
	if !f.IsBalanced() {
		return -1 // 不平衡的函数没有弹性
	}
	return f.CorrelationImmunity()
}

// IsBent 检查函数是否是 bent 函数.
func (f *BooleanFunction) IsBent() bool {
	if f.n%2 != 0 {
		return false
	}
	wht := f.WalshHadamardTransform()
	// bent 函数的 WHT 的绝对值是常数2^(n/2)
	expectedMagnitude := int64(1 << (f.n / 2))
	for _, v := range wht {
		absV := v
		if absV < 0 {
			absV = -absV
		}
		if absV != expectedMagnitude {
			return false
		}
	}
	return true
}

// SumOfSquareIndicator 计算平方和指标.
func (f *BooleanFunction) SumOfSquareIndicator() int64 {
	ac := f.Autocorrelation()
	var sum int64 = 0
	for _, v := range ac {
		sum += v * v
	}
	return sum
}

// AbsoluteWalshSpectrum 计算绝对 Walsh 谱.
// AbsoluteWalshSpectrum 计算绝对 Walsh 谱的频率分布.
// 返回一个 map，键是谱值的绝对值，值是该绝对值出现的次数.
func (f *BooleanFunction) AbsoluteWalshSpectrum() map[int64]int {
	wht := f.WalshHadamardTransform()
	distribution := make(map[int64]int)
	for _, v := range wht {
		absV := v
		if absV < 0 {
			absV = -absV
		}
		distribution[absV]++
	}
	return distribution
}

// AbsoluteAutocorrelation 计算绝对自相关谱的频率分布.
func (f *BooleanFunction) AbsoluteAutocorrelation() map[int64]int {
	ac := f.Autocorrelation()
	distribution := make(map[int64]int)
	for _, v := range ac {
		absV := v
		if absV < 0 {
			absV = -absV
		}
		distribution[absV]++
	}
	return distribution
}

// AbsoluteIndicator 计算绝对指标 (Delta-uniformity).
// Delta_f = max_{a != 0} |AC_f(a)|.   除去第一个元素
// 这是一个衡量函数抵抗差分攻击能力的重要指标，值越小越好.
func (f *BooleanFunction) AbsoluteIndicator() int64 {
	ac := f.Autocorrelation()
	if len(ac) <= 1 {
		return 0 // 对于n = 0的情况或异常情况
	}
	var maxAbs int64 = 0
	// 循环从索引 1 开始，跳过 AC_f(0)
	for _, v := range ac[1:] { // 跳过 a=0 的位置
		absV := v
		if absV < 0 {
			absV = -absV
		}
		if absV > maxAbs {
			maxAbs = absV
		}
	}
	return maxAbs
}

// DifferentialUniformity 计算差分均匀度 (delta_f).
// 该值与绝对指标 (Delta_f) 存在线性关系: delta_f = 2^(n-1) + Delta_f / 2.
// 它同样衡量函数抵抗差分攻击的能力，值越小越好.
// 最小值为 2^(n-1)，达到该值的函数称为完美非线性函数(PN)，对于布尔函数即为Bent函数.
func (f *BooleanFunction) DifferentialUniformity() int64 {
	// 首先，调用已经写好的函数计算 Delta_f
	deltaF := f.AbsoluteIndicator()

	// 然后，应用转换公式
	// delta_f = 2^(n-1) + Delta_f / 2
	return (1 << (f.n - 1)) + (deltaF / 2)
}

// IsRotationSymmetric 检查函数是否具有旋转对称性.检查函数是否是旋转对称函数 .
func (f *BooleanFunction) IsRotationSymmetric() bool {
	// 【文档】此处的循环移位操作基于 x0 为最低有效位 (LSB) 的约定.
	tt := f.TruthTable()
	length := len(tt)
	mask := length - 1 // 例如 n=3, mask=7 (二进制 111)

	for i := 0; i < length; i++ {
		// 循环左移一位: (i << 1) | (i >> (n-1))
		rotatedI := ((i << 1) | (i >> (f.n - 1))) & mask
		if tt[i] != tt[rotatedI] {
			return false
		}
	}
	return true

}

// TransparencyOrder 计算透明度阶.
// To(f) = 1 - (1/(2^n(2^n-1))) * ∑|ac[a]| for a != 0
func (f *BooleanFunction) TransparencyOrder() float64 {
	ac := f.Autocorrelation()
	length := 1 << f.n

	if length <= 1 {
		return 1.0 // 对于 n=0 或 n=1 的平凡情况
	}

	var sumDelta int64 = 0
	// 累加所有非零位置的自相关谱的绝对值
	for i := 1; i < length; i++ { // 从自相关的索引 1 开始，跳过 a=0
		absV := ac[i]
		if absV < 0 {
			absV = -absV
		}
		sumDelta += absV
	}

	denominator := float64(length * (length - 1))
	if denominator == 0 {
		return 1.0
	}

	to := 1.0 - (float64(sumDelta) / denominator)
	return to
}

// AlgebraicImmunity /*以下是计算代数免疫度和零化因子的函数很长*/
// AlgebraicImmunity 计算代数免疫度.
// 参数 findAnnihilator 控制是否需要计算并返回一个具体的最低次零化子表达式.
// 返回值: (代数免疫度, 零化子ANF字符串, 错误)
// ...existing code...

// BitMatrix 使用 uint64 切片来紧凑地存储二元矩阵，以进行高效的位运算.
type BitMatrix struct {
    rows        int
    cols        int
    wordsPerRow int      // 每行需要多少个 uint64
    data        []uint64
}

// NewBitMatrix 创建一个新的位矩阵.
func NewBitMatrix(rows, cols int) *BitMatrix {
    wordsPerRow := (cols + 63) / 64
    return &BitMatrix{
        rows:        rows,
        cols:        cols,
        wordsPerRow: wordsPerRow,
        data:        make([]uint64, rows*wordsPerRow),
    }
}

// Set 设置矩阵在 (r, c) 位置的比特值.
func (m *BitMatrix) Set(r, c int, val byte) {
    wordIndex := r*m.wordsPerRow + c/64
    bitIndex := uint(c % 64)
    if val == 1 {
        m.data[wordIndex] |= (1 << bitIndex)
    } else {
        m.data[wordIndex] &= ^(1 << bitIndex)
    }
}

// Get 获取矩阵在 (r, c) 位置的比特值.
func (m *BitMatrix) Get(r, c int) byte {
    wordIndex := r*m.wordsPerRow + c/64
    bitIndex := uint(c % 64)
    if (m.data[wordIndex]>>bitIndex)&1 == 1 {
        return 1
    }
    return 0
}

// SwapRows 交换两行.
func (m *BitMatrix) SwapRows(r1, r2 int) {
    if r1 == r2 {
        return
    }
    start1 := r1 * m.wordsPerRow
    start2 := r2 * m.wordsPerRow
    for i := 0; i < m.wordsPerRow; i++ {
        m.data[start1+i], m.data[start2+i] = m.data[start2+i], m.data[start1+i]
    }
}

// XorRow 将 srcRow 的值异或到 dstRow 上. 这是性能提升的关键.
func (m *BitMatrix) XorRow(dstRow, srcRow int) {
    startDst := dstRow * m.wordsPerRow
    startSrc := srcRow * m.wordsPerRow
    for i := 0; i < m.wordsPerRow; i++ {
        m.data[startDst+i] ^= m.data[startSrc+i]
    }
}

// Clone 创建矩阵的深拷贝
func (m *BitMatrix) Clone() *BitMatrix {
    newMatrix := &BitMatrix{
        rows:        m.rows,
        cols:        m.cols,
        wordsPerRow: m.wordsPerRow,
        data:        make([]uint64, len(m.data)),
    }
    copy(newMatrix.data, m.data)
    return newMatrix
}

// computeRREF_GF2 在 BitMatrix 上执行高斯消元法，效率极高.
// 返回矩阵的秩
func computeRREF_GF2(m *BitMatrix) int {
    rank := 0
    pivotRow := 0
    for col := 0; col < m.cols && pivotRow < m.rows; col++ {
        // 寻找主元
        i := pivotRow
        for i < m.rows && m.Get(i, col) == 0 {
            i++
        }

        if i < m.rows { // 找到主元
            m.SwapRows(pivotRow, i)

            // 消去当前列的其他行的1
            for j := 0; j < m.rows; j++ {
                if j != pivotRow && m.Get(j, col) == 1 {
                    // 这里的 XorRow 函数会一次性处理 64 bits!
                    m.XorRow(j, pivotRow)
                }
            }
            pivotRow++
        }
    }
    rank = pivotRow
    return rank
}


// computeRREF_GF2WithSolve 执行高斯消元并返回简化行阶梯形矩阵用于求解
func computeRREF_GF2WithSolve(m *BitMatrix) (int, *BitMatrix) {
    rrefMatrix := m.Clone()
    rank := computeRREF_GF2(rrefMatrix) // 使用RREF版本
    return rank, rrefMatrix
}

// AlgebraicImmunity 计算代数免疫度.
// 参数 findAnnihilator 控制是否需要计算并返回一个具体的最低次零化子表达式.
// 返回值: (代数免疫度, 零化子ANF字符串, 错误)
func (f *BooleanFunction) AlgebraicImmunity(findAnnihilator bool) (int, string, error) {
    // 边界情况处理
    if f.n == 0 {
        return 0, "", nil
    }
    if f.n == 1 {
        tt := truthTableFromUint64Slice(f.packedTruthTable, 2)
        if tt[0] == tt[1] {
            // 常数函数
            return 1, "x0", nil
        }
        // 非常数函数,AI=1
        return 1, "x0", nil
    }

    // 理论上代数免疫度的最大值为 ceil(n/2)
    maxDegreeToCheck := (f.n + 1) / 2

    tt := truthTableFromUint64Slice(f.packedTruthTable, 1<<f.n)

    // 计算汉明重量以决定检查顺序
    weight := 0
    for _, v := range tt {
        if v == 1 {
            weight++
        }
    }

    for d := 1; d <= maxDegreeToCheck; d++ {
        // 决定检查顺序：优先检查支撑集更小的函数
        var checkOrder [2]bool
        if weight <= (1 << (f.n - 1)) {
            checkOrder = [2]bool{true, false} // f 的支撑集更小，先检查 f
        } else {
            checkOrder = [2]bool{false, true} // f+1 的支撑集更小，先检查 f+1
        }

        for _, checkF := range checkOrder {
            if findAnnihilator {
                // 需要具体的零化子，使用完整版本
                found, annihilator, err := findLowestDegreeAnnihilatorFull(f.n, d, tt, checkF)
                if err != nil {
                    return -1, "", err
                }
                if found {
                    return d, annihilator, nil
                }
            } else {
                // 只需要判断存在性，使用快速版本
                found, err := findAnnihilatorExists(f.n, d, tt, checkF)
                if err != nil {
                    return -1, "", err
                }
                if found {
                    return d, "", nil
                }
            }
        }
    }

    return maxDegreeToCheck, "", nil
}

// findAnnihilatorExists 快速判断是否存在零化子（不求解具体表达式）
func findAnnihilatorExists(n, d int, tt []byte, useSupportOfF bool) (bool, error) {
    // 1. 生成单项式
    monomials := make([]int, 0)
    for i := 0; i < (1 << n); i++ {
        if bits.OnesCount(uint(i)) <= d {
            monomials = append(monomials, i)
        }
    }
    numVars := len(monomials)

    // 2. 构建支撑集
    support := make([]int, 0)
    targetVal := byte(1)
    if !useSupportOfF {
        targetVal = byte(0)
    }
    for i, val := range tt {
        if val == targetVal {
            support = append(support, i)
        }
    }

    // 常数函数特殊处理
    if len(support) == 0 {
        return d >= 1, nil
    }

    numEqs := len(support)

    // 如果方程数少于未知数个数，必然存在非零解
    if numEqs < numVars {
        return true, nil
    }

    // 3. 构建 BitMatrix
    matrix := NewBitMatrix(numEqs, numVars)
    for i := 0; i < numEqs; i++ {
        inputVec := support[i]
        for j := 0; j < numVars; j++ {
            monomial := monomials[j]
            if (inputVec & monomial) == monomial {
                matrix.Set(i, j, 1)
            }
        }
    }

    // 4. 调用优化版的高斯消元
    rank := computeRREF_GF2(matrix)

    return rank < numVars, nil
}

// findLowestDegreeAnnihilatorFull 完整版本，计算具体的零化子表达式
func findLowestDegreeAnnihilatorFull(n, d int, tt []byte, useSupportOfF bool) (bool, string, error) {
    // 1. 生成单项式
    monomials := make([]int, 0)
    for i := 0; i < (1 << n); i++ {
        if bits.OnesCount(uint(i)) <= d {
            monomials = append(monomials, i)
        }
    }
    numVars := len(monomials)

    // 2. 构建支撑集
    support := make([]int, 0)
    targetVal := byte(1)
    if !useSupportOfF {
        targetVal = byte(0)
    }
    for i, val := range tt {
        if val == targetVal {
            support = append(support, i)
        }
    }

    // 常数函数特殊处理
    if len(support) == 0 {
        if d >= 1 {
            return true, "x0", nil
        }
        return false, "", nil
    }

    numEqs := len(support)

    // 3. 构建 BitMatrix
    matrix := NewBitMatrix(numEqs, numVars)
    for i := 0; i < numEqs; i++ {
        inputVec := support[i]
        for j := 0; j < numVars; j++ {
            monomial := monomials[j]
            if (inputVec & monomial) == monomial {
                matrix.Set(i, j, 1)
            }
        }
    }

    // 4. 调用带求解的高斯消元
    rank, rrefMatrix := computeRREF_GF2WithSolve(matrix)

    if rank < numVars {
        // 从RREF矩阵中求解
        solution := solveRREFFromBitMatrix(rrefMatrix, numVars, rank)
        annihilatorStr := formatAnnihilator(solution, monomials, n)
        return true, annihilatorStr, nil
    }

    return false, "", nil
}

// solveRREFFromBitMatrix 从 BitMatrix 格式的 RREF 矩阵中找到一个非零特解
func solveRREFFromBitMatrix(rref *BitMatrix, numVars, rank int) []byte {
    solution := make([]byte, numVars)
    pivotCols := make([]int, rank)
    for i := range pivotCols {
        pivotCols[i] = -1
    }

    // 识别主元列
    for r := 0; r < rank && r < rref.rows; r++ {
        for c := 0; c < rref.cols; c++ {
            if rref.Get(r, c) == 1 {
                pivotCols[r] = c
                break
            }
        }
    }

    // 找到第一个自由变量（非主元列），并将其设为1
    isPivot := make([]bool, numVars)
    for _, c := range pivotCols {
        if c != -1 {
            isPivot[c] = true
        }
    }

    freeVarIndex := -1
    for c := numVars - 1; c >= 0; c-- {
        if !isPivot[c] {
            freeVarIndex = c
            break
        }
    }

    if freeVarIndex != -1 {
        solution[freeVarIndex] = 1
    } else {
        // 没有自由变量，返回零解
        return solution
    }

    // 回代求解主元变量
    for r := rank - 1; r >= 0; r-- {
        pivotCol := pivotCols[r]
        if pivotCol == -1 {
            continue
        }

        var sum byte = 0
        for c := pivotCol + 1; c < numVars; c++ {
            sum ^= rref.Get(r, c) * solution[c]
        }
        solution[pivotCol] = sum
    }

    return solution
}

// formatAnnihilator 将解向量格式化为人类可读的 ANF 字符串.
func formatAnnihilator(solution []byte, monomials []int, n int) string {
    // 创建一个长度为 2^n 的完整 ANF 系数向量
    hCoeffs := make([]byte, 1<<n)
    for i, coeff := range solution {
        if coeff == 1 {
            monomialIndex := monomials[i]
            hCoeffs[monomialIndex] = 1
        }
    }

    // 格式化为字符串
    var terms []string
    for i, coeff := range hCoeffs {
        if coeff == 1 {
            if i == 0 {
                terms = append(terms, "1")
            } else {
                var termParts []string
                for j := 0; j < n; j++ {
                    if (i>>j)&1 == 1 {
                        termParts = append(termParts, fmt.Sprintf("x%d", j))
                    }
                }
                terms = append(terms, strings.Join(termParts, "*"))
            }
        }
    }

    if len(terms) == 0 {
        return "0"
    }
    return strings.Join(terms, " + ")
}

/*上面包裹的是计算代数免疫度和零化因子的函数*/


