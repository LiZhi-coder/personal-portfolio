#!/bin/bash

# BoolCore 快速测试脚本
# 运行所有基本测试并生成简要报告

echo "🚀 BoolCore 快速测试"
echo "===================="
echo

BASE_URL="http://localhost:8080"

# 1. 服务连通性测试
echo "1️⃣ 服务连通性测试"
if response=$(curl -s "$BASE_URL/api/ping" 2>/dev/null); then
    echo "✅ 后端服务正常: $response"
else
    echo "❌ 后端服务异常，请检查服务是否启动"
    echo "启动命令: cd backend && go run cmd/server/main.go"
    exit 1
fi
echo

# 2. 快速功能测试
echo "2️⃣ 快速功能测试"

test_cases=(
    "真值表输入|{\"type\":\"truthTable\",\"n\":3,\"truthTable\":[0,1,0,1,0,1,1,0]}"
    "十六进制输入|{\"type\":\"hex\",\"n\":3,\"hexValue\":\"96\"}"
    "整数输入|{\"type\":\"int\",\"n\":3,\"intValue\":150}"
    "ANF输入|{\"type\":\"anf\",\"n\":3,\"anfExpression\":\"x0 + x1*x2\"}"
    "Bent函数|{\"type\":\"hex\",\"n\":4,\"hexValue\":\"6996\"}"
)

for test_case in "${test_cases[@]}"; do
    IFS='|' read -r name payload <<< "$test_case"
    
    start_time=$(date +%s%N)
    response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/api/analyze" \
        -H "Content-Type: application/json" \
        -d "$payload")
    end_time=$(date +%s%N)
    
    http_code="${response: -3}"
    json_data="${response%???}"
    response_time=$((($end_time - $start_time) / 1000000))
    
    if [ "$http_code" = "200" ]; then
        hamming=$(echo "$json_data" | jq -r '.hammingWeight // "N/A"')
        balanced=$(echo "$json_data" | jq -r '.isBalanced // "N/A"')
        nonlinearity=$(echo "$json_data" | jq -r '.nonlinearity // "N/A"')
        
        echo "✅ $name (${response_time}ms) - 汉明重量:$hamming, 平衡:$balanced, 非线性度:$nonlinearity"
    else
        error_msg=$(echo "$json_data" | jq -r '.error // "未知错误"')
        echo "❌ $name - HTTP:$http_code, 错误:$error_msg"
    fi
done
echo

# 3. 错误处理测试
echo "3️⃣ 错误处理测试"

error_tests=(
    "无效类型|{\"type\":\"invalid\"}"
    "缺少参数|{\"type\":\"hex\",\"n\":3}"
    "真值表长度错误|{\"type\":\"truthTable\",\"n\":3,\"truthTable\":[0,1]}"
)

for test_case in "${error_tests[@]}"; do
    IFS='|' read -r name payload <<< "$test_case"
    
    response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/api/analyze" \
        -H "Content-Type: application/json" \
        -d "$payload")
    
    http_code="${response: -3}"
    json_data="${response%???}"
    
    if [ "$http_code" = "400" ]; then
        error_msg=$(echo "$json_data" | jq -r '.error // "N/A"')
        echo "✅ $name - 正确返回400错误: $error_msg"
    else
        echo "❌ $name - 错误处理异常 HTTP:$http_code"
    fi
done
echo

# 4. 性能概览
echo "4️⃣ 性能概览测试"

perf_tests=(
    "n=3|{\"type\":\"int\",\"n\":3,\"intValue\":123}"
    "n=5|{\"type\":\"int\",\"n\":5,\"intValue\":12345}"
    "n=6|{\"type\":\"int\",\"n\":6,\"intValue\":123456}"
)

for test_case in "${perf_tests[@]}"; do
    IFS='|' read -r name payload <<< "$test_case"
    
    # 测试3次取平均值
    total_time=0
    for i in {1..3}; do
        start_time=$(date +%s%N)
        curl -s -X POST "$BASE_URL/api/analyze" \
            -H "Content-Type: application/json" \
            -d "$payload" > /dev/null
        end_time=$(date +%s%N)
        
        response_time=$((($end_time - $start_time) / 1000000))
        total_time=$(($total_time + $response_time))
    done
    
    avg_time=$(($total_time / 3))
    
    if [ $avg_time -lt 100 ]; then
        status="✅ 优秀"
    elif [ $avg_time -lt 500 ]; then
        status="⚠️  一般" 
    else
        status="❌ 较慢"
    fi
    
    echo "$status $name 平均响应时间: ${avg_time}ms"
done
echo

# 5. 总结报告
echo "📊 测试总结"
echo "============"
echo "✅ 如果所有测试都显示成功，系统功能正常"
echo "⚠️  如果有警告，建议查看详细的性能测试"
echo "❌ 如果有失败，请检查代码或环境配置"
echo
echo "📖 详细测试报告请运行:"
echo "   ./test_functional.sh  # 详细功能测试"
echo "   ./test_performance.sh # 详细性能测试"
echo
echo "📋 测试结果记录到: 测试文档.md"