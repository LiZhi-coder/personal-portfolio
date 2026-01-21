#!/usr/bin/env node
import { promises as fs } from "node:fs";
import path from "node:path";

const BLOG_DIR = path.join(process.cwd(), "client", "public", "blog");
const PROJECTS_DIR = path.join(process.cwd(), "client", "public", "projects");
const BLOG_OUTPUT_PATH = path.join(process.cwd(), "client", "public", "blog.json");
const PROJECTS_OUTPUT_PATH = path.join(
  process.cwd(),
  "client",
  "public",
  "projects.json"
);
const BLOG_EXCERPT_LENGTH = 150;
const PROJECT_DESCRIPTION_LENGTH = 160;

function normalizeNewlines(text) {
  return text.replace(/\r\n/g, "\n");
}

function cleanValue(value) {
  const trimmed = value.trim();
  if (
    (trimmed.startsWith('"') && trimmed.endsWith('"')) ||
    (trimmed.startsWith("'") && trimmed.endsWith("'"))
  ) {
    return trimmed.slice(1, -1).trim();
  }
  return trimmed;
}

function extractFrontmatter(markdown) {
  const normalized = normalizeNewlines(markdown).replace(/^\uFEFF/, "");
  const frontmatterRegex = /^---\n([\s\S]*?)\n---\n?([\s\S]*)$/;
  const match = normalized.match(frontmatterRegex);

  if (!match) {
    return { metadata: {}, content: normalized };
  }

  const metadata = {};
  const lines = match[1].split("\n");
  for (const line of lines) {
    if (!line.trim() || !line.includes(":")) {
      continue;
    }
    const [key, ...rest] = line.split(":");
    const value = rest.join(":");
    if (key && value) {
      metadata[key.trim()] = cleanValue(value);
    }
  }

  return { metadata, content: match[2] };
}

function parseDateString(value) {
  if (!value) {
    return null;
  }

  const trimmed = value.trim();
  const isoMatch = trimmed.match(/^(\d{4})[-/](\d{1,2})[-/](\d{1,2})/);
  if (isoMatch) {
    const [, year, month, day] = isoMatch;
    return new Date(Number(year), Number(month) - 1, Number(day));
  }

  const cnMatch = trimmed.match(
    /^(\d{4})\s*年\s*(\d{1,2})\s*月\s*(\d{1,2})\s*日/
  );
  if (cnMatch) {
    const [, year, month, day] = cnMatch;
    return new Date(Number(year), Number(month) - 1, Number(day));
  }

  const parsed = new Date(trimmed);
  if (!Number.isNaN(parsed.getTime())) {
    return parsed;
  }

  return null;
}

function formatDate(date) {
  const year = String(date.getFullYear());
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

function stripMarkdown(markdown) {
  let text = normalizeNewlines(markdown);
  text = text.replace(/^\s*#{1,6}\s+.*\n/, "");
  text = text.replace(/```[\s\S]*?```/g, " ");
  text = text.replace(/`[^`]*`/g, " ");
  text = text.replace(/!\[[^\]]*\]\([^)]*\)/g, " ");
  text = text.replace(/\[([^\]]+)\]\([^)]*\)/g, "$1");
  text = text.replace(/[*_~]/g, "");
  text = text.replace(/^\s{0,3}#{1,6}\s+/gm, "");
  text = text.replace(/^\s{0,3}>\s?/gm, "");
  text = text.replace(/^\s{0,3}([-*+]|\d+\.)\s+/gm, "");
  text = text.replace(/^\s*[-*_]{3,}\s*$/gm, " ");
  text = text.replace(/^\s*\|.*\|\s*$/gm, " ");
  text = text.replace(/^.*(发布时间|阅读时间).*$\n?/gm, " ");
  text = text.replace(/\s+/g, " ").trim();
  return text;
}

function sliceFromFirstHeading(markdown) {
  const match = markdown.match(/^\s*#{1,6}\s+.*$/m);
  if (!match || match.index === undefined) {
    return markdown;
  }
  return markdown.slice(match.index);
}

function generateExcerpt(plainText, length) {
  if (!plainText) {
    return "";
  }
  return plainText.length > length
    ? `${plainText.slice(0, length)}...`
    : plainText;
}

function getReadingTime(plainText) {
  if (!plainText) {
    return 1;
  }
  const cjkCount = (plainText.match(/[\u4e00-\u9fff]/g) || []).length;
  const wordCount = plainText
    .replace(/[\u4e00-\u9fff]/g, " ")
    .trim()
    .split(/\s+/)
    .filter(Boolean).length;
  const total = cjkCount + wordCount;
  return Math.max(1, Math.ceil(total / 200));
}

function pickTitle(metadata, content, fallback) {
  if (metadata.title) {
    return metadata.title;
  }
  const headingMatch = content.match(/^#{1,6}\s+(.+)\s*$/m);
  if (headingMatch) {
    return headingMatch[1].trim().replace(/\s+#+$/, "");
  }
  return fallback;
}

function parseReadingTime(value) {
  if (!value) {
    return null;
  }
  const parsed = Number.parseInt(String(value), 10);
  return Number.isNaN(parsed) ? null : parsed;
}

function parseTags(value) {
  if (!value) {
    return [];
  }
  const raw = String(value).trim();
  if (!raw) {
    return [];
  }
  if (raw.startsWith("[") && raw.endsWith("]")) {
    const inner = raw.slice(1, -1).trim();
    if (!inner) {
      return [];
    }
    return inner
      .split(",")
      .map((tag) => cleanValue(tag))
      .filter(Boolean);
  }
  return raw
    .split(/[,，|]/)
    .map((tag) => cleanValue(tag))
    .filter(Boolean);
}

function pickDescription(metadata, content, length) {
  const direct =
    metadata.description || metadata.summary || metadata.excerpt || "";
  if (direct) {
    return String(direct).trim();
  }

  const plain = stripMarkdown(content);
  return generateExcerpt(plain, length);
}

function pickLink(metadata) {
  return metadata.link || metadata.url || metadata.repo || "";
}

function normalizeLine(line) {
  return line
    .replace(/^\s{0,3}>\s?/, "")
    .replace(/[*_`~]/g, "")
    .trim();
}

