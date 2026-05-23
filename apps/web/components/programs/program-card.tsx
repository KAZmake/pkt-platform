import Link from 'next/link';
import type { LoanProgram } from '@pkt/shared';
import { Button } from '@/components/ui/button';
import { formatMoney } from '@/lib/utils';

const ACTIVITY_LABELS: Record<string, string> = {
  crop_farming: 'Растениеводство',
  livestock: 'Животноводство',
  mixed: 'Смешанное',
};

interface ProgramCardProps {
  program: LoanProgram;
  selected?: boolean;
  onCompare?: (id: string) => void;
}

export function ProgramCard({ program, selected, onCompare }: ProgramCardProps) {
  return (
    <div className="flex flex-col rounded-xl border border-gray-200 bg-white shadow-card hover:shadow-card-hover transition-shadow">
      <div className="p-6 flex-1">
        <div className="flex items-start justify-between gap-3 mb-4">
          <h3 className="font-semibold text-gray-900 leading-snug">{program.name}</h3>
          {program.activityTypes?.length > 0 && (
            <span className="shrink-0 text-xs bg-brand-green-100 text-brand-green-700 rounded-full px-2 py-0.5">
              {ACTIVITY_LABELS[program.activityTypes[0]] ?? program.activityTypes[0]}
            </span>
          )}
        </div>

        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <p className="text-xs text-gray-400 mb-0.5">Ставка</p>
            <p className="font-bold text-brand-green text-lg">{program.rate}%</p>
            <p className="text-xs text-gray-400">в год</p>
          </div>
          <div>
            <p className="text-xs text-gray-400 mb-0.5">Срок</p>
            <p className="font-semibold text-gray-900">
              {program.minTermMonths}–{program.maxTermMonths}
            </p>
            <p className="text-xs text-gray-400">месяцев</p>
          </div>
          <div>
            <p className="text-xs text-gray-400 mb-0.5">Мин. сумма</p>
            <p className="font-semibold text-gray-900">{formatMoney(program.minAmount)}</p>
          </div>
          <div>
            <p className="text-xs text-gray-400 mb-0.5">Макс. сумма</p>
            <p className="font-semibold text-gray-900">{formatMoney(program.maxAmount)}</p>
          </div>
        </div>
      </div>

      <div className="px-6 pb-6 flex gap-2">
        <Link href={`/programs/${program.id}`} className="flex-1">
          <Button variant="primary" size="sm" className="w-full">
            Подробнее
          </Button>
        </Link>
        {onCompare && (
          <Button
            variant={selected ? 'secondary' : 'outline'}
            size="sm"
            onClick={() => onCompare(program.id)}
            title={selected ? 'Убрать из сравнения' : 'Добавить к сравнению'}
          >
            {selected ? '✓' : '⇄'}
          </Button>
        )}
      </div>
    </div>
  );
}
