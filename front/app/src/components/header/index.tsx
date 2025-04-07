import Search from '@/components/search';
import Image from 'next/image';
import Link from 'next/link';
import './style.css';

export default function Header() {
    return (
        <div className="grid auto-rows-[minmax(1fr,auto)] h-18 relative z-10">
            <div className="fixed top-0 right-0 left-0">
                <nav className="grid grid-cols-[1fr_2fr_1fr] gap-2 bg-purple-800 h-18">
                    <Link href="/" className="justify-self-center flex items-center space-x-3 rtl:space-x-reverse ml-4 mr-10">
                        <Image width={32} height={32} src="/video.svg" className="h-8" alt="Logo" />
                        <span className="rainbow-animated self-center text-2xl font-semibold whitespace-nowrap dark:text-white">WatchVideo</span>
                    </Link>
                    <div className="self-center">
                        <Search />
                    </div>
                    <div className="flex justify-center items-center">
                        {/*<div className="mr-2 text-white font-semibold">Инфо</div>
                        <svg className="w-[35px] h-[35px] text-white dark:text-white mr" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                            <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 11h2v5m-2 0h4m-2.592-8.5h.01M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
                        </svg>*/}
                    </div>
                </nav>
            </div>
        </div>
    )
}
