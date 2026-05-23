import Link from 'next/link';
import { ApplicationStatusBadge } from '@/components/ui/badge';
import type { ApplicationStatus } from '@pkt/shared';
import type { Metadata } from 'next';

export const metadata: Metadata = { title: 'Экспертиза' };

interface QueueItem {
  id: string;
  borrower: string;
  program: string;
  amount: number;
  term_months: number;
  status: ApplicationStatus;
  priority: 'high' | 'normal';
  created_at: string;
  updated_at: string;
  assigned_to?: string;
}

// Mock — replace with: apiClient.get('/expertise/queue', { token: session.accessToken })
const MOCK_QUEUE: QueueItem[] = [
  {
    id: 'APP-2026-0006',
    borrower: 'Сейткали М.Т.',
    program: 'Животноводство 2025',
    amount: 7_500_000,
    term_months: 36,
    status: 'received',
    priority: 'normal',
    created_at: '2026-05-22',
    updated_at: '2026-05-22',
  },
  {
    id: 'APP-2026-0001',
    borrower: 'Касымов Т.А.',
    program: 'Агробизнес 2025',
    amount: 12_000_000,
    term_months: 48,
    status: 'primary_scoring',
    priority: 'high',
    created_at: '2026-05-20',
    updated_at: '2026-05-23',
    assigned_to: 'Ахметов Р.С.',
  },
  {
    id: 'APP-2026-0002',
    borrower: 'Джаксыбеков А.Б.',
    program: 'Растениеводство',
    amount: 8_500_000,
    term_months: 24,
    status: 'collateral_expertise',
    priority: 'normal',
    created_at: '2026-05-10',
    updated_at: '2026-05-20',
    assigned_to: 'Ахметов Р.С.',
  },
  {
    id: 'APP-2026-0003',
    borrower: 'Ибрагимов Р.К.',
    program: 'Животноводство 2025',
    amount: 5_000_000,
    term_months: 36,
    status: 'credit_analysis',
    priority: 'normal',
    created_at: '2026-05-05',
    updated_at: '2026-05-21',
    assigned_to: 'Нурланова А.К.',
  },
  {
    id: 'APP-2026-0004',
    borrower: 'Нурмагамбетова С.Ж.',
    program: 'Агробизнес 2025',
    amount: 20_000_000,
    term_months: 60,
    status: 'credit_committee',
    priority: 'high',
    created_at: '2026-04-25',
    updated_at: '2026-05-22',
    assigned_to: 'Ахметов Р.С.',
  },
  {
    id: 'APP-2026-0005',
    borrower: 'Бекова А.М.',
    program: 'Растениеводство',
    amount: 3_000_000,
    term_months: 24,
    status: 'documentation',
    priority: 'normal',
    created_at: '2026-04-15',
    updated_at: '2026-05-18',
    assigned_to: 'Нурланова А.К.',
  },
];

const STAGE_ORDER: ApplicationStatus[] = [
  'received',
  'primary_scoring',
  'security_check',
  'collateral_expertise',
  'legal_check',
  'credit_analysis',
  'credit_committee',
  'approved',
  'documentation',
  'issued',
  'revision',
  'rejected',
];

function stageIndex(s: ApplicationStatus) {
  return STAGE_ORDER.indexOf(s);
}

export default function ExpertisePage() {
  const sorted = [...MOCK_QUEUE].sort((a, b) => {
    if (a.priority !== b.priority) return a.priority === 'high' ? -1 : 1;
    return stageIndex(a.status) - stageIndex(b.status);
  });

  const high = sorted.filter((a) => a.priority === 'high');
  const normal = sorted.filter((a) => a.priority === 'normal');

  return (
    <div className="max-w-5xl">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Очередь экспертизы</h1>
          <p className="text-sm text-gray-500 mt-1">
            {MOCK_QUEUE.length} заявок · {high.length} приоритетных
          </p>
        </div>
      </div>

      {high.length > 0 && (
        <section className="mb-6">
          <p className="text-xs font-semibold text-red-500 uppercase tracking-wide mb-2">
            Приоритетные
          </p>
          <QueueList items={high} />
        </section>
      )}

      <section>
        <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2">Обычные</p>
        <QueueList items={normal} />
      </section>
    </div>
  );
}

function QueueList({ items }: { items: QueueItem[] }) {
  return (
    <div className="space-y-2">
      {items.map((item) => (
        <Link
          key={item.id}
          href={`/expertise/${item.id}`}
          className="flex items-center gap-4 bg-white rounded-xl border border-gray-200 px-4 py-3 hover:border-brand-green-100 hover:shadow-card transition-all group"
        >
          {item.priority === 'high' && (
            <span className="w-1.5 h-10 rounded-full bg-red-400 shrink-0" />
          )}
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-0.5">
              <span className="text-xs font-mono text-gray-400">{item.id}</span>
              <ApplicationStatusBadge status={item.status} />
            </div>
            <p className="font-medium text-gray-900 group-hover:text-brand-green truncate">
              {item.borrower}
            </p>
            <p className="text-xs text-gray-500">
              {item.program} · {item.amount.toLocaleString('ru-KZ')} ₸ · {item.term_months} мес.
            </p>
          </div>
          <div className="text-right shrink-0 hidden sm:block">
            {item.assigned_to && <p className="text-xs text-gray-400">{item.assigned_to}</p>}
            <p className="text-xs text-gray-400 mt-0.5">
              {new Date(item.updated_at).toLocaleDateString('ru-KZ', {
                day: 'numeric',
                month: 'short',
              })}
            </p>
          </div>
          <span className="text-gray-300 group-hover:text-brand-green text-sm shrink-0">→</span>
        </Link>
      ))}
    </div>
  );
}
