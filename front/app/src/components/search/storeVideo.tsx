import { create } from 'zustand';

type VideoStoreState = {
    data: {
        video: ElementsVideo[];
    };
    pending: boolean;
    error: Error | null;
}

export type ElementsVideo = {
    id: string;
    thumbnail: string;
    title: string;
    description: string;
    duration: string;
}

type VideoStoreActions = {
  fetchSearch: (query: string) => void;
}

type VideoStore = VideoStoreState & VideoStoreActions;

export const useStoreVideo = create<VideoStore>((set) => ({
  data: {
    video: [],
  },
  pending: false,
  error: null,
  fetchSearch: async (query: string) => {
    const params = new URLSearchParams({
        q: query, limit: '20'
    }).toString();
    set({ pending: true });
    try {
        const data = await fetch(`/api/search?${params}`);
        const res = await data.json();    
        set({data: res, error: null});
    } catch (error) {
        if (error instanceof Error) {
            set({error}); 
        }
    } finally {
        set({pending: false});
    }
  },
}));