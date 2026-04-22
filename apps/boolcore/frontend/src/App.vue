<script setup>
import { ref, computed } from 'vue'
import axios from 'axios'

const message = ref('Click the button to ping the BoolCore backend...')

// 输入控制状态
const inputType = ref('truthTable')
const n = ref(3)
const loading = ref(false)
const analysisResult = ref(null)
const activeTab = ref('basic')

// 输入数据
const truthTableInput = ref('')
const hexInput = ref('')
const intInput = ref(0)
const anfInput = ref('')

// 标签页配置
const tabs = [
  { id: 'basic', label: '基础性质', icon: '🔢' },
  { id: 'spectral', label: '频谱分析', icon: '📊' },
  { id: 'crypto', label: '密码学性质', icon: '🔒' },
  { id: 'export', label: '数据导出', icon: '📁' }
]

// 输入类型选项
const inputTypes = [
  { value: 'truthTable', label: '真值表', needsN: true },
  { value: 'hex', label: '十六进制', needsN: true },
  { value: 'int', label: '整数', needsN: true },
  { value: 'anf', label: 'ANF (代数正规式)', needsN: true }
]

// 计算属性：当前输入类型是否需要n参数
const needsN = computed(() => {
  const currentType = inputTypes.find(type => type.value === inputType.value)
  return currentType?.needsN || false
})

const apiBase = import.meta.env.VITE_BOOLCORE_API_BASE || '/api/boolcore'

const pingBackend = async () => {
  try {
    message.value = 'Pinging...'
    const response = await axios.get(`${apiBase}/ping`)
    message.value = response.data.message
  } catch (error) {
    console.error('Error pinging backend:', error)
    message.value = 'Failed to connect to backend. Is it running?'
  }
}

