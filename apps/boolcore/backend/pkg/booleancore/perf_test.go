package booleancore

import (
	"testing"
)

// BenchmarkCoreAnalyze 针对核心布尔函数分析的基准测试，不经过 HTTP 层。
// 覆盖主要性质计算：ANF/次数、WHT、自相关、非线性度、相关免疫、弹性、Bent、
// 平方和指标、旋转对称、绝对指标、差分均匀度、代数免疫度（快速版本）。
func BenchmarkCoreAnalyze(b *testing.B) {
	// 使用一个中等规模用例（n=6）以兼顾速度与覆盖度
	n := 6
	val := uint64(123456)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf, err := NewFromInt(val, n)
		if err != nil {
			b.Fatalf("NewFromInt error: %v", err)
		}

		// 代数正规式与次数
		_ = bf.AlgebraicNormalForm()
		_ = bf.AlgebraicDegree()

		// Walsh 与自相关
		_ = bf.WalshHadamardTransform()
		_ = bf.Autocorrelation()

		// 密码学性质
		_ = bf.Nonlinearity()
		_ = bf.CorrelationImmunity()
		_ = bf.ResiliencyOrder()
		_ = bf.IsBent()
		_ = bf.SumOfSquareIndicator()
		_ = bf.IsRotationSymmetric()
		_ = bf.AbsoluteWalshSpectrum()
		_ = bf.AbsoluteAutocorrelation()
		_ = bf.AbsoluteIndicator()
		_ = bf.DifferentialUniformity()
		// 代数免疫度（快速版本，不求零化子表达式）
		_, _, _ = bf.AlgebraicImmunity(false)
	}
}
