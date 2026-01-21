---
title: 性能问题
date: 2026-01-06 10:00:00
excerpt: 网页性能加载问题
readingTime: 
---

# 网页性能优化思路
## 为什么性能很重要？

网页性能直接影响用户体验。研究表明，页面加载速度每延迟 1 秒，就会导致 7% 的转化率下降。因此，优化网页性能是每个开发者都应该关注的问题。

## Core Web Vitals

Google 提出的 Core Web Vitals 是衡量网页性能的三个关键指标：

### 1. LCP (Largest Contentful Paint)

LCP 衡量的是页面主要内容加载的速度。理想的 LCP 应该在 2.5 秒以内。

**优化方法**：
- 优化服务器响应时间
- 减少关键渲染路径中的资源
- 使用 CDN 分发内容
- 实现懒加载

### 2. FID (First Input Delay)

FID 衡量的是用户与页面交互时的响应速度。理想的 FID 应该在 100 毫秒以内。

**优化方法**：
- 减少 JavaScript 执行时间
- 将长任务分解为较小的任务
- 使用 Web Workers 处理复杂计算
- 优化第三方脚本

### 3. CLS (Cumulative Layout Shift)

CLS 衡量的是页面布局的稳定性。理想的 CLS 应该在 0.1 以下。

**优化方法**：
- 为图片和视频指定尺寸
- 避免在现有内容上方插入新内容
- 使用 transform 而不是改变 width/height
- 预加载字体

## 性能优化技巧

### 1. 代码分割

使用动态导入将代码分割成更小的块，只在需要时加载。

```javascript
// 动态导入
const HeavyComponent = React.lazy(() => import('./HeavyComponent'));

function App() {
  return (
    <Suspense fallback={<Loading />}>
      <HeavyComponent />
    </Suspense>
  );
}
```

### 2. 图片优化

- 使用现代格式（WebP）
- 响应式图片（srcset）
- 懒加载
- 压缩图片

```html
<picture>
  <source srcset="image.webp" type="image/webp">
  <source srcset="image.jpg" type="image/jpeg">
  <img src="image.jpg" alt="Description" loading="lazy">
</picture>
```

### 3. 缓存策略

- 利用浏览器缓存
- 实现 Service Worker 缓存
- 使用 CDN

### 4. 减少 JavaScript

- 移除未使用的代码
- 使用 Tree Shaking
- 选择更小的库
- 延迟加载非关键脚本

## 性能监测

### 使用 Web Vitals 库

```javascript
import { getCLS, getFID, getFCP, getLCP, getTTFB } from 'web-vitals';

getCLS(console.log);
getFID(console.log);
getFCP(console.log);
getLCP(console.log);
getTTFB(console.log);
```

### 使用 Performance API

```javascript
// 测量特定操作的性能
performance.mark('operation-start');
// ... 执行操作
performance.mark('operation-end');
performance.measure('operation', 'operation-start', 'operation-end');

const measure = performance.getEntriesByName('operation')[0];
console.log(`Operation took ${measure.duration}ms`);
```

## 常见的性能问题

### 1. 阻塞渲染的资源

```html
<!-- ❌ 不好：阻塞渲染 -->
<script src="script.js"></script>

<!-- ✅ 好：异步加载 -->
<script src="script.js" async></script>

<!-- ✅ 好：延迟加载 -->
<script src="script.js" defer></script>
```

### 2. 未优化的 CSS

- 避免内联大量 CSS
- 使用 CSS-in-JS 库时要小心
- 移除未使用的 CSS

### 3. 过大的 Bundle

使用 webpack-bundle-analyzer 分析你的 bundle 大小。

```bash
npm install --save-dev webpack-bundle-analyzer
```


**有用的工具**：
- [Google PageSpeed Insights](https://pagespeed.web.dev/)
- [WebPageTest](https://www.webpagetest.org/)
- [Lighthouse](https://developers.google.com/web/tools/lighthouse)
