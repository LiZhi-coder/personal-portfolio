package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hui-cyber/BoolCore/backend/internal/api"
)

// 设置测试环境
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	api.RegisterRoutes(router)
	return router
}

// 测试请求结构
type TestRequest struct {
	Type          string `json:"type"`
	N             int    `json:"n,omitempty"`
	TruthTable    []byte `json:"truthTable,omitempty"`
	HexValue      string `json:"hexValue,omitempty"`
	IntValue      uint64 `json:"intValue,omitempty"`
	ANFExpression string `json:"anfExpression,omitempty"`
}

// 测试响应结构
type TestResponse struct {
	N                      int    `json:"n"`
	TruthTable             []int  `json:"truthTable"`
	HammingWeight          int    `json:"hammingWeight"`
	IsBalanced             bool   `json:"isBalanced"`
	ANF                    string `json:"anf"`
	AlgebraicDegree        int    `json:"algebraicDegree"`
	Nonlinearity           int64  `json:"nonlinearity"`
	IsBent                 bool   `json:"isBent"`
	CorrelationImmunity    int    `json:"correlationImmunity"`
	ResiliencyOrder        int    `json:"resiliencyOrder"`
	AlgebraicImmunity      int    `json:"algebraicImmunity"`
	DifferentialUniformity int64  `json:"differentialUniformity"`
	Error                  string `json:"error,omitempty"`
}

// 执行API测试的辅助函数
func performAPITest(t *testing.T, router *gin.Engine, req TestRequest) *TestResponse {
	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/analyze", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	var response TestResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("响应JSON解析失败: %v", err)
	}

	return &response
}

// TestPingEndpoint 测试ping接口
func TestPingEndpoint(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ping", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 200, 实际得到 %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] == "" {
		t.Error("响应中缺少message字段")
	}
}

// TestTruthTableInput 测试真值表输入
func TestTruthTableInput(t *testing.T) {
	router := setupRouter()

	testCases := []struct {
		name             string
		request          TestRequest
		expectedHW       int   // 汉明重量
		expectedBalanced bool  // 是否平衡
		expectedDegree   int   // 代数次数
		expectedNL       int64 // 非线性度
	}{
		{
			name: "3元函数-基础测试",
			request: TestRequest{
				Type:       "truthTable",
				N:          3,
				TruthTable: []byte{0, 1, 0, 1, 0, 1, 1, 0},
			},
			expectedHW:       4,
			expectedBalanced: true,
			expectedDegree:   2,
			expectedNL:       2,
		},
		{
			name: "1元函数-边界测试",
			request: TestRequest{
				Type:       "truthTable",
				N:          1,
				TruthTable: []byte{0, 1},
			},
			expectedHW:       1,
			expectedBalanced: true,
			expectedDegree:   1,
			expectedNL:       0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := performAPITest(t, router, tc.request)

			if response.Error != "" {
				t.Errorf("意外错误: %s", response.Error)
				return
			}

			if response.HammingWeight != tc.expectedHW {
				t.Errorf("汉明重量不匹配: 期望 %d, 实际 %d", tc.expectedHW, response.HammingWeight)
			}

			if response.IsBalanced != tc.expectedBalanced {
				t.Errorf("平衡性不匹配: 期望 %t, 实际 %t", tc.expectedBalanced, response.IsBalanced)
			}

			if response.AlgebraicDegree != tc.expectedDegree {
				t.Errorf("代数次数不匹配: 期望 %d, 实际 %d", tc.expectedDegree, response.AlgebraicDegree)
			}

			if response.Nonlinearity != tc.expectedNL {
				t.Errorf("非线性度不匹配: 期望 %d, 实际 %d", tc.expectedNL, response.Nonlinearity)
			}
		})
	}
}

// TestHexInput 测试十六进制输入
func TestHexInput(t *testing.T) {
	router := setupRouter()

	testCases := []struct {
		name       string
		hexValue   string
		n          int
		expectedHW int
	}{
		{"十六进制96", "96", 3, 4},
		{"十六进制FFFF", "FFFF", 4, 16},
		{"十六进制0", "0", 3, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := TestRequest{
				Type:     "hex",
				N:        tc.n,
				HexValue: tc.hexValue,
			}

			response := performAPITest(t, router, request)

			if response.Error != "" {
				t.Errorf("意外错误: %s", response.Error)
				return
			}

			if response.HammingWeight != tc.expectedHW {
				t.Errorf("汉明重量不匹配: 期望 %d, 实际 %d", tc.expectedHW, response.HammingWeight)
			}
		})
	}
}

// TestIntInput 测试整数输入
func TestIntInput(t *testing.T) {
	router := setupRouter()

	testCases := []struct {
		name       string
		intValue   uint64
		n          int
		expectedHW int
	}{
		{"整数150", 150, 3, 4}, // 150 = 0x96
		{"整数255", 255, 3, 8}, // 255 = 0xFF
		{"整数0", 0, 3, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := TestRequest{
				Type:     "int",
				N:        tc.n,
				IntValue: tc.intValue,
			}

			response := performAPITest(t, router, request)

			if response.Error != "" {
				t.Errorf("意外错误: %s", response.Error)
				return
			}

			if response.HammingWeight != tc.expectedHW {
				t.Errorf("汉明重量不匹配: 期望 %d, 实际 %d", tc.expectedHW, response.HammingWeight)
			}
		})
	}
}

