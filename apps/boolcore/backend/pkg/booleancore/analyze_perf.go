package booleancore

import (
	"time"
)

// AnalyzeResult 汇总布尔函数的各项性质，便于一次性计算和比较。
type AnalyzeResult struct {
	N                               int
	TruthTable                      []byte
	HammingWeight                   int
	IsBalanced                      bool
	ANF                             string
	AlgebraicDegree                 int
	WalshSpectrum                   []int64
	AutocorrelationSpectrum         []int64
	TransparencyOrder               float64
	Nonlinearity                    int64
	CorrelationImmunity             int
	ResiliencyOrder                 int
	IsBent                          bool
	SumOfSquareIndicator            int64
	IsRotationSymmetric             bool
	AbsoluteWalshSpectrum           map[int64]int
	AbsoluteAutocorrelationSpectrum map[int64]int
	AbsoluteIndicator               int64
	DifferentialUniformity          int64
	AlgebraicImmunity               int
}

// AnalyzeAll 计算所有核心性质（快速版本：代数免疫度不求零化子表达式）。
func AnalyzeAll(bf *BooleanFunction) AnalyzeResult {
	// 注意：内部方法已经带有缓存（如WHT/自相关），多处复用不会重复计算。
	res := AnalyzeResult{N: bf.N(), TruthTable: bf.TruthTable()}
	res.HammingWeight = bf.HammingWeight()
	res.IsBalanced = bf.IsBalanced()
	res.ANF = bf.AlgebraicNormalForm()
	res.AlgebraicDegree = bf.AlgebraicDegree()
	res.WalshSpectrum = bf.WalshHadamardTransform()
	res.AutocorrelationSpectrum = bf.Autocorrelation()
	res.TransparencyOrder = bf.TransparencyOrder()
	res.Nonlinearity = bf.Nonlinearity()
	res.CorrelationImmunity = bf.CorrelationImmunity()
	res.ResiliencyOrder = bf.ResiliencyOrder()
	res.IsBent = bf.IsBent()
	res.SumOfSquareIndicator = bf.SumOfSquareIndicator()
	res.IsRotationSymmetric = bf.IsRotationSymmetric()
	res.AbsoluteWalshSpectrum = bf.AbsoluteWalshSpectrum()
	res.AbsoluteAutocorrelationSpectrum = bf.AbsoluteAutocorrelation()
	res.AbsoluteIndicator = bf.AbsoluteIndicator()
	res.DifferentialUniformity = bf.DifferentialUniformity()
	if ai, _, err := bf.AlgebraicImmunity(false); err == nil {
		res.AlgebraicImmunity = ai
	} else {
		res.AlgebraicImmunity = -1
	}
	return res
}

// AnalyzeAllTimed 在 AnalyzeAll 基础上返回每一步耗时与总耗时，方便与 SageMath 做时间对比。
func AnalyzeAllTimed(bf *BooleanFunction) (AnalyzeResult, map[string]time.Duration, time.Duration) {
	timings := make(map[string]time.Duration)
	startTotal := time.Now()
	step := func(name string, fn func()) {
		s := time.Now()
		fn()
		timings[name] = time.Since(s)
	}

	var res AnalyzeResult
	step("construct_result", func() {
		res.N = bf.N()
		res.TruthTable = bf.TruthTable()
	})
	step("hamming_weight", func() { res.HammingWeight = bf.HammingWeight() })
	step("is_balanced", func() { res.IsBalanced = bf.IsBalanced() })
	step("anf", func() { res.ANF = bf.AlgebraicNormalForm() })
	step("algebraic_degree", func() { res.AlgebraicDegree = bf.AlgebraicDegree() })
	step("walsh_hadamard", func() { res.WalshSpectrum = bf.WalshHadamardTransform() })
	step("autocorrelation", func() { res.AutocorrelationSpectrum = bf.Autocorrelation() })
	step("transparency_order", func() { res.TransparencyOrder = bf.TransparencyOrder() })
	step("nonlinearity", func() { res.Nonlinearity = bf.Nonlinearity() })
	step("correlation_immunity", func() { res.CorrelationImmunity = bf.CorrelationImmunity() })
	step("resiliency_order", func() { res.ResiliencyOrder = bf.ResiliencyOrder() })
	step("is_bent", func() { res.IsBent = bf.IsBent() })
	step("sum_of_square_indicator", func() { res.SumOfSquareIndicator = bf.SumOfSquareIndicator() })
	step("rotation_symmetric", func() { res.IsRotationSymmetric = bf.IsRotationSymmetric() })
	step("absolute_walsh_spectrum", func() { res.AbsoluteWalshSpectrum = bf.AbsoluteWalshSpectrum() })
	step("absolute_autocorr_spectrum", func() { res.AbsoluteAutocorrelationSpectrum = bf.AbsoluteAutocorrelation() })
	step("absolute_indicator", func() { res.AbsoluteIndicator = bf.AbsoluteIndicator() })
	step("differential_uniformity", func() { res.DifferentialUniformity = bf.DifferentialUniformity() })
	step("algebraic_immunity_fast", func() {
		if ai, _, err := bf.AlgebraicImmunity(false); err == nil {
			res.AlgebraicImmunity = ai
		} else {
			res.AlgebraicImmunity = -1
		}
	})

	total := time.Since(startTotal)
	return res, timings, total
}
