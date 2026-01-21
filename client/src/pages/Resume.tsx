/**
 * Resume 页面
 * 展示教育背景、工作经历和技能
 * 数据驱动渲染，内容来自 /resume.json
 */

import Navigation from "@/components/Navigation";
import Footer from "@/components/Footer";
import TimelineItem from "@/components/TimelineItem";
import { useEffect, useState } from "react";
import { Briefcase, GraduationCap, Lightbulb, FolderPlus, Info, Download, Mail } from "lucide-react";
import type { ResumeData } from "@/types/resume";

export default function Resume() {
  const [data, setData] = useState<ResumeData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadResume = async () => {
      try {
        setLoading(true);
        const response = await fetch("/resume.json");
        if (!response.ok) throw new Error("Failed to load resume");
        const json: ResumeData = await response.json();
        setData(json);
      } catch (err) {
        console.error("Error loading resume:", err);
        setError("加载简历失败");
      } finally {
        setLoading(false);
      }
    };
    loadResume();
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen flex flex-col bg-background">
        <Navigation />
        <main className="flex-1 flex items-center justify-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
        </main>
        <Footer />
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="min-h-screen flex flex-col bg-background">
        <Navigation />
        <main className="flex-1 flex items-center justify-center">
          <p className="text-destructive text-lg">{error || "加载失败"}</p>
        </main>
        <Footer />
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col bg-background">
      <Navigation />

      <main className="flex-1">
        <section className="py-16 px-4 sm:px-6 lg:px-8">
          <div className="max-w-3xl mx-auto">
            <div className="mb-12">
              <h1 className="text-4xl md:text-5xl font-bold text-foreground mb-4">个人简历</h1>
              <p className="text-lg text-muted-foreground">我的专业经历与技能</p>
            </div>

            {/* 工作经历 */}
            <div className="mb-16">
              <h3 className="text-2xl font-semibold text-foreground mb-8 flex items-center gap-2">
                <Briefcase className="w-6 h-6 text-primary" />
                工作经历
              </h3>
              <div className="space-y-8">
                {data.experience.map((exp, idx) => (
                  <TimelineItem
                    key={idx}
                    title={exp.title}
                    subtitle={`${exp.company} • ${exp.period}`}
                    isLast={idx === data.experience.length - 1}
                  >
                    <ul className="text-base text-foreground mt-3 space-y-2 leading-normal">
                      {exp.highlights.map((item, i) => (
                        <li key={i} className="flex gap-2">
                          <span className="text-primary mt-1">•</span>
                          <span>{item}</span>
                        </li>
                      ))}
                    </ul>
                  </TimelineItem>
                ))}
              </div>
            </div>

            {/* 教育背景 */}
            <div className="mb-16">
              <h3 className="text-2xl font-semibold text-foreground mb-8 flex items-center gap-2">
                <GraduationCap className="w-6 h-6 text-primary" />
                教育背景
              </h3>
              <div className="space-y-6">
                {data.education.map((edu, idx) => (
                  <TimelineItem
                    key={idx}
                    title={edu.degree}
                    subtitle={`${edu.school} • ${edu.type} • ${edu.period}`}
                    isLast={idx === data.education.length - 1}
                  >
                    <p className="text-base text-foreground mt-3 leading-normal">{edu.description}</p>
                  </TimelineItem>
                ))}
              </div>
            </div>

            {/* 技能 */}
            <div className="mb-12">
              <h3 className="text-2xl font-semibold text-foreground mb-8 flex items-center gap-2">
                <Lightbulb className="w-6 h-6 text-primary" />
                专业技能
              </h3>
              <div className="grid md:grid-cols-3 gap-6">
                {data.skills.map((skill, idx) => (
                  <div key={idx} className="space-y-4">
                    <h4 className="font-semibold text-foreground text-lg border-b border-border pb-2">{skill.category}</h4>
                    <ul className="text-base text-foreground space-y-2 leading-normal">
                      {skill.items.map((item, i) => (
                        <li key={i} className="flex items-center gap-2">
                          <span className="w-2 h-2 bg-primary rounded-full"></span>
                          {item}
                        </li>
                      ))}
                    </ul>
                  </div>
                ))}
              </div>
            </div>

            {/* 项目作品 */}
            <div className="mb-16">
              <h3 className="text-2xl font-semibold text-foreground mb-8 flex items-center gap-2">
                <FolderPlus className="w-6 h-6 text-primary" />
                项目作品
              </h3>
              <div className="space-y-8">
                {data.projects.map((project, idx) => (
                  <TimelineItem key={idx} title={project.title} isLast={idx === data.projects.length - 1}>
                    <ul className="text-base text-foreground mt-3 space-y-2 leading-normal">
                      {project.highlights.map((item, i) => (
                        <li key={i} className="flex gap-2">
                          <span className="text-primary mt-1">•</span>
                          <span>{item}</span>
                        </li>
                      ))}
                    </ul>
                  </TimelineItem>
                ))}
              </div>
            </div>

            {/* 个人特点 */}
            <div className="mb-16">
              <h3 className="text-2xl font-semibold text-foreground mb-8 flex items-center gap-2">
                <Info className="w-6 h-6 text-primary" />
                个人特点
              </h3>
              <ul className="text-base text-foreground space-y-2 leading-normal">
                {data.traits.map((trait, idx) => (
                  <li key={idx} className="flex gap-2">
                    <span className="text-primary mt-1">•</span>
                    <span>{trait}</span>
                  </li>
                ))}
              </ul>
            </div>

            {/* 下载简历按钮 */}
            <div className="mt-12 flex flex-wrap gap-4">
              <a
                href={data.resumePdf}
                download
                className="inline-flex items-center gap-2 px-6 py-3 text-sm font-semibold text-primary-foreground bg-primary hover:bg-primary/90 transition-all rounded-sm shadow-sm hover:shadow-md"
              >
                <Download className="w-5 h-5" />
                <span>下载完整简历 PDF</span>
              </a>
              <a
                href={`mailto:${data.email}`}
                className="inline-flex items-center gap-2 px-6 py-3 text-sm font-semibold text-foreground border-2 border-primary hover:bg-primary/5 transition-all rounded-sm"
              >
                <Mail className="w-5 h-5" />
                <span>立即联系</span>
              </a>
            </div>
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}
