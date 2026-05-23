import type { Metadata } from 'next';
import { auth } from '@/auth';
import { MapView } from '@/components/map/map-container';

export const metadata: Metadata = {
  title: 'Интерактивная карта ЗКО',
  description:
    'Карта хозяйств Западно-Казахстанской области. Районы, кластеры предприятий, зоны охвата ПКТ.',
};

export default async function MapPage() {
  const session = await auth();
  const role = (session as { role?: string })?.role ?? 'public';

  return (
    <div className="flex flex-col" style={{ height: 'calc(100vh - 112px)' }}>
      {/* Page header */}
      <div className="container-page py-4 shrink-0">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Карта ЗКО</h1>
            <p className="text-xs text-gray-500 mt-0.5">
              Хозяйства, районы, зона охвата ПКТ
              {['employee', 'expert', 'admin'].includes(role) && ' · Режим сотрудника'}
            </p>
          </div>
          <div className="flex items-center gap-2 text-xs text-gray-400">
            <span className="h-2 w-2 rounded-full bg-brand-green" /> Растениеводство
            <span className="h-2 w-2 rounded-full bg-brand-gold ml-2" /> Животноводство
            <span className="h-2 w-2 rounded-full bg-blue-500 ml-2" /> Смешанное
          </div>
        </div>
      </div>

      {/* Full-height map */}
      <div className="flex-1 container-page pb-4">
        <div className="h-full rounded-xl overflow-hidden border border-gray-200 shadow-card">
          <MapView role={role} />
        </div>
      </div>
    </div>
  );
}
