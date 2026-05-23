'use client';

import { useState, useMemo } from 'react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { formatMoney } from '@/lib/utils';

type PaymentType = 'annuity' | 'differentiated';

interface ScheduleRow {
  period: number;
  principal: number;
  interest: number;
  total: number;
  balance: number;
}

function calcAnnuity(amount: number, termMonths: number, annualRate: number): ScheduleRow[] {
  const r = annualRate / 100 / 12;
  if (r === 0) {
    const p = amount / termMonths;
    return Array.from({ length: termMonths }, (_, i) => ({
      period: i + 1,
      principal: p,
      interest: 0,
      total: p,
      balance: amount - p * (i + 1),
    }));
  }
  const payment = (amount * r * Math.pow(1 + r, termMonths)) / (Math.pow(1 + r, termMonths) - 1);
  const rows: ScheduleRow[] = [];
  let balance = amount;
  for (let i = 1; i <= termMonths; i++) {
    const interest = balance * r;
    const principal = payment - interest;
    balance -= principal;
    rows.push({ period: i, principal, interest, total: payment, balance: Math.max(0, balance) });
  }
  return rows;
}

function calcDifferentiated(amount: number, termMonths: number, annualRate: number): ScheduleRow[] {
  const r = annualRate / 100 / 12;
  const principalPayment = amount / termMonths;
  const rows: ScheduleRow[] = [];
  let balance = amount;
  for (let i = 1; i <= termMonths; i++) {
    const interest = balance * r;
    const total = principalPayment + interest;
    balance -= principalPayment;
    rows.push({
      period: i,
      principal: principalPayment,
      interest,
      total,
      balance: Math.max(0, balance),
    });
  }
  return rows;
}

export function LoanCalculator() {
  const [amount, setAmount] = useState('5000000');
  const [term, setTerm] = useState('24');
  const [rate, setRate] = useState('12');
  const [type, setType] = useState<PaymentType>('annuity');
  const [showFull, setShowFull] = useState(false);

  const schedule = useMemo(() => {
    const a = parseFloat(amount) || 0;
    const t = parseInt(term) || 0;
    const r = parseFloat(rate) || 0;
    if (a <= 0 || t <= 0 || r < 0) return [];
    return type === 'annuity' ? calcAnnuity(a, t, r) : calcDifferentiated(a, t, r);
  }, [amount, term, rate, type]);

  const totalInterest = schedule.reduce((s, r) => s + r.interest, 0);
  const totalPayment = schedule.reduce((s, r) => s + r.total, 0);
  const firstPayment = schedule[0]?.total ?? 0;

  const displayRows = showFull ? schedule : schedule.slice(0, 6);

  return (
    <div className="flex flex-col gap-8">
      {/* Inputs */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Input
          label="Сумма займа (₸)"
          type="number"
          min={100000}
          max={500000000}
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          hint="от 100 000 до 500 000 000 ₸"
        />
        <Input
          label="Срок (месяцев)"
          type="number"
          min={1}
          max={120}
          value={term}
          onChange={(e) => setTerm(e.target.value)}
          hint="от 1 до 120 месяцев"
        />
        <Input
          label="Ставка (% год.)"
          type="number"
          min={0}
          max={100}
          step={0.1}
          value={rate}
          onChange={(e) => setRate(e.target.value)}
          hint="годовая процентная ставка"
        />
        <div className="flex flex-col gap-1">
          <label className="text-sm font-medium text-gray-700">Тип платежей</label>
          <div className="flex rounded-lg border border-gray-300 overflow-hidden h-10">
            {(['annuity', 'differentiated'] as const).map((t) => (
              <button
                key={t}
                onClick={() => setType(t)}
                className={`flex-1 text-xs font-medium transition-colors ${
                  type === t
                    ? 'bg-brand-green text-white'
                    : 'bg-white text-gray-600 hover:bg-gray-50'
                }`}
              >
                {t === 'annuity' ? 'Аннуитетный' : 'Дифференц.'}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Summary */}
      {schedule.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          {[
            {
              label: type === 'annuity' ? 'Ежемес. платёж' : 'Первый платёж',
              value: formatMoney(firstPayment),
            },
            { label: 'Переплата (проценты)', value: formatMoney(totalInterest) },
            { label: 'Итого к выплате', value: formatMoney(totalPayment) },
          ].map(({ label, value }) => (
            <div
              key={label}
              className="rounded-xl bg-brand-green-50 border border-brand-green-100 p-4"
            >
              <p className="text-xs text-gray-500 mb-1">{label}</p>
              <p className="text-xl font-bold text-brand-green">{value}</p>
            </div>
          ))}
        </div>
      )}

      {/* Schedule table */}
      {schedule.length > 0 && (
        <div className="overflow-x-auto rounded-xl border border-gray-200">
          <table className="min-w-full text-sm">
            <thead className="bg-gray-50 text-xs text-gray-500 uppercase tracking-wide">
              <tr>
                {['№', 'Основной долг', 'Проценты', 'Платёж', 'Остаток'].map((h) => (
                  <th key={h} className="px-4 py-3 text-right first:text-left font-medium">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {displayRows.map((row) => (
                <tr key={row.period} className="hover:bg-gray-50">
                  <td className="px-4 py-2.5 text-gray-500">{row.period}</td>
                  <td className="px-4 py-2.5 text-right">{formatMoney(row.principal)}</td>
                  <td className="px-4 py-2.5 text-right text-brand-gold">
                    {formatMoney(row.interest)}
                  </td>
                  <td className="px-4 py-2.5 text-right font-medium">{formatMoney(row.total)}</td>
                  <td className="px-4 py-2.5 text-right text-gray-500">
                    {formatMoney(row.balance)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {schedule.length > 6 && (
            <div className="px-4 py-3 border-t border-gray-100 text-center">
              <button
                onClick={() => setShowFull(!showFull)}
                className="text-sm text-brand-green hover:underline"
              >
                {showFull ? 'Скрыть' : `Показать все ${schedule.length} платежей`}
              </button>
            </div>
          )}
        </div>
      )}

      {/* CTA */}
      <div className="flex flex-col sm:flex-row gap-3 pt-2">
        <a href="/apply">
          <Button size="lg">Подать заявку с этими параметрами</Button>
        </a>
        <a href="/programs">
          <Button variant="outline" size="lg">
            Смотреть программы
          </Button>
        </a>
      </div>
    </div>
  );
}
