package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hui-cyber/BoolCore/backend/pkg/booleancore" // 模块名加目录
)

// AnalyzeRequest 定义了前端请求的 JSON 结构.
type AnalyzeRequest struct {
	// `json:"type"`
	//→ 告诉 JSON 库：序列化/反序列化时，这个字段对应 JSON 中的 "type" 键
	//`binding:"required"`
	//→ 告诉 Gin 框架：这个字段是必填的，如果前端没传或为空，就报错
	Type          string `json:"type" binding:"required"` // 用于根据不同的type 进行不同的处理
	N             int    `json:"n" `                      // 对于十六进制和整数输入是必须的
	TruthTable    []byte `json:"truthTable" `             //<--- 关键改动：移除了 "required" 标签 因为用户可以进行其他输入不一定都是真值表
	HexValue      string `json:"hexValue"`
	IntValue      uint64 `json:"intValue"`
	ANFExpression string `json:"anfExpression"` // ANF 代数正规式表达式
	// TODO: 或者其他的输入方式
}

// AnalyzeResponse 定义了返回给前端的 JSON 结构.
type AnalyzeResponse struct {
	// “当转 JSON 时，这个字段用 algebraicDegree 作为键名，
	// 这是因为加上了omitempty，如果是空，则不返回这个字段”
	N                               int           `json:"n"`                               // n元布尔函数
	TruthTable                      []int         `json:"truthTable"`                      // 将传入的[]byte修改为[]int，避免json转换为base64
	HammingWeight                   int           `json:"hammingWeight"`                   // 汉明重量
	IsBalanced                      bool          `json:"isBalanced"`                      // 平衡
	WalshSpectrum                   []int64       `json:"walshSpectrum"`                   // 输出Walsh谱
	ANF                             string        `json:"anf,omitempty"`                   // 代数标准型
	AlgebraicDegree                 int           `json:"algebraicDegree,omitempty"`       // 代数次数
	Nonlinearity                    int64         `json:"nonlinearity,omitempty"`          // 非线性度
	AutocorrelationSpectrum         []int64       `json:"autocorrelationSpectrum"`         // 自相关谱
	CorrelationImmunity             int           `json:"correlationImmunity"`             // 相关免疫度
	ResiliencyOrder                 int           `json:"resiliencyOrder"`                 // 弹性阶数
	TransparencyOrder               float64       `json:"transparencyOrder"`               // 透明度阶
	IsBent                          bool          `json:"isBent"`                          // 是否 bent
	SumOfSquareIndicator            int64         `json:"sumOfSquareIndicator"`            // 平方和指标
	IsRotationSymmetric             bool          `json:"isRotationSymmetric"`             // 是否旋转对称
	AbsoluteWalshSpectrum           map[int64]int `json:"absoluteWalshSpectrum"`           // 绝对walsh谱分布
	AbsoluteAutocorrelationSpectrum map[int64]int `json:"absoluteAutocorrelationSpectrum"` // 绝对自相关谱分布
	AbsoluteIndicator               int64         `json:"absoluteIndicator"`               // 绝对指标
	DifferentialUniformity          int64         `json:"differentialUniformity"`          // 差分均匀度
	AlgebraicImmunity               int           `json:"algebraicImmunity"`               // 代数免疫度
	Annihilator                     string        `json:"annihilator,omitempty"`           // 零化因子ANF表达式
	// TODO: 添加更多字段
}

