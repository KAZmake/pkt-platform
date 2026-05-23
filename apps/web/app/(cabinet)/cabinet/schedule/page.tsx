import { ScheduleChart } from '@/components/cabinet/schedule-chart';
import type { Metadata } from 'next';

export const metadata: Metadata = { title: 'График платежей' };

interface ScheduleRow {
  month: number;
  date: string;
  principal: number;
  interest: number;
  total: number;
  balance: number;
  status: 'paid' | 'upcoming' | 'overdue';
}

function buildSchedule(): ScheduleRow[] {
  const loanAmount = 15_000_000;
  const annualRate = 0.12;
  const months = 48;
  const r = annualRate / 12;
  const annuity = (loanAmount * r * Math.pow(1 + r, months)) / (Math.pow(1 + r, months) - 1);

  let balance = loanAmount;
  const rows: ScheduleRow[] = [];
  const start = new Date('2024-06-01');
  const today = new Date('2026-05-24');

  for (let i = 1; i <= months; i++) {
    const interestPart = balance * r;
    const principalPart = annuity - interestPart;
    balance = Math.max(0, balance - principalPart);

    const d = new Date(start);
    d.setMonth(d.getMonth() + i);

    const status: ScheduleRow['status'] = d <= today ? 'paid' : 'upcoming';

    rows.push({
      month: i,
      date: d.toISOString().slice(0, 10),
      principal: Math.round(principalPart),
      interest: Math.round(interestPart),
      total: Math.round(annuity),
      balance: Math.round(balance),
      status,
    });
  }
  return rows;
}

function fmt(n: number) {
  return n.toLocaleString('ru-KZ') + ' ₸';
}

export default function SchedulePage() {
  const rows = buildSchedule();
  const paid = rows.filter((r) => r.status === 'paid');
  const totalPaid = paid.reduce((s, r) => s + r.total, 0);
  const totalInterestPaid = paid.reduce((s, r) => s + r.interest, 0);
  const totalPrincipalPaid = totalPaid - totalInterestPaid;

  return (
    <div className="max-w-5xl">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">График платежей</h1>
        <p className="text-sm text-gray-500 mt-1">
          Займ LN-2024-0042 · 48 месяцев · 12% годовых · аннуитетные платежи
        </p>
      </div>

      {/* Summary */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <SummaryCard
          label="Выплачено платежей"
          value={`${paid.length} из ${rows.length}`}
          sub="месяцев"
        />
        <SummaryCard label="Погашено основного долга" value={fmt(totalPrincipalPaid)} sub="" />
        <SummaryCard label="Уплачено процентов" value={fmt(totalInterestPaid)} sub="" />
        <SummaryCard
          label="Осталось платежей"
          value={`${rows.length - paid.length} из ${rows.length}`}
          sub="месяцев"
        />
      </div>

      {/* Chart — first 24 months */}
      <div className="bg-white rounded-xl border border-gray-200 p-5 mb-6">
        <p className="text-sm font-medium text-gray-700 mb-4">
          Структура платежей (первые 24 мес.)
        </p>
        <ScheduleChart rows={rows.slice(0, 24)} />
        <div className="flex gap-5 mt-3 text-xs text-gray-500">
          <span className="flex items-center gap-1.5">
            <span className="w-3 h-3 rounded-sm bg-brand-green inline-block" />
            Основной долг
          </span>
          <span className="flex items-center gap-1.5">
            <span className="w-3 h-3 rounded-sm bg-brand-gold inline-block" />
            Проценты
          </span>
          <span className="flex items-center gap-1.5">
            <span className="w-3 h-3 rounded-sm bg-brand-green-100 inline-block" />
            Предстоящие
          </span>
        </div>
      </div>

      {/* Full table */}
      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide">
                  №
                </th>
                <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide">
                  Дата
                </th>
                <th className="text-right px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide">
                  Осн. долг
                </th>
                <th className="text-right px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide">
                  Проценты
                </th>
                <th className="text-right px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide">
                  Итого
                </th>
                <th className="text-right px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide hidden md:table-cell">
                  Остаток
                </th>
                <th className="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide">
                  Статус
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {rows.map((row) => (
                <tr
                  key={row.month}
                  className={`hover:bg-gray-50 ${row.status === 'paid' ? 'text-gray-400' : 'text-gray-700'}`}
                >
                  <td className="px-4 py-2 text-xs">{row.month}</td>
                  <td className="px-4 py-2 whitespace-nowrap">
                    {new Date(row.date).toLocaleDateString('ru-KZ', {
                      day: 'numeric',
                      month: 'short',
                      year: 'numeric',
                    })}
                  </td>
                  <td className="px-4 py-2 text-right tabular-nums">
                    {row.principal.toLocaleString('ru-KZ')}
                  </td>
                  <td className="px-4 py-2 text-right tabular-nums">
                    {row.interest.toLocaleString('ru-KZ')}
                  </td>
                  <td className="px-4 py-2 text-right tabular-nums font-medium">
                    {row.total.toLocaleString('ru-KZ')}
                  </td>
                  <td className="px-4 py-2 text-right tabular-nums hidden md:table-cell">
                    {row.balance.toLocaleString('ru-KZ')}
                  </td>
                  <td className="px-4 py-2 text-center">
                    {row.status === 'paid' ? (
                      <span className="inline-flex items-center gap-1 text-xs text-emerald-600">
                        <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                        Оплачен
                      </span>
                    ) : row.status === 'overdue' ? (
                      <span className="inline-flex items-center gap-1 text-xs text-red-600">
                        <span className="w-1.5 h-1.5 rounded-full bg-red-500" />
                        Просрочен
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1 text-xs text-gray-400">
                        <span className="w-1.5 h-1.5 rounded-full bg-gray-300" />
                        Предстоит
                      </span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function SummaryCard({ label, value, sub }: { label: string; value: string; sub: string }) {
  return (
    <div className="bg-white rounded-xl border border-gray-200 p-4">
      <p className="text-xs text-gray-500 mb-1">{label}</p>
      <p className="text-base font-bold text-gray-900">{value}</p>
      {sub && <p className="text-xs text-gray-400 mt-0.5">{sub}</p>}
    </div>
  );
}
