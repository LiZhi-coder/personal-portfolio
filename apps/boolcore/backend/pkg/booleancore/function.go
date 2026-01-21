package booleancore

import (
	"errors"
	"fmt"
	"math/big"
	"math/bits"
	"strconv"
	"strings"
)

// 支持的构造方式包括：
//   - FromTruthTable([]byte)   真值表构造
//   - FromInt(int)             整数构造（以二进制位为真值）
//   - FromANF(string)          代数正规型构造，如 "x0 + x1*x2 + 1"
//   - FromHex(string)          十六进制构造
//   - FromExpr(string)         布尔表达式构造，如 "(x0 and not x1) xor x2"

//
// 真值表的索引规则与 SageMath 中的 BooleanFunction 保持一致：
//   truthTable[i] = f(x0, x1, ..., xn-1)
// 其中 i 的二进制表示对应输入向量：
//   i = x0 + 2*x1 + 4*x2 + ... + 2^(n-1)*x(n-1)
// 即最低位对应第一个变量 x0。
//
// 例如，对于 n = 3：
// 索引  i  二进制  输入(x2,x1,x0)
//     0    000     f(0,0,0)
//     1    001     f(1,0,0)
//     2    010     f(0,1,0)
//     3    011     f(1,1,0)
//     4    100     f(0,0,1)
//     5    101     f(1,0,1)
//     6    110     f(0,1,1)
//     7    111     f(1,1,1)

// BooleanFunction 代表一个 n 元布尔函数.
// 内部使用位打包的 uint64 切片来高效存储真值表.进行修改
type BooleanFunction struct {
	n                int      // n 元
	packedTruthTable []uint64 // 【升级】使用位打包存储真值表,每个元素存64个布尔值
	// 缓存一些计算结果，避免重复计算
	walshSpectrum           []int64 // 缓存Walsh 频谱
	anfCoefficients         []byte  // 缓存 ANF 系数其实，该方法是布尔函数一种较为重要的表达方法
	autocorrelationSpectrum []int64 // 缓存自相关频谱
}

// NewFromTruthTable 通过真值表创建一个布尔函数.
// 真值表长度必须为 2 的幂，且仅包含 0 或 1。
// 为保证封装安全，输入切片会被复制。

func NewFromTruthTable(tt []byte) (*BooleanFunction, error) {
	length := len(tt)
	if length == 0 || length&(length-1) != 0 {
		return nil, errors.New("truth table length mus t be a power of 2")
	}

	for _, v := range tt {
		if v != 0 && v != 1 {
			return nil, fmt.Errorf("truth table can only contain 0 or 1, found %d", v)
		}
	}

	// bits.TrailingZeros a.k.a Log2 for powers of 2
	n := bits.TrailingZeros(uint(length)) //计算真值表长度为 2ⁿ 时的变量个数 n

	// 将真值表转为位打包的 uint64 切片
	packed := uint64SliceFromTruthTable(tt)

	return &BooleanFunction{
		n:                n,
		packedTruthTable: packed,
	}, nil
}

// TODO:增加其他的构造函数，比如代数标准型和十六进制

// NewFromInt 通过整数创建一个 n 元布尔函数.
// 整数的低位到高位对应真值表的索引 0 到 2^n - 1.
// 例如: n=3, val=170 (二进制 10101010), 真值表为 [0, 1, 0, 1, 0, 1, 0, 1]
// 整数的最低有效位 (LSB - Least Significant Bit) 对应真值表索引为 0 的位置
func NewFromInt(value uint64, n int) (*BooleanFunction, error) {
	if n <= 0 || n > 6 { // 对于这个整数代表真值表，uint64 最多支持 n=6 (2^6=64)
		return nil, fmt.Errorf("n must be between 1 and 6 for uint64 input, got %d", n)
	}

	length := 1 << n
	tt := make([]byte, length)
	for i := 0; i < length; i++ {
		// 检查val的第i位
		if (value>>i)&1 == 1 {
			tt[i] = 1
		} else {
			tt[i] = 0
		}
	}
	// 因为函数内部也是将整数先转为[]byte数组的真值表，所以可以直接调用这个函数
	// ... 其他构造函数类似地在最后调用 NewFromTruthTable ...
	return NewFromTruthTable(tt)
}

