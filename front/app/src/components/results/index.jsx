'use client';

import { useStoreVideo } from '@/components/search';
import Image from 'next/image';
import { useRouter } from 'next/navigation';

export default function VideoComponent() {
    const { video } = useStoreVideo((state) => state.data);
    const router = useRouter();

    const handlerClick = (el) => {
        router.push(`/watch?v=${el.id}`);
    }

    return (
        <div className="grid justify-center" >
            {Array.isArray(video) && video.length > 0 && video.map((el, i) => (
                <div key={i}>
                    <div className="grid grid-cols-2 break-words">
                        <div className="pt-5 pb-5 pr-5 flex justify-center">
                            <Image
                                width={380}
                                height={360}
                                src={el.thumbnail}
                                alt={el.title}
                                onClick={() => handlerClick(el)}
                                className="cursor-pointer border border-gray-700 border-1 rounded-3xl w-full"
                            />
                        </div>
                        <div className="max-w-[500px] pt-5 pb-5">
                            <div className="text-white mb-2 text-xl font-semibold"  dangerouslySetInnerHTML={{ __html: el.title }}></div>
                            <div className="text-gray-300 md:text-sm dark:text-gray-400">{el.description}</div>
                        </div>
                        <div className="border-b border-purple-400 col-start-1 col-end-3"></div>
                    </div>
                </div>
            ))}
        </div>
    );
};