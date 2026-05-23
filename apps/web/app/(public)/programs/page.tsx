import type { Metadata } from 'next';
import { getPrograms } from '@/lib/api/programs';
import { ProgramsClient } from '@/components/programs/programs-client';

export const metadata: Metadata = {
  title: 'Программы кредитования',
  description:
    'Льготные займы для фермеров и агропредприятий ЗКО. Фильтр по направлению, сравнение условий.',
};

export default async function ProgramsPage() {
  const programs = await getPrograms().catch(() => []);

  return (
    <div className="container-page py-12">
      <div className="mb-8">
        <p className="text-brand-gold text-xs font-semibold uppercase tracking-widest mb-2">
          Финансирование
        </p>
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
          Программы кредитования
        </h1>
        <p className="text-gray-500 text-sm max-w-xl">
          Выберите подходящую программу. Нажмите «⇄» на карточке, чтобы добавить до 3 программ к
          сравнению.
        </p>
      </div>
      <ProgramsClient programs={programs} />
    </div>
  );
}
