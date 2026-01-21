/**
 * 项目详情页面
 * 显示单个项目的完整信息
 */

import Footer from "@/components/Footer";
import Navigation from "@/components/Navigation";
import { Skeleton } from "@/components/ui/skeleton";
import { extractFrontmatter, stripLeadingHeading } from "@/lib/markdown";
import { useEffect, useState } from "react";
import { Helmet } from "react-helmet-async";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { useLocation, useParams } from "wouter";

interface ProjectMetadata {
  id: string;
  title: string;
  description: string;
  tags: string[];
  link?: string;
  detailsFile: string;
  date: string;
}

export default function ProjectDetail() {
  const params = useParams<{ id: string }>();
  const [, setLocation] = useLocation();
  const [project, setProject] = useState<ProjectMetadata | null>(null);
  const [content, setContent] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const projectId = params.id;

  useEffect(() => {
    const loadProject = async () => {
      try {
        setLoading(true);
        setError(null);
        setProject(null);
        setContent("");

        // 检查是否有项目 ID
        if (!projectId) {
          setError("未指定项目 ID");
          return;
        }

        // 加载项目列表
        const response = await fetch("/projects.json");
        if (!response.ok) {
          throw new Error(`Failed to fetch projects.json: ${response.status}`);
        }
        const projects: ProjectMetadata[] = await response.json();

        // 查找对应的项目
        const foundProject = projects.find((p) => p.id === projectId);
        if (!foundProject) {
          setError(`项目未找到 (ID: ${projectId})`);
          return;
        }

        setProject(foundProject);

        // 加载项目详情 Markdown 文件
        const markdownUrl = `/projects/${foundProject.detailsFile}`;

        const detailsResponse = await fetch(markdownUrl);
        if (!detailsResponse.ok) {
          throw new Error(`Failed to fetch markdown: ${detailsResponse.status}`);
        }

        const markdown = await detailsResponse.text();
        const { content: markdownContent } = extractFrontmatter(markdown);
        // 直接存储 Markdown 内容，交给 ReactMarkdown 渲染
        setContent(stripLeadingHeading(markdownContent));
      } catch (err) {
        console.error("Error loading project:", err);
        setError(
          `加载项目详情时出错: ${err instanceof Error ? err.message : "未知错误"}`
        );
      } finally {
        setLoading(false);
      }
    };

    void loadProject();
  }, [projectId]);

  if (loading) {
    return (
      <div className="min-h-screen flex flex-col bg-background font-sans selection:bg-primary/10">
        <Navigation />
        <main className="flex-1">
          <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16 space-y-8">
             <div className="space-y-4">
               <Skeleton className="h-12 w-3/4" />
               <Skeleton className="h-6 w-full" />
               <Skeleton className="h-6 w-2/3" />
             </div>
             <div className="flex gap-4">
               <Skeleton className="h-6 w-20" />
               <Skeleton className="h-6 w-20" />
             </div>
             <div className="space-y-4 pt-8">
               <Skeleton className="h-4 w-full" />
               <Skeleton className="h-4 w-full" />
               <Skeleton className="h-4 w-full" />
               <Skeleton className="h-4 w-4/5" />
             </div>
          </div>
        </main>
        <Footer />
      </div>
    );
  }

  if (error || !project) {
    return (
      <div className="min-h-screen flex flex-col bg-background font-sans selection:bg-primary/10">
        <Navigation />
        <main className="flex-1">
          <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
            <div className="text-center">
              <h1 className="text-2xl font-bold text-foreground mb-4">
                {error || "项目未找到"}
              </h1>
              <p className="text-sm text-muted-foreground mb-6">
                {projectId && `项目 ID: ${projectId}`}
              </p>
              <button
                onClick={() => setLocation("/projects")}
                className="inline-block px-6 py-2 text-sm font-semibold text-primary-foreground bg-primary hover:bg-primary/90 transition-colors rounded-sm"
              >
                返回项目集
              </button>
            </div>
          </div>
        </main>
        <Footer />
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col bg-background font-sans selection:bg-primary/10">
      <Helmet>
        <title>{project.title} - 我的项目</title>
        <meta name="description" content={project.description} />
      </Helmet>
      
      <Navigation />
      <main className="flex-1">
        <article className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          {/* 返回按钮 */}
          <div className="mb-8">
            <button
              onClick={() => setLocation("/projects")}
              className="text-sm text-muted-foreground hover:text-primary transition-colors flex items-center gap-1.5"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" /></svg>
              返回项目集
            </button>
          </div>

          {/* 项目头部信息 */}
          <header className="mb-12 pb-8 border-b border-border">
            <h1 className="text-4xl md:text-5xl font-bold text-foreground mb-4">
              {project.title}
            </h1>

            <p className="text-lg text-muted-foreground leading-relaxed mb-6">
              {project.description}
            </p>

            {/* 项目元数据 */}
            <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground mb-6">
              <span>{project.date}</span>
              <span>•</span>
              <div className="flex flex-wrap gap-2">
                {project.tags.map((tag) => (
                  <span
                    key={tag}
                    className="px-2.5 py-0.5 bg-secondary text-secondary-foreground rounded-full text-xs font-medium"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </div>

            {/* 项目链接 */}
            {project.link && (
              <a
                href={project.link}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 px-6 py-2.5 text-sm font-semibold text-primary-foreground bg-primary hover:bg-primary/90 transition-all rounded-lg shadow-sm hover:shadow"
              >
                查看源代码
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" /></svg>
              </a>
            )}
          </header>

          {/* 项目详情内容 */}
          {content ? (
            <div className="prose dark:prose-invert max-w-none prose-headings:font-bold prose-headings:tracking-tight prose-a:text-primary prose-img:rounded-xl prose-img:shadow-lg">
              <ReactMarkdown remarkPlugins={[remarkGfm]}>
                {content}
              </ReactMarkdown>
            </div>
          ) : (
            <p className="text-center text-muted-foreground">暂无项目详情内容</p>
          )}

          {/* 项目底部导航 */}
          <div className="mt-16 pt-8 border-t border-border">
            <button
              onClick={() => setLocation("/projects")}
              className="inline-flex items-center gap-2 px-6 py-2.5 text-sm font-semibold text-primary-foreground bg-primary hover:bg-primary/90 transition-all rounded-lg"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" /></svg>
              返回项目集
            </button>
          </div>
        </article>
      </main>
      <Footer />
    </div>
  );
}
