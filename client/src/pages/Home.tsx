/**
 * Home 页面 - 极简主义落地页
 * 核心理念：Less is More
 */

import Footer from "@/components/Footer";
import HomeCard from "@/components/HomeCard";
import Navigation from "@/components/Navigation";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { siteConfig } from "@/config/site";
import { Archive, FileText, User } from "lucide-react";
import { Helmet } from "react-helmet-async";

const homeCards = [
  {
    href: "/projects",
    icon: Archive,
    title: "主要项目",
    description: "查看最新的开发成果。",
  },
  {
    href: "/blog",
    icon: FileText,
    title: "内容写作",
    description: "随手记录生活点滴。",
  },
  {
    href: "/resume",
    icon: User,
    title: "关于我",
    description: "工作经历、技能栈与联系方式。",
  },
];

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col bg-background font-sans selection:bg-primary/10">
      <Helmet>
        <title>{siteConfig.name}的个人主页</title>
        <meta name="description" content={siteConfig.description} />
      </Helmet>

      <Navigation />

      <main className="flex-1 px-6 sm:px-8">
        <div className="max-w-5xl mx-auto pt-10 sm:pt-12 pb-12 sm:pb-16">
          {/* 顶部信息条 */}
          <div
            className="flex items-center justify-between gap-4 animate-fade-in-up"
            style={{ animationDelay: "0.1s" }}
          >
            <div className="flex items-center gap-3 min-w-0">
              <Avatar className="w-9 h-9 border border-border/50">
                <AvatarImage src={siteConfig.avatar} alt={siteConfig.name} />
                <AvatarFallback className="text-xs bg-muted text-muted-foreground">
                  {siteConfig.name[0]}
                </AvatarFallback>
              </Avatar>
              <div className="min-w-0">
                <div className="text-base font-semibold tracking-tight text-foreground leading-tight">
                  {siteConfig.name}
                </div>
                <div className="text-sm text-muted-foreground/80 truncate">
                  {siteConfig.title}
                </div>
              </div>
            </div>

            <div className="hidden sm:block text-sm text-muted-foreground">
              欢迎您的到来！
            </div>
          </div>

          <div className="mt-5 border-t border-border/50" />

          {/* 内容卡片 - 三分天下 */}
          <div className="mt-7 grid grid-cols-1 md:grid-cols-3 gap-6">
            {homeCards.map((card) => (
              <HomeCard key={card.href} {...card} />
            ))}
          </div>
        </div>
      </main>

      <Footer />
    </div>
  );
}
