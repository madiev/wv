export default function Footer() {
    const currentDate = new Date();
    return (
        <footer className="h-10 bg-purple-400 row-start-3 flex gap-6 flex-wrap items-center justify-center">
            Поделись ссылкой на видео {currentDate.getFullYear()}
        </footer>
    );
};