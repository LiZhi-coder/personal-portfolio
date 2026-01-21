/**
 * Navigation 组件
 * 极简现代主义风格的导航栏
 */

import { siteConfig } from "@/config/site";
import { Menu, X } from "lucide-react";
import { useState } from "react";
import { Link, useLocation } from "wouter";

interface NavItem {
  label: string;
  href: string;
}

const navItems: NavItem[] = [
  { label: "首页", href: "/" },
  { label: "项目集", href: "/projects" },
  { label: "博客", href: "/blog" },
  { label: "简历", href: "/resume" },
];

export default function Navigation() {
  const [isOpen, setIsOpen] = useState(false);
  const [location] = useLocation();

  return (
    <nav className="sticky top-0 z-50 bg-background/80 backdrop-blur-md border-b border-border shadow-sm">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo / 名字 */}
          <Link
            href="/"
            className="text-xl font-bold text-foreground hover:text-primary transition-colors"
          >
            {siteConfig.siteName}
          </Link>

          {/* 桌面导航 */}
          <div className="hidden md:flex items-center gap-8">
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={`text-base transition-colors ${location === item.href
                    ? "text-primary font-semibold"
                    : "text-foreground hover:text-primary"
                  }`}
              >
                {item.label}
              </Link>
            ))}
          </div>

          {/* 移动菜单按钮 */}
          <button
            onClick={() => setIsOpen(!isOpen)}
            className="md:hidden text-foreground hover:text-primary transition-colors"
            aria-label="切换菜单"
            aria-expanded={isOpen}
            aria-controls="mobile-nav"
          >
            {isOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>

        {/* 移动导航菜单 */}
        {isOpen && (
          <div id="mobile-nav" className="md:hidden pb-4 space-y-2">
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={`block px-4 py-2 text-base rounded-sm transition-colors ${location === item.href
                    ? "text-primary bg-primary/10 font-semibold"
                    : "text-foreground hover:text-primary hover:bg-card"
                  }`}
                onClick={() => setIsOpen(false)}
              >
                {item.label}
              </Link>
            ))}
          </div>
        )}
      </div>
    </nav>
  );
}
