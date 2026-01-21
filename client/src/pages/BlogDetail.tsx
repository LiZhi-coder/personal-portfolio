/**
 * 博客详情页面
 * 显示单篇博客文章的完整内容
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

interface BlogPost {
  id: string;
  title: string;
  excerpt: string;
  date: string;
  readingTime: number;
  detailsFile: string;
}

export default function BlogDetail() {
  const params = useParams<{ id: string }>();
  const [, setLocation] = useLocation();
  const [post, setPost] = useState<BlogPost | null>(null);
  const [content, setContent] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const postId = params.id;

  useEffect(() => {
    const loadPost = async () => {
      try {
        setLoading(true);
        setError(null);
        setPost(null);
        setContent("");

        if (!postId) {
          setError("未指定文章 ID");
          return;
        }

        // 加载博客列表
        const response = await fetch("/blog.json");
        if (!response.ok) {
          throw new Error(`Failed to fetch blog.json: ${response.status}`);
        }
        const posts: BlogPost[] = await response.json();

        // 查找对应的文章
        const foundPost = posts.find((p) => p.id === postId);
        if (!foundPost) {
          setError("文章未找到");
          return;
        }

        setPost(foundPost);

        // 加载博客详情 Markdown 文件
        const detailsResponse = await fetch(`/blog/${foundPost.detailsFile}`);
        if (!detailsResponse.ok) {
          setError("无法加载文章内容");
          return;
        }

        const markdown = await detailsResponse.text();
        const { content: markdownContent } = extractFrontmatter(markdown);
        setContent(stripLeadingHeading(markdownContent));
      } catch (err) {
        setError(
          `加载文章时出错: ${err instanceof Error ? err.message : "未知错误"}`
        );
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    void loadPost();
  }, [postId]);

  if (loading) {
    return (
      <div className="min-h-screen flex flex-col bg-background font-sans selection:bg-primary/10">
        <Navigation />
        <main className="flex-1">
          <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16 space-y-8">
            <div className="space-y-4">
              <Skeleton className="h-12 w-3/4" />
              <div className="flex gap-4">
                 <Skeleton className="h-4 w-24" />
                 <Skeleton className="h-4 w-24" />
              </div>
            </div>
            <div className="space-y-4 pt-8">
               <Skeleton className="h-4 w-full" />
               <Skeleton className="h-4 w-full" />
               <Skeleton className="h-4 w-full" />
               <Skeleton className="h-4 w-4/5" />
               <Skeleton className="h-4 w-full" />
            </div>
          </div>
        </main>
        <Footer />
      </div>
    );
  }

  if (error || !post) {
    return (
      <div className="min-h-screen flex flex-col bg-background font-sans selection:bg-primary/10">
        <Navigation />
        <main className="flex-1">
          <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
            <div className="text-center">
              <h1 className="text-2xl font-bold text-foreground mb-4">
                {error || "文章未找到"}
              </h1>
              <button
                onClick={() => setLocation("/blog")}
                className="inline-block px-6 py-2 text-sm font-semibold text-primary-foreground bg-primary hover:bg-primary/90 transition-colors rounded-sm"
              >
                返回博客列表
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
        <title>{post.title} - 我的博客</title>
        <meta name="description" content={post.excerpt} />
      </Helmet>
      
      <Navigation />
      <main className="flex-1">
        <article className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          {/* 返回链接 */}
          <div className="mb-8">
            <button
              onClick={() => setLocation("/blog")}
              className="text-sm text-muted-foreground hover:text-primary transition-colors flex items-center gap-1.5"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" /></svg>
              返回博客列表
            </button>
          </div>

          <header className="mb-12 pb-8 border-b border-border">
            <h1 className="text-4xl md:text-5xl font-bold text-foreground mb-6 tracking-tight">
              {post.title}
            </h1>

            <div className="flex flex-wrap items-center gap-6 text-sm text-muted-foreground">
              <time dateTime={post.date} className="flex items-center gap-1">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                {new Date(post.date).toLocaleDateString("zh-CN", {
                  year: "numeric",
                  month: "long",
                  day: "numeric",
                })}
              </time>
              <span className="flex items-center gap-1">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                {post.readingTime} 分钟阅读
              </span>
            </div>
          </header>

          <div className="prose dark:prose-invert max-w-none prose-headings:font-bold prose-headings:tracking-tight prose-a:text-primary prose-img:rounded-xl prose-img:shadow-lg">
             <ReactMarkdown remarkPlugins={[remarkGfm]}>
               {content}
             </ReactMarkdown>
          </div>

          <div className="mt-16 pt-8 border-t border-border">
            <button
              onClick={() => setLocation("/blog")}
              className="inline-flex items-center gap-2 px-6 py-2.5 text-sm font-semibold text-primary-foreground bg-primary hover:bg-primary/90 transition-all rounded-lg"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" /></svg>
              返回博客列表
            </button>
          </div>
        </article>
      </main>
      <Footer />
    </div>
  );
}