const analyzeFunction = async () => {
  loading.value = true
  analysisResult.value = null

  try {
    // 构建请求数据
    let requestData = {
      type: inputType.value,
      n: n.value
    }

    // 根据输入类型添加对应数据
    switch (inputType.value) {
      case 'truthTable':
        // 清洗真值表输入，自动识别和清理各种格式
        // 支持: [0,1,0,1], {0,1}, "0101", 0 1 0 1, 0,1,0,1 等
        let cleanInput = truthTableInput.value
          .replace(/[\[\]{}"']/g, '')    // 移除括号和引号
          .replace(/[^01,\s]/g, '')      // 只保留0、1、逗号和空格
        
        let values = []
        if (cleanInput.includes(',')) {
          // 逗号分隔格式: "0,1,0,1" 或 "0, 1, 0, 1"
          values = cleanInput.split(',').map(x => x.trim()).filter(x => x === '0' || x === '1')
        } else if (cleanInput.includes(' ')) {
          // 空格分隔格式: "0 1 0 1"
          values = cleanInput.split(/\s+/).filter(x => x === '0' || x === '1')
        } else {
          // 连续字符串格式: "0101"
          values = cleanInput.split('').filter(x => x === '0' || x === '1')
        }
        
        const ttArray = values.map(x => parseInt(x))
        requestData.truthTable = ttArray
        break
      case 'hex':
        // 自动清理十六进制输入，支持各种格式
        // 支持: "0xAB", "AB CD", "[AB,CD]", "AB, CD", "ABCD" 等
        let cleanHex = hexInput.value
          .replace(/0x/gi, '')           // 移除 0x 前缀
          .replace(/[\[\]{}"'\s,]/g, '') // 移除括号、引号、空格、逗号
          .replace(/[^0-9a-fA-F]/g, '')  // 只保留十六进制字符
        requestData.hexValue = cleanHex
        break
      case 'int':
        requestData.intValue = parseInt(intInput.value)
        break
      case 'anf':
        requestData.anfExpression = anfInput.value
        break
    }

    // 发送到后端分析接口
    const response = await axios.post(`${apiBase}/analyze`, requestData)
    analysisResult.value = response.data

  } catch (error) {
    console.error('Analysis error:', error)
    analysisResult.value = {
      error: '分析失败: ' + (error.response?.data?.error || error.message)
    }
  } finally {
    loading.value = false
  }
}

// 复制到剪贴板功能
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    alert('已复制到剪贴板!')
  } catch (err) {
    console.error('复制失败:', err)
  }
}

// 导出功能
const exportJSON = () => {
  const dataStr = JSON.stringify(analysisResult.value, null, 2)
  copyToClipboard(dataStr)
}

const exportSummary = () => {
  const summary = `布尔函数分析摘要
变量数: ${analysisResult.value.n}
汉明重量: ${analysisResult.value.hammingWeight}
是否平衡: ${analysisResult.value.isBalanced ? '是' : '否'}
代数次数: ${analysisResult.value.algebraicDegree}
非线性度: ${analysisResult.value.nonlinearity}
相关免疫度: ${analysisResult.value.correlationImmunity}
弹性阶数: ${analysisResult.value.resiliencyOrder >= 0 ? analysisResult.value.resiliencyOrder : 'N/A'}
是否Bent函数: ${analysisResult.value.isBent ? '是' : '否'}
ANF表达式: ${analysisResult.value.anf}`
  copyToClipboard(summary)
}

const exportSageMath = () => {
  const sageMathFormat = `# SageMath 格式
R.<${Array.from({length: analysisResult.value.n}, (_, i) => `x${i}`).join(',')}> = BooleanPolynomialRing()
f = ${analysisResult.value.anf}
truth_table = ${JSON.stringify(analysisResult.value.truthTable)}`
  copyToClipboard(sageMathFormat)
}

const exportMatlab = () => {
  const matlabFormat = `% Matlab 格式
n = ${analysisResult.value.n};
truthTable = [${analysisResult.value.truthTable.join(', ')}];
walshSpectrum = [${analysisResult.value.walshSpectrum.join(', ')}];
autocorrelationSpectrum = [${analysisResult.value.autocorrelationSpectrum.join(', ')}];`
  copyToClipboard(matlabFormat)
}

// 文件上传处理
const handleFileUpload = (event) => {
  const file = event.target.files[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (e) => {
    const content = e.target.result.trim()

    // 根据当前输入类型处理文件内容
    switch (inputType.value) {
      case 'truthTable':
        // 处理真值表文件：自动识别和清理各种格式
        // 支持: [0,1,0,1], {0,1}, "0101", 0 1 0 1, 0,1,0,1 等
        const cleanContent = content
          .replace(/[\[\]{}"']/g, '')    // 移除括号和引号
          .replace(/[^01,\s]/g, '')      // 只保留0、1、逗号和空格
        
        let values = []
        if (cleanContent.includes(',')) {
          // 逗号分隔格式: "0,1,0,1" 或 "0, 1, 0, 1"
          values = cleanContent.split(',').map(x => x.trim()).filter(x => x === '0' || x === '1')
        } else if (cleanContent.includes(' ')) {
          // 空格分隔格式: "0 1 0 1"
          values = cleanContent.split(/\s+/).filter(x => x === '0' || x === '1')
        } else {
          // 连续字符串格式: "0101"
          values = cleanContent.split('').filter(x => x === '0' || x === '1')
        }

        truthTableInput.value = values.join(',')

        // 自动推导n值
        const length = values.length
        if (length > 0 && (length & (length - 1)) === 0) { // 检查是否为2的幂
          n.value = Math.log2(length)
        }
        break

      case 'hex':
        // 处理十六进制文件：自动识别和清理各种格式
        // 支持: "0xAB", "AB CD", "[AB,CD]", "AB, CD", "ABCD" 等
        let hexContent = content
          .replace(/0x/gi, '')           // 移除 0x 或 0X 前缀
          .replace(/[\[\]{}"'\s,]/g, '') // 移除括号、引号、空格、逗号
          .replace(/[^0-9a-fA-F]/g, '')  // 只保留十六进制字符

        hexInput.value = hexContent
        break

      case 'int':
        // 处理整数文件：提取数字
        const intMatch = content.match(/\d+/)
        if (intMatch) {
          intInput.value = parseInt(intMatch[0])
        }
        break

      case 'anf':
        // 处理ANF文件：直接使用内容
        anfInput.value = content
        break
    }

    // 清空文件输入，允许重复上传同一文件
    event.target.value = ''
  }

  reader.onerror = () => {
    alert('文件读取失败！')
  }

  reader.readAsText(file)
}
</script>

<template>
  <div class="modern-container">
    <!-- Hero Section -->
    <div class="hero-section">
      <div class="hero-background"></div>
      <div class="hero-content">
        <div class="hero-icon">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
            <defs>
              <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" style="stop-color:#0ea5e9"/>
                <stop offset="50%" style="stop-color:#2563eb"/>
                <stop offset="100%" style="stop-color:#7c3aed"/>
              </linearGradient>
            </defs>
            <path d="M12 2L2 7V17L12 22L22 17V7L12 2Z" stroke="url(#gradient)" stroke-width="2" fill="none"/>
            <circle cx="12" cy="12" r="3" fill="url(#gradient)"/>
          </svg>
        </div>
        <h1 class="hero-title">BoolCore</h1>
        <p class="hero-subtitle">Advanced Boolean Function Cryptographic Analysis Platform</p>
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-number">10+</span>
            <span class="stat-label">分析指标</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">4</span>
            <span class="stat-label">输入格式</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">99.99%</span>
            <span class="stat-label">准确率</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <!-- 后端连接测试 -->
      <div class="modern-card connection-card">
        <div class="card-header">
          <div class="card-icon connection-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path d="M12 2L22 8.5V15.5L12 22L2 15.5V8.5L12 2Z" stroke="currentColor" stroke-width="2" fill="none"/>
              <circle cx="12" cy="12" r="3" fill="currentColor"/>
            </svg>
          </div>
          <h3>Backend Connection</h3>
          <div class="connection-status" :class="message.includes('Success') ? 'connected' : 'disconnected'"></div>
        </div>
        <div class="card-content">
          <button @click="pingBackend" class="modern-button secondary">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path d="M13 2L3 14H12L11 22L21 10H12L13 2Z" fill="currentColor"/>
            </svg>
            Test Connection
          </button>
          <div class="status-display">
            <span class="status-label">Status:</span>
            <span class="status-value" :class="message.includes('Success') ? 'success' : 'pending'">{{ message }}</span>
          </div>
        </div>
      </div>

      <!-- 输入区 -->
      <div class="modern-card input-card">
        <div class="card-header">
          <div class="card-icon input-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path d="M14 2H6C4.9 2 4 2.9 4 4V20C4 21.1 4.89 22 5.99 22H18C19.1 22 20 21.1 20 20V8L14 2Z" stroke="currentColor" stroke-width="2" fill="none"/>
              <polyline points="14,2 14,8 20,8" stroke="currentColor" stroke-width="2"/>
            </svg>
          </div>
          <h3>Function Input</h3>
          <div class="input-type-pills">
            <button 
              v-for="type in inputTypes" 
              :key="type.value"
              @click="inputType = type.value"
              :class="['type-pill', { active: inputType === type.value }]"
            >
              {{ type.label }}
            </button>
          </div>
        </div>
        
        <div class="card-content">
          <!-- 变量数 n (需要时显示) -->
          <div v-if="needsN" class="modern-input-group">
            <label class="modern-label">
              <span class="label-text">变量数 (n)</span>
              <span class="label-hint">真值表长度: {{ Math.pow(2, n) }}</span>
            </label>
            <div class="number-spinner-wrapper">
              <input
                  v-model.number="n"
                  type="number"
                  min="1"
                  max="10"
                  class="number-spinner-input"
              >
              <div class="spinner-controls">
                <button @click="n = Math.min(10, n + 1)" class="spinner-btn spinner-up" type="button">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
                    <path d="M7 15L12 10L17 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </button>
                <button @click="n = Math.max(1, n - 1)" class="spinner-btn spinner-down" type="button">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
                    <path d="M7 10L12 15L17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <!-- 动态输入框 -->
          <div class="modern-input-group">
            <!-- 真值表输入 -->
            <div v-if="inputType === 'truthTable'" class="input-section">
              <label class="modern-label">
                <span class="label-text">真值表</span>
                <span class="label-hint">✨ 支持多种格式: 0101 | 0,1,0,1 | [0,1,0,1] | {0,1,0,1} | "0101" 等</span>
              </label>
              <div class="input-group-vertical">
                <input
                    v-model="truthTableInput"
                    type="text"
                    placeholder="输入真值表 (支持各种分隔符和格式)，例如: 0101, 0,1,0,1, [0,1,0,1] 等"
                    class="modern-input code-input"
                >
                <label class="upload-button-block">
                  <input type="file" accept=".txt,.csv,.dat" @change="handleFileUpload" class="file-input">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                    <path d="M21 15V19C21 20.1 20.1 21 19 21H5C3.9 21 3 20.1 3 19V15" stroke="currentColor" stroke-width="2"/>
                    <polyline points="7,10 12,15 17,10" stroke="currentColor" stroke-width="2"/>
                    <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2"/>
                  </svg>
                  或上传文件
                </label>
              </div>
            </div>

            <!-- 十六进制输入 -->
            <div v-if="inputType === 'hex'" class="input-section">
              <label class="modern-label">
                <span class="label-text">十六进制值</span>
                <span class="label-hint">✨ 支持多种格式: 96 | 0x96 | AB CD | [AB, CD] | {AB,CD} 等</span>
              </label>
              <div class="input-group-vertical">
                <input
                    v-model="hexInput"
                    type="text"
                    placeholder="输入十六进制 (支持各种分隔符和格式)，例如: 96, 0xAB, A1 B2 等"
                    class="modern-input code-input"
                >
                <label class="upload-button-block">
                  <input type="file" accept=".txt,.hex" @change="handleFileUpload" class="file-input">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                    <path d="M21 15V19C21 20.1 20.1 21 19 21H5C3.9 21 3 20.1 3 19V15" stroke="currentColor" stroke-width="2"/>
                    <polyline points="7,10 12,15 17,10" stroke="currentColor" stroke-width="2"/>
                    <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2"/>
                  </svg>
                  或上传文件
                </label>
              </div>
            </div>

            <!-- 整数输入 -->
            <div v-if="inputType === 'int'" class="input-section">
              <label class="modern-label">
                <span class="label-text">整数值</span>
                <span class="label-hint">示例: 150</span>
              </label>
              <div class="input-group-vertical">
                <input
                    v-model.number="intInput"
                    type="number"
                    min="0"
                    placeholder="输入整数，例如: 150"
                    class="modern-input code-input"
                >
                <label class="upload-button-block">
                  <input type="file" accept=".txt" @change="handleFileUpload" class="file-input">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                    <path d="M21 15V19C21 20.1 20.1 21 19 21H5C3.9 21 3 20.1 3 19V15" stroke="currentColor" stroke-width="2"/>
                    <polyline points="7,10 12,15 17,10" stroke="currentColor" stroke-width="2"/>
                    <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2"/>
                  </svg>
                  或上传文件
                </label>
              </div>
            </div>

            <!-- ANF输入 -->
            <div v-if="inputType === 'anf'" class="input-section">
              <label class="modern-label">
                <span class="label-text">ANF 表达式</span>
                <span class="label-hint">示例: x0 + x1*x2 + 1</span>
              </label>
              <div class="input-group-vertical">
                <input
                    v-model="anfInput"
                    type="text"
                    placeholder="输入 ANF 表达式，例如: x0 + x1*x2"
                    class="modern-input code-input"
                >
                <label class="upload-button-block">
                  <input type="file" accept=".txt,.anf" @change="handleFileUpload" class="file-input">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                    <path d="M21 15V19C21 20.1 20.1 21 19 21H5C3.9 21 3 20.1 3 19V15" stroke="currentColor" stroke-width="2"/>
                    <polyline points="7,10 12,15 17,10" stroke="currentColor" stroke-width="2"/>
                    <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2"/>
                  </svg>
                  或上传文件
                </label>
              </div>
            </div>
          </div>

          <!-- 开始分析按钮 -->
          <div class="analyze-section">
            <button
                @click="analyzeFunction"
                :disabled="loading"
                class="modern-button primary analyze-btn"
            >
              <div v-if="loading" class="loading-spinner"></div>
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
                <polyline points="12,6 12,12 16,14" stroke="currentColor" stroke-width="2"/>
              </svg>
              {{ loading ? 'Analyzing...' : 'Start Analysis' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 结果展示区 -->
      <div v-if="analysisResult" class="modern-card results-card">
        <div class="card-header">
          <div class="card-icon results-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <rect x="3" y="3" width="18" height="18" rx="2" stroke="currentColor" stroke-width="2" fill="none"/>
              <line x1="9" y1="9" x2="15" y2="15" stroke="currentColor" stroke-width="2"/>
              <line x1="15" y1="9" x2="9" y2="15" stroke="currentColor" stroke-width="2"/>
            </svg>
          </div>
          <h3>Analysis Results</h3>
        </div>

        <!-- 错误显示 -->
        <div v-if="analysisResult.error" class="error-message">
          {{ analysisResult.error }}
        </div>

        <!-- 正常结果显示 -->
        <div v-else>
          <!-- 核心指标仪表板 -->
          <div class="dashboard-overview">
            <h4>核心指标概览</h4>
            <div class="dashboard-grid">
              <div class="dashboard-card">
                <div class="dashboard-icon">🔢</div>
                <div class="dashboard-label">变量数</div>
                <div class="dashboard-value">{{ analysisResult.n }}</div>
              </div>

              <div class="dashboard-card">
                <div class="dashboard-icon">⚖️</div>
                <div class="dashboard-label">汉明重量</div>
                <div class="dashboard-value">{{ analysisResult.hammingWeight }}</div>
              </div>

              <div class="dashboard-card">
                <div class="dashboard-icon">🎯</div>
                <div class="dashboard-label">非线性度</div>
                <div class="dashboard-value highlight">{{ analysisResult.nonlinearity }}</div>
              </div>

              <div class="dashboard-card">
                <div class="dashboard-icon">🛡️</div>
                <div class="dashboard-label">代数免疫度</div>
                <div class="dashboard-value highlight">{{ analysisResult.algebraicImmunity !== undefined && analysisResult.algebraicImmunity >= 0 ? analysisResult.algebraicImmunity : 'N/A' }}</div>
              </div>

              <div class="dashboard-card">
                <div class="dashboard-icon">⚡</div>
                <div class="dashboard-label">FAA</div>
                <div class="dashboard-value highlight">{{ analysisResult.faa !== undefined && analysisResult.faa >= 0 ? analysisResult.faa : 'N/A' }}</div>
              </div>

              <div class="dashboard-card">
                <div class="dashboard-icon">⚙️</div>
                <div class="dashboard-label">FAA (deg(g)≥1)</div>
                <div class="dashboard-value highlight">{{ analysisResult.faaWithPositiveDegree !== undefined && analysisResult.faaWithPositiveDegree >= 0 ? analysisResult.faaWithPositiveDegree : 'N/A' }}</div>
              </div>

              <div class="dashboard-card">
                <div class="dashboard-icon">🧩</div>
                <div class="dashboard-label">FAI</div>
                <div class="dashboard-value highlight">{{ analysisResult.fai !== undefined && analysisResult.fai >= 0 ? analysisResult.fai : 'N/A' }}</div>
              </div>

              <div class="dashboard-card">
                <div class="dashboard-icon">🧮</div>
                <div class="dashboard-label">代数次数</div>
                <div class="dashboard-value">{{ analysisResult.algebraicDegree }}</div>
              </div>

              <div class="dashboard-card">
                <div class="dashboard-icon">🔀</div>
                <div class="dashboard-label">差分均匀度</div>
                <div class="dashboard-value info">{{ analysisResult.differentialUniformity !== undefined ? analysisResult.differentialUniformity : 'N/A' }}</div>
              </div>

              <div class="dashboard-card">
                <div class="dashboard-icon">⚖️</div>
                <div class="dashboard-label">是否平衡</div>
                <div class="dashboard-value" :class="analysisResult.isBalanced ? 'success' : 'warning'">
                  {{ analysisResult.isBalanced ? '是' : '否' }}
                </div>
              </div>

              <div class="dashboard-card">
                <div class="dashboard-icon">🌟</div>
                <div class="dashboard-label">Bent函数</div>
                <div class="dashboard-value" :class="analysisResult.isBent ? 'success' : 'info'">
                  {{ analysisResult.isBent ? '是' : '否' }}
                </div>
              </div>
            </div>
          </div>

          <!-- 分类标签页 -->
          <div class="tab-container">
            <div class="tab-headers">
              <button
                  v-for="tab in tabs"
                  :key="tab.id"
                  @click="activeTab = tab.id"
                  :class="['tab-button', { active: activeTab === tab.id }]"
              >
                {{ tab.icon }} {{ tab.label }}
              </button>
            </div>

            <div class="tab-content">
              <!-- Tab 1: 基础性质 -->
              <div v-if="activeTab === 'basic'" class="tab-pane">
                <div class="property-section">
                  <h5>🔢 基础信息</h5>
                  <div class="property-grid">
                    <div class="property-item">
                      <span class="property-label">变量数 (n):</span>
                      <span class="property-value">{{ analysisResult.n }}</span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">真值表长度:</span>
                      <span class="property-value">{{ analysisResult.truthTable.length }}</span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">汉明重量:</span>
                      <span class="property-value">{{ analysisResult.hammingWeight }}</span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">是否平衡:</span>
                      <span class="property-value" :class="analysisResult.isBalanced ? 'success' : 'warning'">
                        {{ analysisResult.isBalanced ? '是' : '否' }}
                      </span>
                    </div>
                  </div>
                </div>

                <div class="property-section">
                  <h5>🧮 代数性质</h5>
                  <div class="property-item full-width">
                    <span class="property-label">ANF表达式:</span>
                    <div class="anf-display">
                      {{ analysisResult.anf || 'N/A' }}
                    </div>
                  </div>
                  <div class="property-grid">
                    <div class="property-item">
                      <span class="property-label">代数次数:</span>
                      <span class="property-value">{{ analysisResult.algebraicDegree }}</span>
                    </div>
                  </div>
                </div>

                <div class="property-section">
                  <h5>📊 真值表数据</h5>
                  <div class="data-display">
                    <div class="data-header">
                      <span>真值表</span>
                      <button @click="copyToClipboard(analysisResult.truthTable.join(','))" class="copy-btn">
                        📋 复制
                      </button>
                    </div>
                    <div class="data-content">
                      {{ analysisResult.truthTable.join(', ') }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Tab 2: 频谱分析 -->
              <div v-if="activeTab === 'spectral'" class="tab-pane">
                <div class="property-section">
                  <h5>📈 Walsh 频谱分析</h5>
                  <div class="data-display">
                    <div class="data-header">
                      <span>Walsh 频谱</span>
                      <button @click="copyToClipboard(analysisResult.walshSpectrum.join(','))" class="copy-btn">
                        📋 复制
                      </button>
                    </div>
                    <div class="data-content">
                      {{ analysisResult.walshSpectrum.join(', ') }}
                    </div>
                  </div>

                  <div class="distribution-display">
                    <h6>绝对Walsh谱分布</h6>
                    <table class="distribution-table">
                      <thead>
                      <tr>
                        <th>谱值</th>
                        <th>频次</th>
                      </tr>
                      </thead>
                      <tbody>
                      <tr v-for="(freq, value) in analysisResult.absoluteWalshSpectrum" :key="value">
                        <td>{{ value }}</td>
                        <td>{{ freq }}</td>
                      </tr>
                      </tbody>
                    </table>
                  </div>
                </div>

                <div class="property-section">
                  <h5>📊 自相关分析</h5>
                  <div class="data-display">
                    <div class="data-header">
                      <span>自相关谱</span>
                      <button @click="copyToClipboard(analysisResult.autocorrelationSpectrum.join(','))" class="copy-btn">
                        📋 复制
                      </button>
                    </div>
                    <div class="data-content">
                      {{ analysisResult.autocorrelationSpectrum.join(', ') }}
                    </div>
                  </div>

                  <div class="distribution-display">
                    <h6>绝对自相关谱分布</h6>
                    <table class="distribution-table">
                      <thead>
                      <tr>
                        <th>谱值</th>
                        <th>频次</th>
                      </tr>
                      </thead>
                      <tbody>
                      <tr v-for="(freq, value) in analysisResult.absoluteAutocorrelationSpectrum" :key="value">
                        <td>{{ value }}</td>
                        <td>{{ freq }}</td>
                      </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>

              <!-- Tab 3: 密码学性质 -->
              <div v-if="activeTab === 'crypto'" class="tab-pane">
                <div class="property-section">
                  <h5>🔒 抗攻击性质</h5>
                  <div class="property-grid">
                    <div class="property-item">
                      <span class="property-label">非线性度:</span>
                      <span class="property-value highlight">{{ analysisResult.nonlinearity }}</span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">代数免疫度:</span>
                      <span class="property-value highlight">{{ analysisResult.algebraicImmunity !== undefined && analysisResult.algebraicImmunity >= 0 ? analysisResult.algebraicImmunity : 'N/A' }}</span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">FAA:</span>
                      <span class="property-value highlight">{{ analysisResult.faa !== undefined && analysisResult.faa >= 0 ? analysisResult.faa : 'N/A' }}</span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">FAA (deg(g)≥1):</span>
                      <span class="property-value highlight">{{ analysisResult.faaWithPositiveDegree !== undefined && analysisResult.faaWithPositiveDegree >= 0 ? analysisResult.faaWithPositiveDegree : 'N/A' }}</span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">FAI:</span>
                      <span class="property-value highlight">{{ analysisResult.fai !== undefined && analysisResult.fai >= 0 ? analysisResult.fai : 'N/A' }}</span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">差分均匀度:</span>
                      <span class="property-value info">{{ analysisResult.differentialUniformity !== undefined ? analysisResult.differentialUniformity : 'N/A' }}</span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">相关免疫度:</span>
                      <span class="property-value">{{ analysisResult.correlationImmunity }}</span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">弹性阶数:</span>
                      <span class="property-value">{{ analysisResult.resiliencyOrder >= 0 ? analysisResult.resiliencyOrder : 'N/A' }}</span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">透明度阶:</span>
                      <span class="property-value">{{ analysisResult.transparencyOrder !== undefined && analysisResult.transparencyOrder !== null ? analysisResult.transparencyOrder : 'N/A' }}</span>
                    </div>
                  </div>
                </div>

                <div class="property-section">
                  <h5>🎯 特殊函数类型</h5>
                  <div class="property-grid">
                    <div class="property-item">
                      <span class="property-label">是否Bent函数:</span>
                      <span class="property-value" :class="analysisResult.isBent ? 'success' : 'info'">
                        {{ analysisResult.isBent ? '是' : '否' }}
                      </span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">是否旋转对称:</span>
                      <span class="property-value" :class="analysisResult.isRotationSymmetric ? 'success' : 'info'">
                        {{ analysisResult.isRotationSymmetric ? '是' : '否' }}
                      </span>
                    </div>
                  </div>
                </div>

                <div class="property-section">
                  <h5>📐 其他指标</h5>
                  <div class="property-grid">
                    <div class="property-item">
                      <span class="property-label">平方和指标:</span>
                      <span class="property-value">{{ analysisResult.sumOfSquareIndicator }}</span>
                    </div>
                    <div class="property-item">
                      <span class="property-label">绝对指标:</span>
                      <span class="property-value info">{{ analysisResult.absoluteIndicator }} (开发中)</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Tab 4: 数据导出 -->
              <div v-if="activeTab === 'export'" class="tab-pane">
                <div class="property-section">
                  <h5>📁 快速复制</h5>
                  <div class="export-grid">
                    <button @click="copyToClipboard(analysisResult.truthTable.join(','))" class="export-btn">
                      📋 复制真值表
                    </button>
                    <button @click="copyToClipboard(analysisResult.walshSpectrum.join(','))" class="export-btn">
                      📋 复制Walsh谱
                    </button>
                    <button @click="copyToClipboard(analysisResult.autocorrelationSpectrum.join(','))" class="export-btn">
                      📋 复制自相关谱
                    </button>
                    <button @click="copyToClipboard(analysisResult.anf)" class="export-btn">
                      📋 复制ANF表达式
                    </button>
                  </div>
                </div>

                <div class="property-section">
                  <h5>🔽 批量导出</h5>
                  <div class="export-grid">
                    <button @click="exportJSON()" class="export-btn">
                      📄 导出完整JSON
                    </button>
                    <button @click="exportSummary()" class="export-btn">
                      📊 导出摘要信息
                    </button>
                  </div>
                </div>

                <div class="property-section">
                  <h5>🔗 外部工具格式</h5>
                  <div class="export-grid">
                    <button @click="exportSageMath()" class="export-btn">
                      🧮 SageMath格式
                    </button>
                    <button @click="exportMatlab()" class="export-btn">
                      📈 Matlab格式
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- 关闭 main-content -->
  </div>
  <!-- 关闭 modern-container -->
</template>

<style scoped>
/* 现代化白色背景 */
.modern-container {
  min-height: 100vh;
  background: #f8fafc;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 'Noto Sans SC', sans-serif;
  font-feature-settings: 'tnum';
}

/* Hero Section */
.hero-section {
  position: relative;
  padding: 60px 24px 40px;
  text-align: center;
  overflow: hidden;
  background: linear-gradient(135deg, #0ea5e9 0%, #2563eb 50%, #7c3aed 100%);
  border-radius: 0 0 32px 32px;
  margin-bottom: 40px;
}

.hero-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 25% 25%, rgba(255,255,255,0.1) 0%, transparent 50%),
    radial-gradient(circle at 75% 75%, rgba(255,255,255,0.08) 0%, transparent 50%);
  animation: heroFloat 20s ease-in-out infinite;
}

@keyframes heroFloat {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-10px) rotate(1deg); }
}

.hero-content {
  position: relative;
  max-width: 800px;
  margin: 0 auto;
  z-index: 1;
}

.hero-icon {
  margin-bottom: 24px;
  animation: iconFloat 3s ease-in-out infinite;
}

@keyframes iconFloat {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-8px); }
}

.hero-title {
  font-size: 3.5rem;
  font-weight: 700;
  color: white;
  margin: 0 0 16px 0;
  letter-spacing: -2px;
  line-height: 1.1;
}

.hero-subtitle {
  font-size: 1.125rem;
  color: rgba(255, 255, 255, 0.95);
  margin: 0 0 32px 0;
  font-weight: 400;
  line-height: 1.6;
}

.hero-stats {
  display: flex;
  justify-content: center;
  gap: 48px;
  margin-top: 48px;
}

.stat-item {
  text-align: center;
}

.stat-number {
  display: block;
  font-size: 2.5rem;
  font-weight: 700;
  color: white;
  line-height: 1;
}

.stat-label {
  display: block;
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 8px;
  font-weight: 500;
}

/* Main Content */
.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px 80px;
}

/* 现代化卡片 */
.modern-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 32px;
  margin-bottom: 32px;
  box-shadow: 
    0 8px 32px rgba(0, 0, 0, 0.12),
    0 2px 16px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.modern-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #0ea5e9, #2563eb, #7c3aed);
}

.modern-card:hover {
  transform: translateY(-4px);
  box-shadow: 
    0 20px 40px rgba(0, 0, 0, 0.15),
    0 8px 32px rgba(0, 0, 0, 0.12);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
  flex-wrap: wrap;
}

.card-icon {
  width: 48px;
  height: 48px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0ea5e9, #2563eb);
  color: white;
  flex-shrink: 0;
}

.card-header h3 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
  flex: 1;
}

/* 连接状态 */
.connection-status {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #ef4444;
  animation: pulse 2s infinite;
}

.connection-status.connected {
  background: #10b981;
}

@keyframes pulse {
  0% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.2); opacity: 0.7; }
  100% { transform: scale(1); opacity: 1; }
}