// AnalyzeFunctionHandler 是处理性质分析请求的 Gin Handler.  是/api/analyze的处理函数
func AnalyzeFunctionHandler(c *gin.Context) {
	var req AnalyzeRequest // 定义请求结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bf *booleancore.BooleanFunction
	var err error

	// 1. 使用核心库创建一个布尔函数实例，由于可以输入不同的类型，输出真值表啊，十六进制啊，整数啊，所以这里根据不同的type 进行不同的处理
	switch req.Type {
	case "truthTable":
		if req.TruthTable == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'truthTable' is required for type 'truthTable'"})
			return
		}
		bf, err = booleancore.NewFromTruthTable(req.TruthTable)
	case "hex":
		if req.N == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'n' is required for type 'hex'"})
			return
		}
		if req.HexValue == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'hexValue' is required for type 'hex'"})
			return
		}
		bf, err = booleancore.NewFromHex(req.HexValue, req.N)
	case "int":
		if req.N == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'n' is required for type 'int'"})
			return
		}
		// 对于int整数类型，intValue 为 0 是一个有效值, 所以我们只检查 n
		bf, err = booleancore.NewFromInt(req.IntValue, req.N)
	case "anf":
		if req.N == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'n' is required for type 'anf'"})
			return
		}
		if req.ANFExpression == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'anfExpression' is required for type 'anf'"})
			return
		}
		bf, err = booleancore.NewFromANF(req.N, req.ANFExpression)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'type' specified, must be one of [truthTable, hex, int, anf]"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 调用核心库进行性性质分析
	// 【解决 Base64 问题】将 []byte 转换为 []int
	originalTT := bf.TruthTable()
	ttAsInt := make([]int, len(originalTT))
	for i, v := range originalTT {
		ttAsInt[i] = int(v)
	}

	// 3. 调用计算分析函数返回结果
	// 优化：如果输入是ANF，直接使用原始ANF，避免重复计算
	var anfString string
	var algebraicDegree int

	if req.Type == "anf" {
		anfString = req.ANFExpression // 直接使用输入的ANF
		// 直接从输入的ANF字符串计算代数次数，避免重复计算ANF系数
		algebraicDegree = calculateDegreeFromANF(req.ANFExpression)
	} else {
		anfString = bf.AlgebraicNormalForm()   // 其他输入类型需要计算ANF
		algebraicDegree = bf.AlgebraicDegree() // 正常计算代数次数
	}

	// 计算代数免疫度（快速版本，不计算零化因子）
	var algebraicImmunity int
	ai, _, err := bf.AlgebraicImmunity(false) // false = 只计算代数免疫度值，不求解零化因子
	if err != nil {
		algebraicImmunity = -1
	} else {
		algebraicImmunity = ai
	}

	// 【可选功能 - 已注释】计算零化因子表达式（比较耗时）
	// 如果需要零化因子，取消下面的注释并注释掉上面的快速版本
	/*
		var algebraicImmunity int
		var annihilator string
		findAnnihilator := bf.N() <= 8 // 可根据需要调整阈值
		ai, ann, err := bf.AlgebraicImmunity(findAnnihilator)
		if err != nil {
			algebraicImmunity = -1
			annihilator = "calculation error: " + err.Error()
		} else {
			algebraicImmunity = ai
			annihilator = ann
		}
	*/

	resp := AnalyzeResponse{
		N:          bf.N(),
		TruthTable: ttAsInt,

		HammingWeight:   bf.HammingWeight(),
		IsBalanced:      bf.IsBalanced(),
		ANF:             anfString,       // 这进行个预处理，如果输入就是ANF就直接用
		AlgebraicDegree: algebraicDegree, // 直接从ANF字符串计算次数

		AutocorrelationSpectrum:         bf.Autocorrelation(),
		WalshSpectrum:                   bf.WalshHadamardTransform(),
		AbsoluteWalshSpectrum:           bf.AbsoluteWalshSpectrum(),
		AbsoluteAutocorrelationSpectrum: bf.AbsoluteAutocorrelation(),

		Nonlinearity:           bf.Nonlinearity(),
		CorrelationImmunity:    bf.CorrelationImmunity(),
		TransparencyOrder:      bf.TransparencyOrder(),
		ResiliencyOrder:        bf.ResiliencyOrder(),
		IsBent:                 bf.IsBent(),
		SumOfSquareIndicator:   bf.SumOfSquareIndicator(),
		IsRotationSymmetric:    bf.IsRotationSymmetric(),
		AbsoluteIndicator:      bf.AbsoluteIndicator(),
		DifferentialUniformity: bf.DifferentialUniformity(),
		AlgebraicImmunity:      algebraicImmunity, // 使用预计算的值
		// Annihilator:                  annihilator,       // 【已禁用】如需启用，取消注释并启用上面的完整计算版本
	}

	c.JSON(http.StatusOK, resp)
}

// calculateDegreeFromANF 直接从ANF字符串计算代数次数，避免重新解析
func calculateDegreeFromANF(anfString string) int {
	// 清理字符串
	cleanANF := strings.ReplaceAll(strings.ToLower(anfString), " ", "")

	// 处理空字符串或 "0"
	if cleanANF == "" || cleanANF == "0" {
		return 0
	}

	// 分割项
	terms := strings.Split(cleanANF, "+")
	maxDegree := 0

	for _, term := range terms {
		term = strings.TrimSpace(term)
		if term == "" {
			continue
		}

		// 计算当前项的次数
		degree := 0
		if term == "1" {
			degree = 0 // 常数项次数为0
		} else {
			// 计算变量个数（用*分割）
			vars := strings.Split(term, "*")
			for _, variable := range vars {
				variable = strings.TrimSpace(variable)
				if variable != "" && strings.HasPrefix(variable, "x") {
					degree++
				}
			}
		}

		if degree > maxDegree {
			maxDegree = degree
		}
	}

	return maxDegree
}

// PingHandler 是 /api/ping 的处理器
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong from BoolCore backend单独!"}) //这个StatusOK 是一个常量，表示 HTTP 状态码 200
}