// TestANFInput 测试ANF输入
func TestANFInput(t *testing.T) {
	router := setupRouter()

	testCases := []struct {
		name           string
		anfExpression  string
		n              int
		expectedDegree int
	}{
		{"线性函数", "x0 + x1 + x2", 3, 1},
		{"二次函数", "x0 + x1*x2", 3, 2},
		{"常数函数", "1", 3, 0},
		{"零函数", "0", 3, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := TestRequest{
				Type:          "anf",
				N:             tc.n,
				ANFExpression: tc.anfExpression,
			}

			response := performAPITest(t, router, request)

			if response.Error != "" {
				t.Errorf("意外错误: %s", response.Error)
				return
			}

			if response.AlgebraicDegree != tc.expectedDegree {
				t.Errorf("代数次数不匹配: 期望 %d, 实际 %d", tc.expectedDegree, response.AlgebraicDegree)
			}
		})
	}
}

// TestSpecialFunctions 测试特殊函数类型
func TestSpecialFunctions(t *testing.T) {
	router := setupRouter()

	// Bent函数测试 (n=4, hex=6996)
	t.Run("Bent函数测试", func(t *testing.T) {
		request := TestRequest{
			Type:     "hex",
			N:        4,
			HexValue: "6996",
		}

		response := performAPITest(t, router, request)

		if response.Error != "" {
			t.Errorf("意外错误: %s", response.Error)
			return
		}

		if !response.IsBent {
			t.Error("应该识别为Bent函数")
		}

		if response.Nonlinearity != 6 {
			t.Errorf("Bent函数非线性度应为6, 实际 %d", response.Nonlinearity)
		}

		if !response.IsBalanced {
			t.Error("Bent函数应该是平衡的")
		}
	})
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	router := setupRouter()

	errorCases := []struct {
		name    string
		request TestRequest
	}{
		{
			name: "无效类型",
			request: TestRequest{
				Type: "invalid",
			},
		},
		{
			name: "真值表缺少n参数",
			request: TestRequest{
				Type:       "truthTable",
				TruthTable: []byte{0, 1, 0, 1},
			},
		},
		{
			name: "十六进制缺少hexValue",
			request: TestRequest{
				Type: "hex",
				N:    3,
			},
		},
		{
			name: "ANF缺少表达式",
			request: TestRequest{
				Type: "anf",
				N:    3,
			},
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tc.request)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/analyze", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("期望状态码 400, 实际得到 %d", w.Code)
			}

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)

			if response["error"] == "" {
				t.Error("错误响应中应包含error字段")
			}
		})
	}
}

// TestInputFormatConsistency 测试不同输入格式的一致性
func TestInputFormatConsistency(t *testing.T) {
	router := setupRouter()

	// 测试相同函数的不同表示方式应产生相同结果
	// 例: 真值表[0,1,1,0,1,0,0,1] = hex 96 = int 150

	truthTableReq := TestRequest{
		Type:       "truthTable",
		N:          3,
		TruthTable: []byte{0, 1, 1, 0, 1, 0, 0, 1},
	}

	hexReq := TestRequest{
		Type:     "hex",
		N:        3,
		HexValue: "96",
	}

	intReq := TestRequest{
		Type:     "int",
		N:        3,
		IntValue: 150,
	}

	ttResp := performAPITest(t, router, truthTableReq)
	hexResp := performAPITest(t, router, hexReq)
	intResp := performAPITest(t, router, intReq)

	// 检查关键属性是否一致
	if ttResp.HammingWeight != hexResp.HammingWeight || hexResp.HammingWeight != intResp.HammingWeight {
		t.Errorf("汉明重量不一致: TT=%d, Hex=%d, Int=%d",
			ttResp.HammingWeight, hexResp.HammingWeight, intResp.HammingWeight)
	}

	if ttResp.IsBalanced != hexResp.IsBalanced || hexResp.IsBalanced != intResp.IsBalanced {
		t.Errorf("平衡性不一致: TT=%t, Hex=%t, Int=%t",
			ttResp.IsBalanced, hexResp.IsBalanced, intResp.IsBalanced)
	}

	if ttResp.Nonlinearity != hexResp.Nonlinearity || hexResp.Nonlinearity != intResp.Nonlinearity {
		t.Errorf("非线性度不一致: TT=%d, Hex=%d, Int=%d",
			ttResp.Nonlinearity, hexResp.Nonlinearity, intResp.Nonlinearity)
	}
}

// BenchmarkAnalyzeFunction 性能基准测试
func BenchmarkAnalyzeFunction(b *testing.B) {
	router := setupRouter()

	request := TestRequest{
		Type:     "int",
		N:        6,
		IntValue: 123456,
	}

	jsonData, _ := json.Marshal(request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/analyze", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}
}
