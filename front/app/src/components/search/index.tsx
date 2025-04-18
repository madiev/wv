'use client';

import { useState, useEffect, useCallback, FormEvent } from 'react';
import { create } from 'zustand';
import { useRouter } from 'next/navigation';

type VideoStoreState = {
    data: {
        video: ElementsVideo[];
    };
    pending: boolean;
}

export type ElementsVideo = {
    id: string;
    thumbnail: string;
    title: string;
    description: string;
    duration: string;
}

type VideoStoreActions = {
  setVideo: ( res: VideoStoreState['data'] ) => void;
  setPending: ( bool: VideoStoreState['pending'] ) => void;
}

type VideoStore = VideoStoreState & VideoStoreActions;

export const useStoreVideo = create<VideoStore>((set) => ({
  data: {
    video: [],
  },
  pending: true,
  setVideo: (res) => set(() => ({ data: res })),
  setPending: (bool) => set(() => ({ pending: bool })),
}));

export default function Search() {
    const [q, setQ] = useState('');
    const router = useRouter();
    const { setVideo, setPending } = useStoreVideo((state) => state);
  
    const fetchSearch = useCallback(async (query: string) => {
        const params = new URLSearchParams({
            q: query, limit: '20'
        }).toString();
        setPending(true);
        const data = await fetch(`/api/search?${params}`);
        const res = await data.json();
        setVideo(res);
        setPending(false);
    }, [setVideo, setPending]);
    
    function handlerSubmit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault();
        fetchSearch(q);
        router.push('/');
    }

    useEffect(() => {
        fetchSearch('Музыка');
    }, [fetchSearch]);

    return (
        <form className="max-w-2xl mx-auto w-full pr-3" onSubmit={handlerSubmit} >
            <label htmlFor="default-search" className="mb-2 text-sm font-medium text-gray-900 sr-only dark:text-white">Искать</label>
            <div className="relative">
                <div className="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
                    <svg className="w-4 h-4 text-gray-500 dark:text-gray-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20">
                        <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"/>
                    </svg>
                </div>
                <input
                    type="search"
                    id="default-search"
                    className="block w-full p-4 ps-10 text-sm text-base text-black border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                    placeholder="Поиск видео..."
                    onChange={(e) => setQ(e.target.value)}
                />
                <button type="submit" className="text-white absolute end-2.5 bottom-2.5 bg-gradient-to-r from-orange-500 via-orange-600 to-orange-700 hover:bg-gradient-to-br focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Поиск</button>
            </div>
        </form>
    );
};