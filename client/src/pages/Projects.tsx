/**
 * Projects 页面
 * 展示所有项目的集合
 */

import ProjectCard from "@/components/ProjectCard";
import Navigation from "@/components/Navigation";
import Footer from "@/components/Footer";
import { useFetchData } from "@/hooks/useFetchData";
import { Helmet } from "react-helmet-async";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { siteConfig } from "@/config/site";

interface Project {
  id: string;
  title: string;
  description: string;
  tags: string[];
  link?: string;
  date: string;
  detailsFile: string;
}

export default function Projects() {
  const { data: projects, loading, error } = useFetchData<Project>("/projects.json", "加载项目列表失败");

  return (
    <div className="min-h-screen flex flex-col bg-background font-sans selection:bg-primary/10">
      <Helmet>
        <title>项目集 - 我的作品</title>
        <meta name="description" content="浏览我的开源项目、个人作品和技术实验。" />
      </Helmet>
      <Navigation />

      <main className="flex-1">
        <section className="py-20 px-6 sm:px-8">
          <div className="max-w-3xl mx-auto space-y-12">
            {/* 极简头部 */}
            <div className="space-y-4 animate-fade-in-up" style={{ animationDelay: "0.1s" }}>
              <div className="flex items-center gap-3">
                <Avatar className="w-10 h-10 border border-border/50">
                  <AvatarImage src={siteConfig.avatar} alt={siteConfig.name} />
                  <AvatarFallback className="text-sm bg-muted text-muted-foreground">{siteConfig.name[0]}</AvatarFallback>
                </Avatar>
                <div className="leading-tight">
                  <div className="text-sm font-medium text-foreground">{siteConfig.name}</div>
                  <div className="text-xs text-muted-foreground">{siteConfig.title}</div>
                </div>
              </div>
              <h1 className="text-3xl sm:text-4xl font-bold text-foreground tracking-tight">项目集</h1>
              <p className="text-lg text-muted-foreground font-light max-w-2xl">
                这里是项目内容和一些作品。
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

            {!loading && !error && projects.length === 0 && (
              <div className="text-center py-20">
                <p className="text-muted-foreground text-lg">暂无项目</p>
              </div>
            )}

            {!loading && !error && projects.length > 0 && (
              <div className="space-y-0">
                {projects.map((project) => (
                  <ProjectCard
                    key={project.id}
                    id={project.id}
                    title={project.title}
                    description={project.description}
                    tags={project.tags}
                    link={project.link}
                    date={project.date}
                  />
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