function extractInlineTags(content) {
  const lines = normalizeNewlines(content).split("\n");
  for (const line of lines) {
    const cleaned = normalizeLine(line);
    const match = cleaned.match(
      /^(项目关键词|关键词|Tags?|Tag)\s*[:：]\s*(.+)$/i
    );
    if (match) {
      return parseTags(match[2]);
    }
  }
  return [];
}

function extractInlineLink(content) {
  const lines = normalizeNewlines(content).split("\n");
  for (const line of lines) {
    const cleaned = normalizeLine(line);
    const mdLinkMatch = cleaned.match(/\[[^\]]+]\((https?:\/\/[^)]+)\)/i);
    if (mdLinkMatch) {
      return mdLinkMatch[1].trim();
    }
    const urlMatch = cleaned.match(/https?:\/\/\S+/i);
    if (urlMatch) {
      return urlMatch[0].replace(/[),.]+$/, "");
    }
  }
  return "";
}

function removeProjectMetaLines(content) {
  const lines = normalizeNewlines(content).split("\n");
  const filtered = lines.filter((line) => {
    const cleaned = normalizeLine(line);
    return !/^(项目关键词|关键词|Tags?|Tag|源码仓库|仓库|Repo|Repository)\s*[:：]/i.test(
      cleaned
    );
  });
  return filtered.join("\n");
}

async function readMarkdownEntries(dir) {
  const entries = await fs.readdir(dir, { withFileTypes: true });
  const files = [];
  for (const entry of entries) {
    if (!entry.isFile() || !entry.name.toLowerCase().endsWith(".md")) {
      continue;
    }
    if (entry.name.toLowerCase() === "template.md") {
      continue;
    }
    const filePath = path.join(dir, entry.name);
    const raw = await fs.readFile(filePath, "utf8");
    const stats = await fs.stat(filePath);
    files.push({
      id: path.basename(entry.name, ".md"),
      fileName: entry.name,
      raw,
      stats,
    });
  }
  return files;
}

function sortByDateDesc(items) {
  return items.sort((a, b) => {
    const timeA = Number.isNaN(Date.parse(a.date)) ? 0 : Date.parse(a.date);
    const timeB = Number.isNaN(Date.parse(b.date)) ? 0 : Date.parse(b.date);
    if (timeA !== timeB) {
      return timeB - timeA;
    }
    return a.id.localeCompare(b.id);
  });
}

async function main() {
  const blogEntries = await readMarkdownEntries(BLOG_DIR);
  const posts = blogEntries.map((entry) => {
    const { metadata, content } = extractFrontmatter(entry.raw);
    const body = content.trim();
    const title = pickTitle(metadata, body, entry.id);
    const plain = stripMarkdown(body);
    const excerpt = metadata.excerpt
      ? String(metadata.excerpt).trim()
      : generateExcerpt(plain, BLOG_EXCERPT_LENGTH);
    const frontmatterDate = parseDateString(metadata.date);
    const date = formatDate(frontmatterDate ?? entry.stats.mtime);
    const readingTime =
      parseReadingTime(metadata.readingTime) ?? getReadingTime(plain);

    return {
      id: entry.id,
      title,
      excerpt,
      date,
      readingTime,
      detailsFile: entry.fileName,
    };
  });

  const projectEntries = await readMarkdownEntries(PROJECTS_DIR);
  const projects = projectEntries.map((entry) => {
    const { metadata, content } = extractFrontmatter(entry.raw);
    const body = content.trim();
    const title = pickTitle(metadata, body, entry.id);
    const projectBody = removeProjectMetaLines(sliceFromFirstHeading(body));
    const description = pickDescription(
      metadata,
      projectBody,
      PROJECT_DESCRIPTION_LENGTH
    );
    const frontmatterDate = parseDateString(metadata.date);
    const date = formatDate(frontmatterDate ?? entry.stats.mtime);
    const tagsFromMeta = parseTags(metadata.tags || metadata.tag);
    const tags = tagsFromMeta.length > 0 ? tagsFromMeta : extractInlineTags(body);
    const link = pickLink(metadata) || extractInlineLink(body);

    return {
      id: entry.id,
      title,
      description,
      tags,
      link: link || undefined,
      detailsFile: entry.fileName,
      date,
    };
  });

  sortByDateDesc(posts);
  sortByDateDesc(projects);

  await fs.writeFile(
    BLOG_OUTPUT_PATH,
    `${JSON.stringify(posts, null, 2)}\n`,
    "utf8"
  );
  await fs.writeFile(
    PROJECTS_OUTPUT_PATH,
    `${JSON.stringify(projects, null, 2)}\n`,
    "utf8"
  );

  console.log(
    `Generated ${posts.length} posts -> ${BLOG_OUTPUT_PATH}\n` +
      `Generated ${projects.length} projects -> ${PROJECTS_OUTPUT_PATH}`
  );
}

main().catch((error) => {
  console.error("Failed to generate blog.json:", error);
  process.exit(1);
});
