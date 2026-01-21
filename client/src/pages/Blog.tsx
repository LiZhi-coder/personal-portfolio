/**
 * Blog 页面
 * 展示碎碎念和日常分享
 */

import Navigation from "@/components/Navigation";
import Footer from "@/components/Footer";
import { useFetchData } from "@/hooks/useFetchData";
import { Link } from "wouter";
import { Helmet } from "react-helmet-async";

interface BlogPost {
  id: string;
  title: string;
  excerpt: string;
  date: string;
  readingTime: number;
  detailsFile: string;
}

export default function Blog() {
  const { data: posts, loading, error } = useFetchData<BlogPost>("/blog.json", "加载博客列表失败");

  return (
    <div className="min-h-screen flex flex-col bg-background font-sans selection:bg-primary/10">
      <Helmet>
        <title>博客 - 碎碎念</title>
        <meta name="description" content="分享我的技术心得、开发经验和生活感悟。" />
      </Helmet>
      <Navigation />

      <main className="flex-1">
        <section className="py-20 px-6 sm:px-8">
          <div className="max-w-3xl mx-auto space-y-12">
            {/* 极简头部 */}
            <div className="space-y-4 animate-fade-in-up" style={{ animationDelay: "0.1s" }}>
              <h1 className="text-3xl sm:text-4xl font-bold text-foreground tracking-tight">碎碎念</h1>
              <p className="text-lg text-muted-foreground font-light max-w-2xl">
                胡言乱语，记录生活点滴。
              </p>
            </div>

            {loading && (
              <div className="flex items-center justify-center py-20">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
              </div>
            )}

            {error && (
              <div className="text-center py-20">
                <p className="text-destructive text-lg">{error}</p>
              </div>
            )}

            {!loading && !error && posts.length === 0 && (
              <div className="text-center py-20">
                <p className="text-muted-foreground text-lg">暂无文章</p>
              </div>
            )}

            {!loading && !error && posts.length > 0 && (
              <div className="space-y-6 animate-fade-in-up" style={{ animationDelay: "0.2s" }}>
                {posts.map((post) => (
                  <article
                    key={post.id}
                    className="group relative p-6 bg-card border border-border/40 hover:border-border/60 hover:bg-card/80 hover:shadow-lg hover:-translate-y-1 rounded-2xl transition-all duration-300 cursor-pointer"
                  >
                    <div className="space-y-3">
                      <Link href={`/blog/${post.id}`}>
                        <span className="absolute inset-0" aria-hidden="true" />
                      </Link>

                      <div className="flex items-center justify-between gap-4">
                        <h3 className="text-xl font-semibold text-foreground group-hover:text-primary transition-colors">
                          {post.title}
                        </h3>
                        <time dateTime={post.date} className="text-xs text-muted-foreground whitespace-nowrap px-3 py-1 bg-muted/50 rounded-full border border-border/20">
                          {new Date(post.date).toLocaleDateString("zh-CN", {
                            month: "short",
                            day: "numeric",
                          })}
                        </time>
                      </div>

                      <p className="text-sm text-muted-foreground leading-relaxed line-clamp-2">
                        {post.excerpt}
                      </p>

                      <div className="flex items-center text-xs font-medium text-primary opacity-0 group-hover:opacity-100 transition-opacity transform translate-y-2 group-hover:translate-y-0">
                        阅读全文 →
                      </div>
                    </div>
                  </article>
                ))}
              </div>
            )}
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}
