import { auth } from '@/auth';
import { ApplicationStatusBadge, Badge } from '@/components/ui/badge';
import Link from 'next/link';
import type { Metadata } from 'next';

export const metadata: Metadata = { title: 'Дашборд' };

// Mock summary — replace with: apiClient.get('/cabinet/summary', { token: session.accessToken })
const MOCK = {
  loan: {
    id: 'LN-2024-0042',
    program: 'Программа «Агробизнес 2025»',
    total: 15_000_000,
    remaining: 11_250_000,
    next_payment_date: '2026-06-01',
    next_payment_amount: 312_500,
    status: 'issued' as const,
    term_months: 48,
    paid_months: 12,
  },
  documents_count: 7,
  unread_notifications: 3,
  open_tickets: 1,
};

function fmt(n: number) {
  return n.toLocaleString('ru-KZ') + ' ₸';
}

function fmtDate(s: string) {
  return new Date(s).toLocaleDateString('ru-KZ', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
}

export default async function CabinetDashboardPage() {
  const session = await auth();
  const firstName = session?.user?.name?.split(' ')[0] ?? 'Заёмщик';
  const loan = MOCK.loan;
  const paidPercent = Math.round((loan.paid_months / loan.term_months) * 100);
  const paidAmount = loan.total - loan.remaining;

  return (
    <div className="max-w-5xl">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Добро пожаловать, {firstName}</h1>
        <p className="text-sm text-gray-500 mt-1">Личный кабинет заёмщика</p>
      </div>

      {/* Stat cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <StatCard label="Остаток долга" value={fmt(loan.remaining)} sub={`из ${fmt(loan.total)}`} />
        <StatCard
          label="Следующий платёж"
          value={fmt(loan.next_payment_amount)}
          sub={fmtDate(loan.next_payment_date)}
        />
        <StatCard label="Документов" value={String(MOCK.documents_count)} sub="в архиве" />
        <StatCard
          label="Уведомлений"
          value={String(MOCK.unread_notifications)}
          sub="непрочитанных"
          highlight={MOCK.unread_notifications > 0}
        />
      </div>

      {/* Loan card */}
      <div className="bg-white rounded-xl border border-gray-200 p-5 mb-6">
        <div className="flex items-start justify-between gap-4 mb-4">
          <div>
            <p className="text-xs text-gray-400 mb-0.5">Займ {loan.id}</p>
            <h2 className="text-base font-semibold text-gray-900">{loan.program}</h2>
          </div>
          <ApplicationStatusBadge status={loan.status} />
        </div>

        <div className="mb-1">
          <div className="flex justify-between text-xs text-gray-500 mb-1">
            <span>Выплачено: {fmt(paidAmount)}</span>
            <span>{paidPercent}%</span>
          </div>
          <div className="h-2 bg-gray-100 rounded-full overflow-hidden">
            <div
              className="h-full bg-brand-green rounded-full"
              style={{ width: `${paidPercent}%` }}
            />
          </div>
          <div className="flex justify-between text-xs text-gray-400 mt-1">
            <span>{loan.paid_months} мес. выплачено</span>
            <span>{loan.term_months - loan.paid_months} мес. осталось</span>
          </div>
        </div>

        <div className="flex gap-4 mt-4 pt-3 border-t border-gray-100">
          <Link
            href="/cabinet/schedule"
            className="text-sm text-brand-green font-medium hover:underline"
          >
            График платежей →
          </Link>
          <Link
            href="/cabinet/documents"
            className="text-sm text-brand-green font-medium hover:underline"
          >
            Документы →
          </Link>
        </div>
      </div>

      {/* Quick actions */}
      <div className="grid sm:grid-cols-2 gap-4">
        <QuickActionCard
          href="/cabinet/tickets"
          icon="💬"
          title="Создать обращение"
          desc="Задайте вопрос или сообщите о проблеме"
          badge={MOCK.open_tickets > 0 ? `${MOCK.open_tickets} открытых` : undefined}
        />
        <QuickActionCard
          href="/cabinet/notifications"
          icon="🔔"
          title="Уведомления"
          desc="Последние сообщения от кредитного товарищества"
          badge={MOCK.unread_notifications > 0 ? `${MOCK.unread_notifications} новых` : undefined}
        />
      </div>
    </div>
  );
}

function StatCard({
  label,
  value,
  sub,
  highlight,
}: {
  label: string;
  value: string;
  sub: string;
  highlight?: boolean;
}) {
  return (
    <div
      className={`bg-white rounded-xl border p-4 ${highlight ? 'border-brand-green-100' : 'border-gray-200'}`}
    >
      <p className="text-xs text-gray-500 mb-1">{label}</p>
      <p className={`text-lg font-bold ${highlight ? 'text-brand-green' : 'text-gray-900'}`}>
        {value}
      </p>
      <p className="text-xs text-gray-400 mt-0.5">{sub}</p>
    </div>
  );
}

function QuickActionCard({
  href,
  icon,
  title,
  desc,
  badge,
}: {
  href: string;
  icon: string;
  title: string;
  desc: string;
  badge?: string;
}) {
  return (
    <Link
      href={href}
      className="bg-white rounded-xl border border-gray-200 p-4 flex items-start gap-3 hover:border-brand-green-100 hover:shadow-card transition-all group"
    >
      <span className="text-2xl">{icon}</span>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 flex-wrap">
          <p className="text-sm font-medium text-gray-900 group-hover:text-brand-green">{title}</p>
          {badge && <Badge variant="gold">{badge}</Badge>}
        </div>
        <p className="text-xs text-gray-500 mt-0.5">{desc}</p>
      </div>
      <span className="text-gray-300 group-hover:text-brand-green text-sm shrink-0">→</span>
    </Link>
  );
}