// NewFromHex 通过十六进制字符串创建一个 n 元布尔函数.
// 使用 big.Int 支持超过 64 位的任意长度.
func NewFromHex(hexString string, n int) (*BooleanFunction, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be positive, got %d", n)
	}
	// 若有0x 前缀，则移除
	cleanHexString := strings.TrimPrefix(hexString, "0x") // TrimPrefix 区分大小写， "0X" 不会被去掉，后面再改

	// 使用 big.Int来处理很长的十六进制字符串
	val := new(big.Int)
	_, success := val.SetString(cleanHexString, 16)
	if !success {
		return nil, fmt.Errorf("invalid hex string: %s", hexString)
	}
	length := 1 << n
	tt := make([]byte, length)

	for i := 0; i < length; i++ {
		// big.Int.Bit(i) 检查第 i 位是否为 1
		if val.Bit(i) == 1 {
			tt[i] = 1
		} else {
			tt[i] = 0
		}
	}

	// 因为这里的是先将HEX转换为[]byte数组的真值表，所以可以直接调用第一个函数
	return NewFromTruthTable(tt)
}

// NewFromANF 通过代数正规式 (ANF) 字符串创建布尔函数.
func NewFromANF(n int, anfString string) (*BooleanFunction, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be positive, got %d", n)
	}

	// 解析 ANF 字符串并转换为系数向量
	coefficients, err := parseANF(anfString, n)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ANF: %v", err)
	}

	// 使用莫比乌斯逆变换将 ANF 系数转换为真值表
	length := 1 << n
	truthTable := make([]byte, length)
	copy(truthTable, coefficients)
	
	// 调用逆变换
	fmtInverseInplace(truthTable)

	// 创建布尔函数
	return NewFromTruthTable(truthTable)
}

// parseANF 解析 ANF 字符串并返回系数向量
func parseANF(anfString string, n int) ([]byte, error) {
	length := 1 << n
	coefficients := make([]byte, length)

	// 清理字符串：移除空格，统一小写
	cleanANF := strings.ReplaceAll(strings.ToLower(anfString), " ", "")
	
	// 处理空字符串或 "0"
	if cleanANF == "" || cleanANF == "0" {
		return coefficients, nil // 全零函数
	}

	// 分割项 (用 + 分割)
	terms := strings.Split(cleanANF, "+")
	
	for _, term := range terms {
		term = strings.TrimSpace(term)
		if term == "" {
			continue
		}
		
		// 解析单个项
		index, err := parseTerm(term, n)
		if err != nil {
			return nil, fmt.Errorf("failed to parse term '%s': %v", term, err)
		}
		
		// 设置系数 (在 GF(2) 中，+ 就是 XOR)
		coefficients[index] ^= 1
	}

	return coefficients, nil
}

// parseTerm 解析单个 ANF 项，返回对应的索引
func parseTerm(term string, n int) (int, error) {
	// 处理常数项
	if term == "1" {
		return 0, nil
	}
	
	// 处理变量项
	index := 0
	
	// 分割变量 (用 * 分割)
	vars := strings.Split(term, "*")
	
	for _, variable := range vars {
		variable = strings.TrimSpace(variable)
		if variable == "" {
			continue
		}
		
		// 解析变量格式 x0, x1, x2, ...
		if !strings.HasPrefix(variable, "x") {
			return 0, fmt.Errorf("invalid variable format: %s", variable)
		}
		
		varNumStr := variable[1:]
		varNum, err := strconv.Atoi(varNumStr)
		if err != nil {
			return 0, fmt.Errorf("invalid variable number: %s", varNumStr)
		}
		
		if varNum < 0 || varNum >= n {
			return 0, fmt.Errorf("variable x%d out of range (n=%d)", varNum, n)
		}
		
		// 设置对应的位
		index |= 1 << varNum
	}
	
	return index, nil
}

// --- 基础方法 ---

// N 返回函数的变量个数 n.
func (f *BooleanFunction) N() int { // 方法绑定结构体Boolean Function
	return f.n
}

// TruthTable 返回函数的真值表副本。
// 修改返回的切片不会影响原函数。、

// 这个返回类型是 []byte，在go中json序列化会把他当成Base64 编码的二进制数据，然后将这个二进制数据用base64输出
//func (f *BooleanFunction) TruthTable() []byte {
//	return append([]byte(nil), f.truthTable...)
//}

//func (f *BooleanFunction) TruthTable() []int {
//	tt := make([]int, len(f.truthTable))
//	for i, v := range f.truthTable {
//		tt[i] = int(v)
//	}
//	return tt
//}

// TruthTable 返回解包后的[]byte 真值表
func (f *BooleanFunction) TruthTable() []byte {
	return truthTableFromUint64Slice(f.packedTruthTable, 1<<f.n)
}

// 辅助函数
// uint64SliceFromTruthTable 将 truth table (0/1 bytes) 打包到 uint64 切片中
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

// truthTableFromUint64Slice 将打包后的 bitset 解包成 truth table
func truthTableFromUint64Slice(bitslice []uint64, length int) []byte {
	tt := make([]byte, length)
	for i := 0; i < length; i++ {
		if (bitslice[i>>6]>>uint(i&63))&1 == 1 {
			tt[i] = 1
		}
	}
	return tt
}
