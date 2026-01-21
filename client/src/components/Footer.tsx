/**
 * Footer 组件
 * 极简现代主义风格的页脚
 */

import { siteConfig } from "@/config/site";
import { Github, Mail, BookOpen, MessageCircle, Eye } from "lucide-react";
import { useEffect } from "react";

const socialLinks = [
  { href: siteConfig.social.github, icon: Github, label: "GitHub" },
  { href: siteConfig.social.bilibili, icon: BookOpen, label: "Bilibili" },
  { href: siteConfig.social.zhihu, icon: MessageCircle, label: "Zhihu" },
  { href: `mailto:${siteConfig.email}`, icon: Mail, label: "Email" },
];

export default function Footer() {
  const currentYear = new Date().getFullYear();

  // 加载不蒜子计数器脚本
  useEffect(() => {
    const script = document.createElement("script");
    script.src = "//busuanzi.ibruce.info/busuanzi/2.3/busuanzi.pure.mini.js";
    script.async = true;
    document.body.appendChild(script);
    return () => {
      document.body.removeChild(script);
    };
  }, []);

  return (
    <footer className="border-t border-border mt-16 py-12 bg-muted/30">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* 社交媒体链接 */}
        <div className="flex items-center justify-center gap-8 mb-8">
          {socialLinks.map(({ href, icon: Icon, label }) => (
            <a
              key={label}
              href={href}
              target={href.startsWith("mailto:") ? undefined : "_blank"}
              rel={href.startsWith("mailto:") ? undefined : "noopener noreferrer"}
              className="p-2 text-muted-foreground hover:text-foreground hover:bg-background transition-all rounded-sm"
              aria-label={label}
            >
              <Icon className="w-5 h-5" />
            </a>
          ))}
        </div>

        {/* 访客计数 */}
        <div className="flex items-center justify-center gap-2 mb-6 text-sm text-muted-foreground">
          <Eye className="w-4 h-4" />
          <span>本站已被访问</span>
          <span id="busuanzi_value_site_pv" className="font-semibold text-foreground">--</span>
          <span>次</span>
        </div>

        {/* 分割线 */}
        <div className="border-t border-border my-8"></div>

        {/* 版权信息 */}
        <div className="text-center space-y-2">
          <p className="text-sm text-muted-foreground">
            ©{currentYear} <span className="font-semibold">{siteConfig.name}</span>。保留所有权利。
          </p>
          <p className="text-xs text-muted-foreground">
            Built with React, TypeScript & Tailwind CSS
          </p>
        </div>
      </div>
    </footer>
  );
}
