#!/bin/bash

# BoolCore 性能测试脚本
# 测试响应时间和并发性能

echo "=== BoolCore 性能测试开始 ==="
echo "时间: $(date)"
echo

BASE_URL="http://localhost:8080"

# 检查服务状态
if ! curl -s "$BASE_URL/api/ping" > /dev/null; then
    echo "❌ 后端服务未启动"
    exit 1
fi
echo "✅ 后端服务正常"
echo

# 响应时间测试函数
test_response_time() {
    local test_name="$1"
    local payload="$2"
    local iterations=5
    
    echo "--- $test_name ---"
    echo "测试 $iterations 次请求的响应时间..."
    
    local total_time=0
    local min_time=999999
    local max_time=0
    
    for i in $(seq 1 $iterations); do
        start_time=$(date +%s%N)
        response=$(curl -s -X POST "$BASE_URL/api/analyze" \
            -H "Content-Type: application/json" \
            -d "$payload")
        end_time=$(date +%s%N)
        
        response_time=$((($end_time - $start_time) / 1000000))
        total_time=$(($total_time + $response_time))
        
        if [ $response_time -lt $min_time ]; then
            min_time=$response_time
        fi
        if [ $response_time -gt $max_time ]; then
            max_time=$response_time
        fi
        
        echo "  第${i}次: ${response_time}ms"
    done
    
    avg_time=$(($total_time / $iterations))
    
    echo "结果统计:"
    echo "  平均响应时间: ${avg_time}ms"
    echo "  最小响应时间: ${min_time}ms"  
    echo "  最大响应时间: ${max_time}ms"
    echo
}

# 内存使用测试函数
test_memory_usage() {
    local test_name="$1"
    local payload="$2"
    
    echo "--- $test_name 内存使用测试 ---"
    
    # 获取Go进程PID
    go_pid=$(pgrep -f "go run.*main.go" | head -1)
    if [ -z "$go_pid" ]; then
        echo "❌ 找不到Go进程"
        return
    fi
    
    # 测试前内存使用
    mem_before=$(ps -p $go_pid -o rss= | tr -d ' ')
    echo "测试前内存使用: ${mem_before}KB"
    
    # 发送请求
    curl -s -X POST "$BASE_URL/api/analyze" \
        -H "Content-Type: application/json" \
        -d "$payload" > /dev/null
    
    sleep 1
    
    # 测试后内存使用
    mem_after=$(ps -p $go_pid -o rss= | tr -d ' ')
    echo "测试后内存使用: ${mem_after}KB"
    echo "内存增长: $((mem_after - mem_before))KB"
    echo
}

# 1. 响应时间测试
echo "=== 响应时间测试 ==="

test_response_time "小规模函数(n=3)" \
    '{"type":"int","n":3,"intValue":150}'

test_response_time "中等规模函数(n=5)" \
    '{"type":"int","n":5,"intValue":12345}'

test_response_time "较大规模函数(n=6)" \
    '{"type":"int","n":6,"intValue":123456}'

# 2. 内存使用测试  
echo "=== 内存使用测试 ==="

test_memory_usage "n=3函数" \
    '{"type":"int","n":3,"intValue":150}'

test_memory_usage "n=6函数" \
    '{"type":"int","n":6,"intValue":123456}'

# 3. 并发测试(如果有wrk)
echo "=== 并发性能测试 ==="

if command -v wrk &> /dev/null; then
    echo "使用wrk进行并发测试..."
    
    # 创建POST脚本
    cat > /tmp/post_test.lua << 'EOF'
wrk.method = "POST"
wrk.body = '{"type":"int","n":4,"intValue":12345}'
wrk.headers["Content-Type"] = "application/json"
EOF
    
    echo "测试配置: 4线程, 20并发连接, 持续10秒"
    wrk -t4 -c20 -d10s -s /tmp/post_test.lua "$BASE_URL/api/analyze"
    
    rm -f /tmp/post_test.lua
    
elif command -v ab &> /dev/null; then
    echo "使用Apache Bench进行并发测试..."
    
    # 创建测试数据
    echo '{"type":"int","n":4,"intValue":12345}' > /tmp/test_payload.json
    
    echo "测试配置: 500总请求, 10并发"
    ab -n 500 -c 10 -p /tmp/test_payload.json -T application/json \
       "$BASE_URL/api/analyze"
    
    rm -f /tmp/test_payload.json
    
else
    echo "⚠️  未找到wrk或ab工具，跳过并发测试"
    echo "安装方法:"
    echo "  Ubuntu: sudo apt install wrk apache2-utils"
    echo "  macOS: brew install wrk"
fi

echo
echo "=== 性能测试完成 ==="
echo "请将测试结果记录到测试文档.md的性能测试结果表中"

# 4. 纯后端函数性能（可选）
echo
echo "=== 纯后端函数性能（无HTTP开销） ==="
if [ -f "../backend/cmd/perf/main.go" ] || [ -f "backend/cmd/perf/main.go" ]; then
    echo "运行 n=6, int=123456 的核心分析..."
    (cd ../backend 2>/dev/null && go run ./cmd/perf -type int -n 6 -int 123456 -repeat 1 -format text) \
    || (cd backend 2>/dev/null && go run ./cmd/perf -type int -n 6 -int 123456 -repeat 1 -format text)
else
    echo "跳过：未找到 backend/cmd/perf/main.go"
fi