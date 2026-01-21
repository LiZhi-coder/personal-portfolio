/**
 * useFetchData Hook
 * 通用数据获取 Hook，用于 Blog 和 Projects 页面
 */

import { useEffect, useState } from "react";

interface FetchState<T> {
    data: T[];
    loading: boolean;
    error: string | null;
}

export function useFetchData<T>(url: string, errorMessage: string = "加载失败"): FetchState<T> {
    const [data, setData] = useState<T[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                setError(null);

                const response = await fetch(url);
                if (!response.ok) {
                    throw new Error("Failed to load data");
                }

                const json: T[] = await response.json();
                setData(json);
            } catch (err) {
                console.error("Error loading data:", err);
                setError(errorMessage);
                setData([]);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [url, errorMessage]);

    return { data, loading, error };
}