/* 输入类型pills */
.input-type-pills {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.type-pill {
  padding: 8px 16px;
  border-radius: 20px;
  border: 2px solid #e5e7eb;
  background: white;
  color: #6b7280;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.type-pill:hover {
  border-color: #0ea5e9;
  color: #0ea5e9;
  transform: translateY(-1px);
}

.type-pill.active {
  border-color: #0ea5e9;
  background: linear-gradient(135deg, #0ea5e9, #2563eb);
  color: white;
  box-shadow: 0 4px 12px rgba(14, 165, 233, 0.4);
}

/* 现代化输入组件 */
.modern-input-group {
  margin-bottom: 32px;
}

.modern-label {
  display: block;
  margin-bottom: 12px;
}

.label-text {
  display: block;
  font-size: 1rem;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.label-hint {
  display: block;
  font-size: 0.875rem;
  color: #6b7280;
  font-weight: 400;
}

.input-wrapper {
  position: relative;
}

.input-wrapper.auto-width {
  flex: 0 0 auto;
}

.modern-input {
  width: 100%;
  padding: 16px 20px;
  border: 2px solid #e5e7eb;
  border-radius: 16px;
  font-size: 1rem;
  font-family: inherit;
  background: white;
  transition: all 0.2s ease;
  outline: none;
}

.modern-input.adaptive-input {
  min-width: 100px;
  max-width: 400px;
  text-align: center;
  font-weight: 600;
  font-size: 1.125rem;
  font-variant-numeric: tabular-nums;
}

.modern-input:focus {
  border-color: #0ea5e9;
  box-shadow: 0 0 0 4px rgba(14, 165, 233, 0.1);
}

.modern-input.number-input {
  appearance: textfield;
  -moz-appearance: textfield;
}

.modern-input.number-input::-webkit-outer-spin-button,
.modern-input.number-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

/* 自定义数字选择器 */
.number-spinner-wrapper {
  position: relative;
  display: inline-flex;
  align-items: center;
  background: white;
  border: 2px solid #e5e7eb;
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.2s ease;
  min-width: 140px;
  max-width: 200px;
}

.number-spinner-wrapper:focus-within {
  border-color: #0ea5e9;
  box-shadow: 0 0 0 4px rgba(14, 165, 233, 0.1);
}

.number-spinner-input {
  flex: 1;
  border: none;
  outline: none;
  padding: 14px 16px;
  font-size: 1.125rem;
  font-weight: 600;
  text-align: center;
  background: transparent;
  color: #1a1a1a;
  font-variant-numeric: tabular-nums;
  -moz-appearance: textfield;
  appearance: textfield;
}

.number-spinner-input::-webkit-outer-spin-button,
.number-spinner-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.spinner-controls {
  display: flex;
  flex-direction: column;
  border-left: 1px solid #e5e7eb;
}

.spinner-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 24px;
  border: none;
  background: #f8fafc;
  color: #64748b;
  cursor: pointer;
  transition: all 0.15s ease;
  padding: 0;
}

.spinner-btn:hover {
  background: linear-gradient(135deg, #0ea5e9, #2563eb);
  color: white;
}

.spinner-btn:active {
  transform: scale(0.95);
}

.spinner-up {
  border-bottom: 1px solid #e5e7eb;
}

/* 垂直布局输入组 */
.input-group-vertical {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 代码风格输入框（真值表和ANF） */
.modern-input.code-input {
  font-family: 'JetBrains Mono', 'SF Mono', Monaco, monospace;
  font-size: 0.95rem;
  font-weight: 500;
  letter-spacing: 0.02em;
  padding: 14px 18px;
}

/* 块级上传按钮 */
.upload-button-block {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 20px;
  background: white;
  border: 2px dashed #cbd5e1;
  border-radius: 12px;
  color: #64748b;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  width: 100%;
}

.upload-button-block:hover {
  background: #f1f5f9;
  border-color: #0ea5e9;
  color: #0ea5e9;
}

.input-decoration {
  position: absolute;
  bottom: 0;
  left: 20px;
  right: 20px;
  height: 2px;
  background: linear-gradient(90deg, #0ea5e9, #2563eb);
  transform: scaleX(0);
  transition: transform 0.2s ease;
  border-radius: 1px;
}

.modern-input:focus + .input-decoration {
  transform: scaleX(1);
}

/* 文件上传现代化 */
.input-with-upload-modern {
  display: flex;
  gap: 16px;
  align-items: flex-end;
}

.input-with-upload-modern .input-wrapper {
  flex: 1;
}

.upload-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 24px;
  background: linear-gradient(135deg, #f8fafc, #f1f5f9);
  border: 2px solid #e2e8f0;
  border-radius: 16px;
  color: #475569;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.upload-button:hover {
  background: linear-gradient(135deg, #0ea5e9, #2563eb);
  border-color: #0ea5e9;
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(14, 165, 233, 0.3);
}

/* 现代化按钮 */
.modern-button {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  padding: 16px 32px;
  border-radius: 16px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
  position: relative;
  overflow: hidden;
}

.modern-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
  transition: left 0.5s ease;
}

.modern-button:hover::before {
  left: 100%;
}

.modern-button.primary {
  background: linear-gradient(135deg, #0ea5e9, #2563eb);
  color: white;
  box-shadow: 0 8px 20px rgba(14, 165, 233, 0.4);
}

.modern-button.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 30px rgba(14, 165, 233, 0.5);
}

.modern-button.secondary {
  background: white;
  color: #0ea5e9;
  border: 2px solid #0ea5e9;
}

.modern-button.secondary:hover {
  background: #0ea5e9;
  color: white;
  transform: translateY(-2px);
}

.modern-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
}

.analyze-section {
  margin-top: 40px;
  text-align: center;
}

.analyze-btn {
  padding: 20px 48px;
  font-size: 1.125rem;
  border-radius: 20px;
}

/* 加载动画 */
.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 状态显示 */
.status-display {
  margin-top: 16px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 12px;
  border-left: 4px solid #e2e8f0;
}

.status-label {
  font-weight: 500;
  color: #64748b;
  margin-right: 8px;
}

.status-value {
  font-weight: 600;
  color: #1e293b;
}

.status-value.success {
  color: #10b981;
}

.status-value.pending {
  color: #f59e0b;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .hero-title {
    font-size: 2.5rem;
  }
  
  .hero-stats {
    gap: 24px;
  }
  
  .stat-number {
    font-size: 2rem;
  }
  
  .main-content {
    padding: 0 16px 60px;
  }
  
  .modern-card {
    padding: 24px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .input-with-upload-modern {
    flex-direction: column;
    align-items: stretch;
  }
  
  .input-type-pills {
    width: 100%;
    justify-content: flex-start;
  }
}

/* 结果展示现代化样式 */
.results-card {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(25px);
}

.error-message {
  background: linear-gradient(135deg, #fecaca, #fca5a5);
  color: #b91c1c;
  padding: 20px;
  border-radius: 16px;
  border: 1px solid #f87171;
  font-size: 0.975rem;
  font-weight: 500;
}

/* 仪表板现代化 */
.dashboard-overview {
  margin-bottom: 40px;
}

.dashboard-overview h4 {
  font-size: 1.75rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 32px 0;
  text-align: center;
  background: linear-gradient(135deg, #0ea5e9, #2563eb, #7c3aed);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.dashboard-card {
  background: white;
  padding: 32px 24px;
  border-radius: 16px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border: 2px solid #e2e8f0;
  transition: all 0.2s ease;
  position: relative;
  overflow: hidden;
}

.dashboard-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #0ea5e9, #2563eb, #7c3aed);
}

.dashboard-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(14, 165, 233, 0.15);
  border-color: #0ea5e9;
}

.dashboard-icon {
  font-size: 2.5rem;
  margin-bottom: 12px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dashboard-label {
  font-size: 0.813rem;
  color: #94a3b8;
  margin-bottom: 12px;
  font-weight: 600;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.dashboard-value {
  font-family: 'JetBrains Mono', 'SF Mono', Monaco, monospace;
  font-size: 2.75rem;
  font-weight: 700;
  color: #0f172a;
  line-height: 1;
  font-variant-numeric: tabular-nums;
}

.dashboard-value.success {
  color: #10b981;
}

.dashboard-value.warning {
  color: #f59e0b;
}

.dashboard-value.info {
  color: #3b82f6;
}

.dashboard-value.highlight {
  background: linear-gradient(135deg, #0ea5e9, #2563eb);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 800;
  font-size: 3rem;
}

/* 标签页现代化 */
.tab-container {
  margin-top: 40px;
}

.tab-headers {
  display: flex;
  border-bottom: 1px solid #e2e8f0;
  margin-bottom: 32px;
  overflow-x: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.tab-headers::-webkit-scrollbar {
  display: none;
}

.tab-button {
  background: none;
  border: none;
  padding: 16px 32px;
  cursor: pointer;
  color: #64748b;
  font-size: 0.975rem;
  font-weight: 500;
  border-bottom: 3px solid transparent;
  transition: all 0.2s ease;
  white-space: nowrap;
  position: relative;
  border-radius: 0;
  height: auto;
  min-width: auto;
}

.tab-button:hover:not(.active) {
  background: rgba(102, 126, 234, 0.05);
  color: #667eea;
}

.tab-button.active {
  color: #667eea;
  border-bottom-color: #667eea;
  background: rgba(102, 126, 234, 0.05);
  font-weight: 600;
}

.tab-content {
  min-height: 400px;
}

.tab-pane {
  animation: fadeInUp 0.4s ease;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 属性展示现代化 */
.property-section {
  margin-bottom: 40px;
  background: white;
  padding: 32px;
  border-radius: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(229, 231, 235, 0.8);
}

.property-section h5 {
  margin: 0 0 24px 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1e293b;
  display: flex;
  align-items: center;
  gap: 12px;
}

.property-section h6 {
  margin: 24px 0 16px 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: #374151;
}

.property-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.property-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: linear-gradient(135deg, #f8fafc, #f1f5f9);
  border-radius: 16px;
  border-left: 4px solid #667eea;
  transition: all 0.2s ease;
}

.property-item:hover {
  background: linear-gradient(135deg, #f1f5f9, #e2e8f0);
  transform: translateX(4px);
}

.property-item.full-width {
  grid-column: 1 / -1;
  flex-direction: column;
  align-items: stretch;
}

.property-label {
  font-weight: 600;
  color: #64748b;
  font-size: 0.938rem;
  letter-spacing: 0.01em;
}

.property-value {
  font-family: 'JetBrains Mono', 'SF Mono', Monaco, monospace;
  font-weight: 700;
  color: #0f172a;
  font-size: 1.375rem;
  font-variant-numeric: tabular-nums;
}

.property-value.success {
  color: #10b981;
}

.property-value.warning {
  color: #f59e0b;
}

.property-value.info {
  color: #3b82f6;
}

.property-value.highlight {
  background: linear-gradient(135deg, #0ea5e9, #2563eb);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 800;
  font-size: 1.5rem;
}

/* ANF 显示现代化 */
.anf-display {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  background: white;
  padding: 20px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  font-size: 0.975rem;
  color: #1e293b;
  margin-top: 16px;
  word-break: break-all;
  line-height: 1.6;
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.06);
  max-height: 200px;
  overflow-y: auto;
}

.anf-display::-webkit-scrollbar {
  width: 8px;
}

.anf-display::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 4px;
}

.anf-display::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}

.anf-display::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* 数据显示现代化 */
.data-display {
  margin-bottom: 32px;
}

.data-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.data-header span {
  font-weight: 600;
  color: #374151;
  font-size: 1.125rem;
}

.copy-btn {
  background: linear-gradient(135deg, #f8fafc, #f1f5f9);
  color: #667eea;
  border: 1px solid #e2e8f0;
  font-size: 0.875rem;
  padding: 10px 20px;
  border-radius: 12px;
  height: auto;
}

.copy-btn:hover {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-color: #667eea;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.data-content {
  font-family: 'JetBrains Mono', 'SF Mono', Monaco, monospace;
  background: #f8fafc;
  padding: 20px 24px;
  border: 2px solid #e2e8f0;
  border-radius: 12px;
  max-height: 200px;
  overflow-y: auto;
  font-size: 1rem;
  font-weight: 500;
  line-height: 1.8;
  word-break: break-all;
  color: #0f172a;
  font-variant-numeric: tabular-nums;
}

.data-content::-webkit-scrollbar {
  width: 8px;
}

.data-content::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 4px;
}

.data-content::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}

.data-content::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* 表格现代化 */
.distribution-display {
  margin-top: 24px;
}

.distribution-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  border: 1px solid #e2e8f0;
}

.distribution-table th,
.distribution-table td {
  padding: 18px 24px;
  text-align: center;
  border-bottom: 1px solid #e2e8f0;
  font-size: 1.063rem;
}

.distribution-table th {
  background: linear-gradient(135deg, #f1f5f9, #e2e8f0);
  font-weight: 700;
  color: #475569;
  letter-spacing: 0.025em;
  text-transform: uppercase;
  font-size: 0.875rem;
}

.distribution-table td {
  font-family: 'JetBrains Mono', 'SF Mono', Monaco, monospace;
  color: #0f172a;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.distribution-table tr:last-child td {
  border-bottom: none;
}

.distribution-table tr:hover {
  background: rgba(102, 126, 234, 0.04);
}

/* 导出按钮现代化 */
.export-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
}

.export-btn {
  background: linear-gradient(135deg, #f8fafc, #f1f5f9);
  color: #667eea;
  border: 1px solid #e2e8f0;
  font-size: 0.975rem;
  padding: 16px 20px;
  text-align: center;
  border-radius: 12px;
  height: auto;
  font-weight: 500;
  transition: all 0.2s ease;
}

.export-btn:hover {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-color: #667eea;
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
}

h1 {
  font-size: 32px;
  font-weight: 400;
  color: #1a73e8;
  margin: 0 0 8px 0;
  letter-spacing: -0.5px;
}

p {
  font-size: 16px;
  color: #5f6368;
  margin: 0 0 32px 0;
  font-weight: 400;
}

.card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 1px 3px 0 rgba(60,64,67,.3), 0 4px 8px 3px rgba(60,64,67,.15);
  border: none;
  text-align: left;
  transition: box-shadow 0.2s cubic-bezier(0.4, 0.0, 0.2, 1);
}

.card:hover {
  box-shadow: 0 2px 6px 2px rgba(60,64,67,.15), 0 8px 24px 4px rgba(60,64,67,.15);
}

.card h3 {
  font-size: 20px;
  font-weight: 500;
  color: #202124;
  margin: 0 0 16px 0;
  letter-spacing: 0.25px;
}

/* Material Design 按钮 */
button {
  background: #1a73e8;
  color: white;
  border: none;
  border-radius: 24px;
  padding: 12px 24px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0.0, 0.2, 1);
  text-transform: none;
  letter-spacing: 0.25px;
  font-family: inherit;
  min-width: 64px;
  height: 40px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

button:hover:not(:disabled) {
  background: #1557b0;
  box-shadow: 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
}

button:active:not(:disabled) {
  background: #1557b0;
  box-shadow: 0 1px 2px 0 rgba(60,64,67,.3), 0 2px 6px 2px rgba(60,64,67,.15);
}

button:disabled {
  background: #f1f3f4;
  color: #80868b;
  cursor: default;
  box-shadow: none;
}

.response {
  margin-top: 16px;
  font-size: 14px;
  color: #5f6368;
}

/* Material Design 输入组件 */
.input-group {
  margin-bottom: 20px;
}

.input-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #3c4043;
  letter-spacing: 0.25px;
}

.input-select, .input-text, .input-number {
  width: 100%;
  max-width: 400px;
  padding: 12px 16px;
  border: 1px solid #dadce0;
  border-radius: 8px;
  font-size: 14px;
  font-family: inherit;
  background: white;
  transition: border-color 0.2s cubic-bezier(0.4, 0.0, 0.2, 1), box-shadow 0.2s cubic-bezier(0.4, 0.0, 0.2, 1);
}

.input-select:focus, .input-text:focus, .input-number:focus {
  outline: none;
  border-color: #1a73e8;
  box-shadow: 0 0 0 1px #1a73e8;
}

.input-hint {
  display: block;
  margin-top: 6px;
  color: #5f6368;
  font-size: 12px;
  line-height: 16px;
}

/* Material Design 文件上传 */
.input-with-upload {
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.input-with-upload .input-text,
.input-with-upload .input-number {
  flex: 1;
  max-width: none;
}

.file-input {
  display: none;
}

.file-label {
  background: #f8f9fa;
  color: #1a73e8;
  border: 1px solid #dadce0;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  transition: all 0.2s cubic-bezier(0.4, 0.0, 0.2, 1);
  text-decoration: none;
  height: 40px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.file-label:hover {
  background: #f1f3f4;
  box-shadow: 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
}

.analyze-button {
  background: #137333;
  font-size: 14px;
  padding: 12px 32px;
  margin-top: 24px;
  height: 48px;
  border-radius: 24px;
}

.analyze-button:hover:not(:disabled) {
  background: #0d652d;
}

/* Material Design 结果展示 */
.results-card {
  background: white;
}

.error-message {
  background: #fce8e6;
  color: #d93025;
  padding: 16px;
  border-radius: 8px;
  border-left: 4px solid #d93025;
  font-size: 14px;
}

/* Material Design 仪表板 */
.dashboard-overview {
  margin-bottom: 32px;
}

.dashboard-overview h4 {
  font-size: 24px;
  font-weight: 400;
  color: #202124;
  margin: 0 0 24px 0;
  text-align: center;
  letter-spacing: 0;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.dashboard-card {
  background: white;
  padding: 24px;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
  border: 1px solid #e8eaed;
  transition: all 0.2s cubic-bezier(0.4, 0.0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.dashboard-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #1a73e8, #4285f4);
}

.dashboard-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px 3px rgba(60,64,67,.15), 0 1px 3px 0 rgba(60,64,67,.3);
}

.dashboard-icon {
  font-size: 32px;
  margin-bottom: 12px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dashboard-label {
  font-size: 14px;
  color: #5f6368;
  margin-bottom: 8px;
  font-weight: 400;
  letter-spacing: 0.25px;
}

.dashboard-value {
  font-size: 28px;
  font-weight: 400;
  color: #202124;
  line-height: 1.2;
}

.dashboard-value.success {
  color: #137333;
}

.dashboard-value.warning {
  color: #f9ab00;
}

.dashboard-value.info {
  color: #1a73e8;
}

.dashboard-value.highlight {
  color: #1a73e8;
  font-weight: 500;
}

/* Material Design 标签页 */
.tab-container {
  margin-top: 32px;
}

.tab-headers {
  display: flex;
  border-bottom: 1px solid #e8eaed;
  margin-bottom: 24px;
  overflow-x: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.tab-headers::-webkit-scrollbar {
  display: none;
}

.tab-button {
  background: none;
  border: none;
  padding: 16px 24px;
  cursor: pointer;
  color: #5f6368;
  font-size: 14px;
  font-weight: 500;
  border-bottom: 2px solid transparent;
  transition: all 0.2s cubic-bezier(0.4, 0.0, 0.2, 1);
  white-space: nowrap;
  position: relative;
  letter-spacing: 0.25px;
  border-radius: 0;
  height: auto;
  min-width: auto;
}

.tab-button:hover:not(.active) {
  background: rgba(26, 115, 232, 0.04);
  color: #1a73e8;
}

.tab-button.active {
  color: #1a73e8;
  border-bottom-color: #1a73e8;
  background: rgba(26, 115, 232, 0.04);
}

.tab-content {
  min-height: 400px;
}

.tab-pane {
  animation: materialFadeIn 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
}

@keyframes materialFadeIn {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Material Design 属性展示 */
.property-section {
  margin-bottom: 32px;
  background: white;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
  border: 1px solid #e8eaed;
}

.property-section h5 {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 500;
  color: #202124;
  letter-spacing: 0.25px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.property-section h6 {
  margin: 20px 0 12px 0;
  font-size: 16px;
  font-weight: 500;
  color: #3c4043;
}

.property-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
}

.property-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  border-left: 4px solid #1a73e8;
  transition: all 0.2s cubic-bezier(0.4, 0.0, 0.2, 1);
}

.property-item:hover {
  background: #f1f3f4;
  transform: translateX(2px);
}

.property-item.full-width {
  grid-column: 1 / -1;
  flex-direction: column;
  align-items: stretch;
}

.property-label {
  font-weight: 500;
  color: #3c4043;
  font-size: 14px;
  letter-spacing: 0.25px;
}

.property-value {
  font-weight: 500;
  color: #202124;
  font-size: 16px;
}

.property-value.success {
  color: #137333;
}

.property-value.warning {
  color: #f9ab00;
}

.property-value.info {
  color: #1a73e8;
}

.property-value.highlight {
  color: #1a73e8;
  font-weight: 600;
}

/* Material Design ANF 显示 */
.anf-display {
  font-family: 'Roboto Mono', 'Courier New', monospace;
  background: white;
  padding: 16px;
  border: 1px solid #e8eaed;
  border-radius: 8px;
  font-size: 14px;
  color: #202124;
  margin-top: 12px;
  word-break: break-all;
  line-height: 1.6;
  box-shadow: inset 0 1px 2px 0 rgba(60,64,67,.3);
}

/* Material Design 数据显示 */
.data-display {
  margin-bottom: 24px;
}

.data-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.data-header span {
  font-weight: 500;
  color: #3c4043;
  font-size: 16px;
}

.copy-btn {
  background: #f8f9fa;
  color: #1a73e8;
  border: 1px solid #dadce0;
  font-size: 12px;
  padding: 8px 16px;
  border-radius: 20px;
  height: 32px;
}

.copy-btn:hover {
  background: #f1f3f4;
  box-shadow: 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
}

.data-content {
  font-family: 'Roboto Mono', 'Courier New', monospace;
  background: white;
  padding: 16px;
  border: 1px solid #e8eaed;
  border-radius: 8px;
  max-height: 140px;
  overflow-y: auto;
  font-size: 13px;
  line-height: 1.6;
  word-break: break-all;
  box-shadow: inset 0 1px 2px 0 rgba(60,64,67,.3);
}

.data-content::-webkit-scrollbar {
  width: 8px;
}

.data-content::-webkit-scrollbar-track {
  background: #f1f3f4;
  border-radius: 4px;
}

.data-content::-webkit-scrollbar-thumb {
  background: #dadce0;
  border-radius: 4px;
}

.data-content::-webkit-scrollbar-thumb:hover {
  background: #bdc1c6;
}

/* Material Design 表格 */
.distribution-display {
  margin-top: 20px;
}

.distribution-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
  border: 1px solid #e8eaed;
}

.distribution-table th,
.distribution-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #e8eaed;
  font-size: 14px;
}

.distribution-table th {
  background: #f8f9fa;
  font-weight: 500;
  color: #3c4043;
  letter-spacing: 0.25px;
}

.distribution-table td {
  color: #202124;
}

.distribution-table tr:last-child td {
  border-bottom: none;
}

.distribution-table tr:hover {
  background: rgba(26, 115, 232, 0.04);
}

/* Material Design 导出按钮 */
.export-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}

.export-btn {
  background: #f8f9fa;
  color: #1a73e8;
  border: 1px solid #dadce0;
  font-size: 14px;
  padding: 12px 16px;
  text-align: center;
  border-radius: 8px;
  height: 40px;
}

.export-btn:hover {
  background: #f1f3f4;
  box-shadow: 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
}

/* Material Design 响应式 */
@media (max-width: 768px) {
  .container {
    padding: 16px;
  }

  .card {
    padding: 20px;
    margin-bottom: 16px;
  }

  .dashboard-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .dashboard-card {
    padding: 20px;
  }

  .dashboard-value {
    font-size: 24px;
  }

  .tab-headers {
    justify-content: flex-start;
  }

  .tab-button {
    font-size: 13px;
    padding: 14px 20px;
  }

  .property-grid {
    grid-template-columns: 1fr;
  }

  .property-section {
    padding: 20px;
  }

  .input-select, .input-text, .input-number {
    max-width: none;
  }

  .input-with-upload {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .export-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .container {
    padding: 12px;
  }

  h1 {
    font-size: 28px;
  }

  .dashboard-grid {
    grid-template-columns: 1fr;
  }

  .dashboard-card {
    padding: 16px;
  }

  .tab-button {
    font-size: 12px;
    padding: 12px 16px;
  }

  .property-item {
    padding: 12px;
  }
}

/* Material Design 加载状态 */
@keyframes materialSpin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

button:disabled::before {
  content: '';
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid #80868b;
  border-top: 2px solid transparent;
  border-radius: 50%;
  animation: materialSpin 1s linear infinite;
  margin-right: 8px;
}

/* Material Design 焦点环 */
button:focus-visible,
.tab-button:focus-visible {
  outline: 2px solid #1a73e8;
  outline-offset: 2px;
}

input:focus-visible,
select:focus-visible {
  outline: none;
}

/* Material Design elevation */
.elevation-1 {
  box-shadow: 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
}

.elevation-2 {
  box-shadow: 0 1px 2px 0 rgba(60,64,67,.3), 0 2px 6px 2px rgba(60,64,67,.15);
}

.elevation-3 {
  box-shadow: 0 4px 8px 3px rgba(60,64,67,.15), 0 1px 3px 0 rgba(60,64,67,.3);
}

/* 字体优化 - 数字和代码使用等宽字体 */
.dashboard-value,
.stat-number,
.property-value,
.number-spinner-input,
.adaptive-input,
input[type="number"],
.data-value {
  font-family: 'JetBrains Mono', 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, monospace;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
}

/* ANF 表达式使用等宽字体 */
.anf-display,
.data-display code,
pre {
  font-family: 'JetBrains Mono', 'SF Mono', Monaco, monospace;
  font-variant-numeric: tabular-nums;
}

/* 中文标题和文本优化 */
h1, h2, h3, h4, h5, h6,
.label-text,
.hero-title,
.hero-subtitle {
  font-family: 'Noto Sans SC', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  font-weight: 600;
  letter-spacing: -0.01em;
}

/* 正文文本 */
p, span, .label-hint, .property-label {
  font-family: 'Noto Sans SC', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  letter-spacing: 0.01em;
}
</style>
