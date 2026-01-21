#!/bin/bash

# BoolCore 功能测试脚本
# 运行所有功能测试用例并记录结果

echo "=== BoolCore 功能测试开始 ==="
echo "时间: $(date)"
echo

BASE_URL="http://localhost:8080"

# 检查服务是否启动
echo "检查服务状态..."
if ! curl -s "$BASE_URL/api/ping" > /dev/null; then
    echo "❌ 后端服务未启动，请先运行: cd backend && go run cmd/server/main.go"
    exit 1
fi
echo "✅ 后端服务正常"
echo

# 测试用例函数
run_test() {
    local test_name="$1"
    local payload="$2"
    local expected_desc="$3"
    
    echo "--- $test_name ---"
    echo "请求: $payload"
    echo "预期: $expected_desc"
    
    start_time=$(date +%s%N)
    response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/analyze" \
        -H "Content-Type: application/json" \
        -d "$payload")
    end_time=$(date +%s%N)
    
    http_code=$(echo "$response" | tail -n1)
    json_response=$(echo "$response" | head -n -1)
    
    response_time=$((($end_time - $start_time) / 1000000))
    
    echo "HTTP状态码: $http_code"
    echo "响应时间: ${response_time}ms"
    
    if [ "$http_code" -eq 200 ]; then
        echo "✅ 请求成功"
        echo "关键结果:"
        echo "$json_response" | jq -r '
            "  变量数(n): " + (.n // "N/A" | tostring) +
            "\n  汉明重量: " + (.hammingWeight // "N/A" | tostring) +
            "\n  是否平衡: " + (.isBalanced // "N/A" | tostring) +
            "\n  ANF: " + (.anf // "N/A") +
            "\n  代数次数: " + (.algebraicDegree // "N/A" | tostring) +
            "\n  非线性度: " + (.nonlinearity // "N/A" | tostring) +
            "\n  是否Bent: " + (.isBent // "N/A" | tostring)
        '
    else
        echo "❌ 请求失败"
        echo "错误信息: $json_response"
    fi
    echo
}

# FT-002: 真值表输入
run_test "FT-002: 真值表输入" \
    '{"type":"truthTable","n":3,"truthTable":[0,1,0,1,0,1,1,0]}' \
    "汉明重量=4, 平衡=true, ANF=x0+x1*x2, 非线性度=2"

# FT-003: 十六进制输入
run_test "FT-003: 十六进制输入" \
    '{"type":"hex","n":3,"hexValue":"96"}' \
    "与真值表[0,1,1,0,1,0,0,1]等价"

# FT-004: 整数输入
run_test "FT-004: 整数输入" \
    '{"type":"int","n":3,"intValue":150}' \
    "整数150 = 0x96, 应与FT-003结果相同"

# FT-005: ANF输入
run_test "FT-005: ANF输入" \
    '{"type":"anf","n":3,"anfExpression":"x0 + x1*x2"}' \
    "直接使用ANF表达式"

# FT-006: 最小函数
run_test "FT-006: 最小函数(n=1)" \
    '{"type":"truthTable","n":1,"truthTable":[0,1]}' \
    "汉明重量=1, 平衡=true, ANF=x0, 非线性度=0"

# FT-009: Bent函数
run_test "FT-009: Bent函数测试" \
    '{"type":"hex","n":4,"hexValue":"6996"}' \
    "汉明重量=8, 平衡=true, 非线性度=6, Bent=true"

# FT-010: 线性函数
run_test "FT-010: 线性函数" \
    '{"type":"anf","n":3,"anfExpression":"x0 + x1 + x2"}' \
    "代数次数=1, 非线性度=0"

# FT-011: 常数函数
run_test "FT-011: 常数函数" \
    '{"type":"int","n":3,"intValue":255}' \
    "汉明重量=8, 平衡=false, ANF=1, 代数次数=0"

echo "=== 边界条件测试 ==="

# 错误输入测试
echo "--- 错误输入测试 ---"

echo "测试无效type:"
curl -s -w "\nHTTP: %{http_code}\n" -X POST "$BASE_URL/api/analyze" \
    -H "Content-Type: application/json" \
    -d '{"type":"invalid"}' | jq -r '.error // .'
echo

echo "测试真值表长度错误:"
curl -s -w "\nHTTP: %{http_code}\n" -X POST "$BASE_URL/api/analyze" \
    -H "Content-Type: application/json" \
    -d '{"type":"truthTable","n":3,"truthTable":[0,1,0]}' | jq -r '.error // .'
echo

echo "测试缺少必需参数:"
curl -s -w "\nHTTP: %{http_code}\n" -X POST "$BASE_URL/api/analyze" \
    -H "Content-Type: application/json" \
    -d '{"type":"hex","n":3}' | jq -r '.error // .'
echo

echo "=== 功能测试完成 ==="
echo "请查看上述结果，将实际数据填入测试文档.md中的表格"