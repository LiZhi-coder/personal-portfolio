package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hui-cyber/BoolCore/backend/pkg/booleancore"
)

type stepTiming struct {
	Name   string  `json:"name"`
	Micros float64 `json:"us"`
}

type perfResult struct {
	N          int          `json:"n"`
	Input      string       `json:"input"`
	Repeat     int          `json:"repeat"`
	TotalUs    float64      `json:"total_us"`
	Steps      []stepTiming `json:"steps"`
	Properties struct {
		HammingWeight          int     `json:"hammingWeight"`
		IsBalanced             bool    `json:"isBalanced"`
		ANFLen                 int     `json:"anfLen"`
		AlgebraicDegree        int     `json:"algebraicDegree"`
		TransparencyOrder      float64 `json:"transparencyOrder"`
		WalshLen               int     `json:"walshLen"`
		AutocorrLen            int     `json:"autocorrelationLen"`
		Nonlinearity           int64   `json:"nonlinearity"`
		CorrelationImmunity    int     `json:"correlationImmunity"`
		ResiliencyOrder        int     `json:"resiliencyOrder"`
		IsBent                 bool    `json:"isBent"`
		SumOfSquareIndicator   int64   `json:"sumOfSquareIndicator"`
		IsRotationSymmetric    bool    `json:"isRotationSymmetric"`
		AbsoluteWalshKinds     int     `json:"absoluteWalshKinds"`
		AbsoluteAutocorrKinds  int     `json:"absoluteAutocorrKinds"`
		AbsoluteIndicator      int64   `json:"absoluteIndicator"`
		DifferentialUniformity int64   `json:"differentialUniformity"`
		AlgebraicImmunity      int     `json:"algebraicImmunity"`
	} `json:"properties"`
}

func main() {
	var (
		inType   string
		n        int
		intValue uint64
		hexValue string
		anf      string
		repeat   int
		format   string
	)

	flag.StringVar(&inType, "type", "int", "input type: int|hex|anf|truth (truth not implemented in CLI)")
	flag.IntVar(&n, "n", 6, "number of variables")
	flag.Uint64Var(&intValue, "int", 123456, "integer value for truth table (low bit = index 0)")
	flag.StringVar(&hexValue, "hex", "", "hex value for truth table")
	flag.StringVar(&anf, "anf", "", "ANF expression, e.g. 'x0 + x1*x2 + 1'")
	flag.IntVar(&repeat, "repeat", 1, "repeat runs to warm cache and average")
	flag.StringVar(&format, "format", "json", "output format: json|text")
	flag.Parse()

	var bf *booleancore.BooleanFunction
	var err error

	inputDesc := ""
	switch strings.ToLower(inType) {
	case "int":
		bf, err = booleancore.NewFromInt(intValue, n)
		inputDesc = fmt.Sprintf("int:%d", intValue)
	case "hex":
		bf, err = booleancore.NewFromHex(hexValue, n)
		inputDesc = fmt.Sprintf("hex:%s", hexValue)
	case "anf":
		bf, err = booleancore.NewFromANF(n, anf)
		inputDesc = fmt.Sprintf("anf:%s", anf)
	default:
		fmt.Fprintf(os.Stderr, "unsupported type: %s\n", inType)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to construct BooleanFunction: %v\n", err)
		os.Exit(2)
	}

	// 每次重复都重新构造，确保缓存影响与真实使用一致（首次昂贵，后续受益于缓存）
	var last perfResult
	totalStart := time.Now()
	for r := 0; r < repeat; r++ {
		// 构造新的实例
		switch strings.ToLower(inType) {
		case "int":
			bf, err = booleancore.NewFromInt(intValue, n)
		case "hex":
			bf, err = booleancore.NewFromHex(hexValue, n)
		case "anf":
			bf, err = booleancore.NewFromANF(n, anf)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "construct run %d error: %v\n", r, err)
			os.Exit(2)
		}

		res := perfResult{N: n, Input: inputDesc, Repeat: repeat}
		_, timings, _ := booleancore.AnalyzeAllTimed(bf)

		// 重新计算一次用于取属性输出（避免复制大切片），这里用 AnalyzeAll 更简洁
		props := booleancore.AnalyzeAll(bf)

		// 基本属性
		res.Properties.HammingWeight = props.HammingWeight
		res.Properties.IsBalanced = props.IsBalanced
		res.Properties.ANFLen = len(props.ANF)
		res.Properties.AlgebraicDegree = props.AlgebraicDegree
		res.Properties.TransparencyOrder = props.TransparencyOrder
		res.Properties.Nonlinearity = props.Nonlinearity
		res.Properties.CorrelationImmunity = props.CorrelationImmunity
		res.Properties.ResiliencyOrder = props.ResiliencyOrder
		res.Properties.IsBent = props.IsBent
		res.Properties.SumOfSquareIndicator = props.SumOfSquareIndicator
		res.Properties.IsRotationSymmetric = props.IsRotationSymmetric
		res.Properties.AbsoluteWalshKinds = len(props.AbsoluteWalshSpectrum)
		res.Properties.AbsoluteAutocorrKinds = len(props.AbsoluteAutocorrelationSpectrum)
		res.Properties.AbsoluteIndicator = props.AbsoluteIndicator
		res.Properties.DifferentialUniformity = props.DifferentialUniformity
		res.Properties.AlgebraicImmunity = props.AlgebraicImmunity

		// 长度信息
		res.Properties.WalshLen = len(props.WalshSpectrum)
		res.Properties.AutocorrLen = len(props.AutocorrelationSpectrum)

		// 步骤时间
		for name, dur := range timings {
			// 使用纳秒转微秒(浮点)以保留子微秒级的小数部分
			res.Steps = append(res.Steps, stepTiming{Name: name, Micros: float64(dur.Nanoseconds()) / 1000.0})
		}

		last = res
	}
	// 总耗时同样用纳秒换算为微秒(浮点)，避免小于1微秒时显示为0
	totalUs := float64(time.Since(totalStart).Nanoseconds()) / 1000.0
	last.TotalUs = totalUs

	switch strings.ToLower(format) {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(last)
	default:
		fmt.Printf("n=%d input=%s repeat=%d total=%.3fμs\n", last.N, last.Input, last.Repeat, last.TotalUs)
		for _, s := range last.Steps {
			fmt.Printf("  %-28s %8.3f μs\n", s.Name+":", s.Micros)
		}
	}
}
