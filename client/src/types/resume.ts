/**
 * Resume 数据类型定义
 */

export interface Experience {
    title: string;
    company: string;
    period: string;
    highlights: string[];
}

export interface Education {
    degree: string;
    school: string;
    type: string;
    period: string;
    description: string;
}

export interface Skill {
    category: string;
    items: string[];
}

export interface Project {
    title: string;
    highlights: string[];
}

export interface ResumeData {
    name: string;
    title: string;
    subtitle: string;
    email: string;
    resumePdf: string;
    experience: Experience[];
    education: Education[];
    skills: Skill[];
    projects: Project[];
    traits: string[];
}
