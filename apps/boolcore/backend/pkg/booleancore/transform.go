package booleancore

import (
	"math"
	"runtime"
	"sync"
)

// WalshHadamardTransform 计算函数的沃尔什-哈达玛变换谱.
// 结果会被缓存. 优化返回类型为 int64.
func (f *BooleanFunction) WalshHadamardTransform() []int64 {
	if f.walshSpectrum != nil {
		return f.walshSpectrum // 使用缓存
	}

	length := 1 << f.n
	// s[x] = (-1)^{f(x)}
	s := make([]int64, length)
	for i := 0; i < length; i++ {
		// 从位打包数据中读取第 i 个 bit
		bit := (f.packedTruthTable[i>>6] >> uint(i&63)) & 1
		if bit == 0 {
			s[i] = 1
		} else {
			s[i] = -1
		}
	}

	fwhtInplace(s)      // 执行变换
	f.walshSpectrum = s // 缓存结果
	return f.walshSpectrum
}

// WalshHadamardTransformParallel FWHT 的并行版本.
func (f *BooleanFunction) WalshHadamardTransformParallel(workers int) []int64 {
	// 并行版本不使用缓存，以确保每次都执行计算
	length := 1 << f.n
	s := make([]int64, length)
	for i := 0; i < length; i++ {
		bit := (f.packedTruthTable[i>>6] >> uint(i&63)) & 1
		if bit == 0 {
			s[i] = 1
		} else {
			s[i] = -1
		}
	}
	parallelFWHTInplace(s, workers)
	return s
}

// Autocorrelation 计算函数的自相关谱.
// 采用 float64 中间计算保证精度和范围，并增加缓存
func (f *BooleanFunction) Autocorrelation() []int64 {
	if f.autocorrelationSpectrum != nil {
		return f.autocorrelationSpectrum
	}

	length := 1 << f.n
	s := make([]int64, length)
	for i := 0; i < length; i++ {
		bit := (f.packedTruthTable[i>>6] >> uint(i&63)) & 1
		if bit == 0 {
			s[i] = 1
		} else {
			s[i] = -1
		}
	}

	fwhtInplace(s)

	// 平方操作，转换为 float64 以防止溢出并保证精度
	sSquaredFloat := make([]float64, length)
	for i := 0; i < length; i++ {
		sSquaredFloat[i] = float64(s[i]) * float64(s[i])
	}

	// 在 float64 上再次 FWHT (相当于逆变换)
	fwhtFloatInplace(sSquaredFloat)

	// 归一化并转换回 int64
	res := make([]int64, length)
	nFloat := float64(length)
	for i := 0; i < length; i++ {
		res[i] = int64(math.Round(sSquaredFloat[i] / nFloat))
	}

	f.autocorrelationSpectrum = res // 缓存结果
	return res
}

// --- 私有实现 ( int64 改为 int) ---

func fwhtInplace(data []int64) {
	n := len(data)
	for step := 1; step < n; step <<= 1 {
		for i := 0; i < n; i += step << 1 {
			for j := 0; j < step; j++ {
				a := data[i+j]
				b := data[i+j+step]
				data[i+j] = a + b
				data[i+j+step] = a - b
			}
		}
	}
}

// 浮点数版本用于计算可能超过2的64位整数的函数
func fwhtFloatInplace(data []float64) {
	n := len(data)
	for step := 1; step < n; step <<= 1 {
		for i := 0; i < n; i += step << 1 {
			for j := 0; j < step; j++ {
				a := data[i+j]
				b := data[i+j+step]
				data[i+j] = a + b
				data[i+j+step] = a - b
			}
		}
	}
}

func parallelFWHTInplace(data []int64, workers int) {
	n := len(data)
	if workers <= 1 {
		fwhtInplace(data)
		return
	}
	maxWorkers := runtime.GOMAXPROCS(0)
	if workers > maxWorkers {
		workers = maxWorkers
	}
	for step := 1; step < n; step <<= 1 {
		blockSize := step << 1
		numBlocks := n / blockSize
		wg := sync.WaitGroup{}
		wg.Add(workers)
		for w := 0; w < workers; w++ {
			go func(wid int) {
				defer wg.Done()
				for b := wid; b < numBlocks; b += workers {
					base := b * blockSize
					for j := 0; j < step; j++ {
						a := data[base+j]
						b := data[base+j+step]
						data[base+j] = a + b
						data[base+j+step] = a - b
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

// Fast Möbius Transform 计算 ANF、逆变换
