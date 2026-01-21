/**
 * ProjectCard 组件
 * 极简现代主义风格的项目展示卡片
 * 特点：细线分割、克制的设计、纯文本展示
 */

import { Link } from "wouter";

interface ProjectCardProps {
  id: string;
  title: string;
  description: string;
  tags: string[];
  link?: string;
  date?: string;
}

export default function ProjectCard({
  id,
  title,
  description,
  tags,
  link,
  date,
}: ProjectCardProps) {
  return (
    <article className="group relative p-6 mb-6 bg-card border border-border/40 hover:border-border/60 hover:bg-card/80 hover:shadow-lg hover:-translate-y-1 rounded-2xl transition-all duration-300 cursor-pointer">
      <Link href={`/projects/${id}`} aria-label={`查看项目：${title}`}>
        <span className="absolute inset-0" aria-hidden="true" />
      </Link>

      <div className="space-y-4">
        {/* 项目标题和日期 */}
        <div className="flex items-start justify-between gap-4">
          <h3 className="text-xl font-semibold text-foreground flex-1 group-hover:text-primary transition-colors">
            {title}
          </h3>
          {date && (
            <time className="relative z-10 text-xs text-muted-foreground whitespace-nowrap px-3 py-1 bg-muted/50 rounded-full border border-border/20">
              {date}
            </time>
          )}
        </div>

        {/* 项目描述 */}
        <p className="text-sm text-muted-foreground leading-relaxed">
          {description}
        </p>

        {/* 技术标签 */}
        {tags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {tags.map((tag) => (
              <span
                key={tag}
                className="inline-block px-3 py-1 text-xs font-medium bg-primary/10 text-primary border border-primary/20 rounded-full"
              >
                {tag}
              </span>
            ))}
          </div>
        )}

        {/* 项目链接 */}
        <div className="flex gap-4 pt-2">
          <Link
            href={`/projects/${id}`}
            className="relative z-10 inline-flex items-center gap-1 text-sm text-primary hover:gap-2 transition-all font-semibold"
          >
            查看详情
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
          </Link>
          {link && (
            <a
              href={link}
              target="_blank"
              rel="noopener noreferrer"
              className="relative z-10 inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-primary hover:gap-2 transition-all"
            >
              源代码
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
            </a>
          )}
        </div>
      </div>
    </article>
  );
}
