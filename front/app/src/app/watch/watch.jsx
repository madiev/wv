import { Suspense } from 'react';
import WatchComponent from './watch';

export const metadata = {
  title: "Watch Video",
  description: "Просмотр видео",
};

export default function Watch() {
  return (
    <Suspense>
      <WatchComponent />
    </Suspense>
  );
}
