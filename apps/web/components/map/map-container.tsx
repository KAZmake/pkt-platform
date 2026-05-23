/**
 * Dynamic import wrapper — MapLibre does not support SSR.
 * Renders a skeleton until the client-side component is ready.
 */
import dynamic from 'next/dynamic';
import { Spinner } from '@/components/ui/spinner';

const MapView = dynamic(() => import('./map-view').then((m) => ({ default: m.MapView })), {
  ssr: false,
  loading: () => (
    <div className="h-full w-full flex items-center justify-center bg-gray-100 rounded-xl">
      <Spinner size="lg" />
    </div>
  ),
});

export { MapView };
