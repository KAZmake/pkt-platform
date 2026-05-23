import type { Metadata } from 'next';
import Link from 'next/link';
import { notFound } from 'next/navigation';
import { getProgram, getPrograms } from '@/lib/api/programs';
import { Button } from '@/components/ui/button';
import { formatMoney } from '@/lib/utils';

interface Props {
  params: { id: string };
}

export async function generateStaticParams() {
  const programs = await getPrograms().catch(() => []);
  return programs.map((p) => ({ id: p.id }));
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const program = await getProgram(params.id).catch(() => null);
  return { title: program?.name ?? 'Программа кредитования' };
}

export default async function ProgramDetailPage({ params }: Props) {
  const program = await getProgram(params.id).catch(() => null);
  if (!program) notFound();

  const ACTIVITY_LABELS: Record<string, string> = {
    crop_farming: 'Растениеводство',
    livestock: 'Животноводство',
    mixed: 'Смешанное хозяйство',
  };

  return (
    <div className="container-page py-12">
      <nav className="text-sm text-gray-400 mb-6 flex items-center gap-2">
        <Link href="/programs" className="hover:text-brand-green">
          Программы
        </Link>
        <span>/</span>
        <span className="text-gray-700">{program.name}</span>
      </nav>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main info */}
        <div className="lg:col-span-2 space-y-6">
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-3">{program.name}</h1>
            {(program.activityTypes ?? []).length > 0 && (
              <div className="flex flex-wrap gap-2">
                {program.activityTypes.map((a) => (
                  <span
                    key={a}
                    className="text-xs bg-brand-green-100 text-brand-green-700 rounded-full px-3 py-1 font-medium"
                  >
                    {ACTIVITY_LABELS[a] ?? a}
                  </span>
                ))}
              </div>
            )}
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
            {[
              { label: 'Ставка', value: `${program.rate}%`, sub: 'в год' },
              {
                label: 'Срок',
                value: `${program.maxTermMonths} мес.`,
                sub: `от ${program.minTermMonths} мес.`,
              },
              { label: 'Мин. сумма', value: formatMoney(program.minAmount), sub: '' },
              { label: 'Макс. сумма', value: formatMoney(program.maxAmount), sub: '' },
            ].map(({ label, value, sub }) => (
              <div
                key={label}
                className="rounded-xl border border-gray-200 bg-white p-4 text-center"
              >
                <p className="text-xs text-gray-400 mb-1">{label}</p>
                <p className="text-xl font-bold text-brand-green">{value}</p>
                {sub && <p className="text-xs text-gray-400 mt-0.5">{sub}</p>}
              </div>
            ))}
          </div>
        </div>

        {/* Sidebar CTA */}
        <div className="rounded-xl border border-brand-green-100 bg-brand-green-50 p-6 flex flex-col gap-4 h-fit">
          <h2 className="font-semibold text-gray-900">Готовы подать заявку?</h2>
          <p className="text-sm text-gray-500 leading-relaxed">
            Заполните заявку онлайн — рассматриваем в течение 3 рабочих дней.
          </p>
          <Link href={`/apply?program=${program.id}`}>
            <Button className="w-full">Подать заявку</Button>
          </Link>
          <Link href="/calculator">
            <Button variant="outline" className="w-full">
              Калькулятор займа
            </Button>
          </Link>
          <p className="text-xs text-gray-400 text-center">
            Или позвоните:{' '}
            <a href="tel:+77112000000" className="text-brand-green font-medium">
              8 (7112) 00-00-00
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}
