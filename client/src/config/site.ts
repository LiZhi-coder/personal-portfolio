/**
 * 站点配置
 * 统一管理站点信息，避免硬编码散落各处
 */

export const siteConfig = {
    // 个人信息
    name: "呵呜微",
    title: "Backend Developer · Tech Writer",
    siteName: "Hui Studio",
    email: "lizhihui916@163.com",
    avatar: "/avatar1.png",

    // 社交链接
    social: {
        github: "https://github.com/hui-cyber/",
        bilibili: "https://www.bilibili.com/",
        zhihu: "https://www.zhihu.com/",
    },

    // 简历下载
    resumePdf: "/resume.pdf",

    // SEO
    description: "欢迎来到我的个人主页，这里展示了我的项目作品、技术博客和个人简历。",
} as const;

export type SiteConfig = typeof siteConfig;
