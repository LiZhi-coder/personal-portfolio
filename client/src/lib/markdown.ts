/**
 * Markdown 工具库
 * 提供 Markdown 处理辅助函数
 */

/**
 * 从 Markdown 前置元数据中提取标题、日期等信息
 * 格式: ---\ntitle: ...\ndate: ...\n---\n内容
 */
export function extractFrontmatter(markdown: string): {
  metadata: Record<string, string>;
  content: string;
} {
  const frontmatterRegex = /^---\n([\s\S]*?)\n---\n([\s\S]*)$/;
  const match = markdown.match(frontmatterRegex);

  if (!match) {
    return { metadata: {}, content: markdown };
  }

  const metadata: Record<string, string> = {};
  const frontmatterLines = match[1].split('\n');

  frontmatterLines.forEach((line) => {
    const [key, ...valueParts] = line.split(':');
    if (key && valueParts.length > 0) {
      metadata[key.trim()] = valueParts.join(':').trim();
    }
  });

  return {
    metadata,
    content: match[2],
  };
}

/**
 * 移除开头的标题（通常是与页面标题重复的 H1）
 */
export function stripLeadingHeading(markdown: string): string {
  return markdown.replace(/^\s*#{1,6}\s+.*(?:\r?\n)+/, "");
}

/**
 * 生成摘要（截取前 N 个字符）
 * 去除 HTML 标签和 Markdown 符号
 */
export function generateExcerpt(content: string, length: number = 150): string {
  // 简单去除 HTML 标签
  let plainText = content.replace(/<[^>]*>/g, '');
  // 去除常见 Markdown 符号 (如 #, *, `, [], ())
  plainText = plainText.replace(/[#*`\[\]()]/g, '');
  // 压缩空白
  plainText = plainText.replace(/\s+/g, ' ').trim();
  
  return plainText.length > length ? `${plainText.substring(0, length)}...` : plainText;
}

/**
 * 获取阅读时间估计（基于字数）
 */
export function getReadingTime(content: string): number {
  const wordCount = content.split(/\s+/).length;
  // 假设平均阅读速度为 200 词/分钟 (中文可以按字符数估算，这里简化处理)
  return Math.max(1, Math.ceil(wordCount / 200));
}
