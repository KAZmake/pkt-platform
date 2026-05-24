import Link from 'next/link';
import type { LoanProgram } from '@pkt/shared';
import { Button } from '@/components/ui/button';
import { formatMoney } from '@/lib/utils';

function ProgramCard({ program }: { program: LoanProgram }) {
  return (
    <div className="flex flex-col rounded-xl border border-gray-200 bg-white p-6 shadow-card hover:shadow-card-hover transition-shadow">
      <h3 className="font-semibold text-gray-900 mb-2">{program.name}</h3>
      <div className="mt-auto pt-4 grid grid-cols-2 gap-3 text-sm">
        <div>
          <p className="text-xs text-gray-500">Ставка</p>
          <p className="font-semibold text-brand-green">{program.rate}% год.</p>
        </div>
        <div>
          <p className="text-xs text-gray-500">Срок</p>
          <p className="font-semibold text-gray-900">до {program.maxTermMonths} мес.</p>
        </div>
        <div className="col-span-2">
          <p className="text-xs text-gray-500">Макс. сумма</p>
          <p className="font-semibold text-gray-900">{formatMoney(program.maxAmount)}</p>
        </div>
      </div>
      <Link href={`/programs/${program.id}`} className="mt-4">
        <Button variant="outline" size="sm" className="w-full">
          Подробнее
        </Button>
      </Link>
    </div>
  );
}

export function ProgramsTeaser({ programs }: { programs: LoanProgram[] }) {
  const featured = (Array.isArray(programs) ? programs : []).slice(0, 3);
  return (
    <section className="py-16 bg-white">
      <div className="container-page">
        <div className="flex items-end justify-between mb-8">
          <div>
            <p className="text-brand-gold text-xs font-semibold uppercase tracking-widest mb-2">
              Финансирование
            </p>
            <h2 className="text-2xl sm:text-3xl font-bold text-gray-900">Программы кредитования</h2>
          </div>
          <Link href="/programs" className="hidden sm:block">
            <Button variant="ghost" size="sm">
              Все программы →
            </Button>
          </Link>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {featured.map((p) => (
            <ProgramCard key={p.id} program={p} />
          ))}
        </div>

        <div className="mt-6 sm:hidden text-center">
          <Link href="/programs">
            <Button variant="outline">Все программы</Button>
          </Link>
        </div>
      </div>
    </section>
  );
}
