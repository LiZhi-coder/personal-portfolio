package booleancore

import (
	"fmt"
	"math/bits"
	"strings"
)

// 这里的计算的是一些基础性质，简单的

// HammingWeight 计算真值表中1的数量 因为使用了uint64打包所以可以使用CPU popcount加速
func (f *BooleanFunction) HammingWeight() int {
	weight := 0
	for _, v := range f.packedTruthTable {
		weight += bits.OnesCount64(v) // 使用OneCount64
	}
	return weight
}

// IsBalanced 检查函数是否是平衡的.
// 平衡意味着真值表中 0 和 1 的数量相等.
func (f *BooleanFunction) IsBalanced() bool {
	// 真值表从长度为2^n
	return f.HammingWeight() == (1 << (f.n - 1)) // 平衡所以是一半
}

// AlgebraicNormalFormCoefficients 使用快速莫比乌斯变换 (FMT) 计算 ANF 系数.
// 返回的也是一个真值表形式的 []byte.
func (f *BooleanFunction) AlgebraicNormalFormCoefficients() []byte {
	if f.anfCoefficients != nil {
		return f.anfCoefficients
	}

	ttCopy := f.TruthTable() // 调用TruthTable()函数将真值表复制一份
	fmtInplace(ttCopy)       // 使用fmtInplace()函数进行快速莫比乌斯变换

	f.anfCoefficients = ttCopy // 缓存结果
	return f.anfCoefficients
}

// AlgebraicNormalForm 将ANF系数转换为我们可读的字符串形式
// TODO:检查这里应该搞的是x0最高位，但sage中的是x0最低位
func (f *BooleanFunction) AlgebraicNormalForm() string {
	coeffs := f.AlgebraicNormalFormCoefficients()
	//  length := len(coeffs)
	var terms []string
	for i, coeff := range coeffs {
		if coeff == 1 {
			if i == 0 {
				terms = append(terms, "1")
			} else {
				var termParts []string
				for j := 0; j < f.n; j++ {
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

// AlgebraicDegree 计算代数次数.
func (f *BooleanFunction) AlgebraicDegree() int {
	coeffs := f.AlgebraicNormalFormCoefficients()
	maxDegree := 0
	for i, coeff := range coeffs {
		if coeff == 1 {
			// 计算 i 的二进制表示中 1 的个数 (popcount)
			degree := bits.OnesCount(uint(i))
			if degree > maxDegree {
				maxDegree = degree
			}
		}
	}
	return maxDegree
}

// --- 私有实现
// fmtInplace 是快速莫比乌斯变换的原地实现.
func fmtInplace(tt []byte) {
	n := len(tt)
	k := bits.TrailingZeros(uint(n))
	for i := 0; i < k; i++ {
		bit := 1 << i
		for j := 0; j < n; j++ {
			if (j & bit) != 0 {
				tt[j] ^= tt[j^bit]
			}
		}
	}
}

// fmtInverseInplace 是快速莫比乌斯逆变换的原地实现.
// 将 ANF 系数转换回真值表
func fmtInverseInplace(anf []byte) {
	n := len(anf)
	k := bits.TrailingZeros(uint(n))
	// 逆变换：反向执行相同的操作
	for i := k - 1; i >= 0; i-- {
		bit := 1 << i
		for j := 0; j < n; j++ {
			if (j & bit) != 0 {
				anf[j] ^= anf[j^bit]
			}
		}
	}
}

// TODO: 还有其他很多性质到时候再说，先把框架搭起来
