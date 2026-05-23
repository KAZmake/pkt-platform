'use client';

import { useState, useMemo } from 'react';
import type { LoanProgram } from '@pkt/shared';
import { ProgramCard } from './program-card';
import { formatMoney } from '@/lib/utils';

type ActivityFilter = 'all' | 'crop_farming' | 'livestock' | 'mixed';

const ACTIVITY_OPTIONS: { value: ActivityFilter; label: string }[] = [
  { value: 'all', label: 'Все направления' },
  { value: 'crop_farming', label: 'Растениеводство' },
  { value: 'livestock', label: 'Животноводство' },
  { value: 'mixed', label: 'Смешанное' },
];

export function ProgramsClient({ programs }: { programs: LoanProgram[] }) {
  const [activity, setActivity] = useState<ActivityFilter>('all');
  const [compareIds, setCompareIds] = useState<string[]>([]);

  const filtered = useMemo(
    () =>
      programs.filter(
        (p) =>
          activity === 'all' ||
          (p.activityTypes ?? []).includes(activity as LoanProgram['activityTypes'][0]),
      ),
    [programs, activity],
  );

  const comparePrograms = programs.filter((p) => compareIds.includes(p.id));

  function toggleCompare(id: string) {
    setCompareIds((prev) =>
      prev.includes(id) ? prev.filter((x) => x !== id) : prev.length < 3 ? [...prev, id] : prev,
    );
  }

  return (
    <>
      {/* Filter tabs */}
      <div className="flex flex-wrap gap-2 mb-6">
        {ACTIVITY_OPTIONS.map((opt) => (
          <button
            key={opt.value}
            onClick={() => setActivity(opt.value)}
            className={`px-4 py-1.5 rounded-full text-sm font-medium transition-colors ${
              activity === opt.value
                ? 'bg-brand-green text-white'
                : 'bg-gray-100 text-gray-600 hover:bg-brand-green-100'
            }`}
          >
            {opt.label}
          </button>
        ))}
      </div>

      {/* Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {filtered.map((p) => (
          <ProgramCard
            key={p.id}
            program={p}
            selected={compareIds.includes(p.id)}
            onCompare={toggleCompare}
          />
        ))}
        {filtered.length === 0 && (
          <p className="col-span-3 text-gray-400 text-sm py-8 text-center">
            По выбранному фильтру программ не найдено.
          </p>
        )}
      </div>

      {/* Comparison table */}
      {comparePrograms.length >= 2 && (
        <div className="mt-12 overflow-x-auto rounded-xl border border-gray-200">
          <table className="min-w-full text-sm">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-xs text-gray-500 font-medium w-32">
                  Параметр
                </th>
                {comparePrograms.map((p) => (
                  <th
                    key={p.id}
                    className="px-4 py-3 text-left text-xs font-semibold text-gray-900"
                  >
                    {p.name}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {[
                { label: 'Ставка', fn: (p: LoanProgram) => `${p.rate}% год.` },
                {
                  label: 'Срок',
                  fn: (p: LoanProgram) => `${p.minTermMonths}–${p.maxTermMonths} мес.`,
                },
                { label: 'Мин. сумма', fn: (p: LoanProgram) => formatMoney(p.minAmount) },
                { label: 'Макс. сумма', fn: (p: LoanProgram) => formatMoney(p.maxAmount) },
              ].map(({ label, fn }) => (
                <tr key={label}>
                  <td className="px-4 py-3 text-xs text-gray-500 font-medium bg-gray-50">
                    {label}
                  </td>
                  {comparePrograms.map((p) => (
                    <td key={p.id} className="px-4 py-3 font-medium text-gray-900">
                      {fn(p)}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
          <div className="px-4 py-3 border-t border-gray-100 text-right">
            <button
              onClick={() => setCompareIds([])}
              className="text-xs text-gray-400 hover:text-red-500 transition-colors"
            >
              Очистить сравнение
            </button>
          </div>
        </div>
      )}
    </>
  );
}
