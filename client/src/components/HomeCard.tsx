/**
 * HomeCard 组件
 * 首页的可复用卡片组件
 */

import { Link } from "wouter";
import type { LucideIcon } from "lucide-react";

interface HomeCardProps {
    href: string;
    icon: LucideIcon;
    title: string;
    description: string;
}

export default function HomeCard({ href, icon: Icon, title, description }: HomeCardProps) {
    return (
        <Link href={href} className="block h-full">
            <div className="group h-full p-8 rounded-2xl bg-card border border-border/40 hover:border-border/60 hover:bg-card/80 hover:shadow-lg hover:-translate-y-1 transition-all duration-300 cursor-pointer flex flex-col items-center text-center">
                <div className="h-11 w-11 bg-primary/5 rounded-2xl flex items-center justify-center mb-4 text-primary group-hover:scale-110 transition-transform duration-300">
                    <Icon className="w-6 h-6" strokeWidth={1.5} />
                </div>
                <h3 className="text-lg font-semibold tracking-tight mb-2.5 group-hover:text-primary transition-colors">
                    {title}
                </h3>
                <p className="text-sm text-muted-foreground leading-relaxed">{description}</p>
            </div>
        </Link>
    );
}
