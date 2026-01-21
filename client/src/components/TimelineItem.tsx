/**
 * TimelineItem 组件
 * 用于 Resume 页面的时间线项目
 */

import type { ReactNode } from "react";

interface TimelineItemProps {
    title: string;
    subtitle?: string;
    children: ReactNode;
    isLast?: boolean;
}

export default function TimelineItem({ title, subtitle, children, isLast = false }: TimelineItemProps) {
    return (
        <div className={`relative border-l-2 border-primary pl-8 ${!isLast ? 'pb-8' : ''} hover:border-primary/70 transition-colors`}>
            <div className="absolute w-4 h-4 bg-primary rounded-full -left-[9px] top-0"></div>
            <div className="space-y-2">
                <h4 className="text-xl font-semibold text-foreground">{title}</h4>
                {subtitle && (
                    <p className="text-sm text-muted-foreground font-medium">{subtitle}</p>
                )}
                {children}
            </div>
        </div>
    );
}
