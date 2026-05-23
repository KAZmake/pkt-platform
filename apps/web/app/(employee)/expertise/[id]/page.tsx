import Link from 'next/link';
import { ApplicationStatusBadge } from '@/components/ui/badge';
import { ApplicationTabs } from '@/components/expertise/application-tabs';
import { FsmPanel } from '@/components/expertise/fsm-panel';
import type { ApplicationDetail } from '@/components/expertise/application-tabs';
import type { ApplicationStatus } from '@pkt/shared';
import type { Metadata } from 'next';
import { notFound } from 'next/navigation';

export const metadata: Metadata = { title: 'Карточка заявки' };

// Mock DB — replace with: apiClient.get(`/expertise/applications/${id}`, { token })
const MOCK_APPS: Record<string, ApplicationDetail> = {
  'APP-2026-0002': {
    id: 'APP-2026-0002',
    borrower: 'Джаксыбеков Аманбай Бекович',
    borrower_id: 'BRW-0015',
    inn: '850312301234',
    phone: '+7 701 234 56 78',
    program: 'Растениеводство',
    amount: 8_500_000,
    term_months: 24,
    activity_type: 'Растениеводство (пшеница, ячмень)',
    district: 'Бурлинский район',
    status: 'collateral_expertise',
    created_at: '2026-05-10T09:15:00Z',
    collaterals: [
      {
        id: 'c1',
        type: 'land',
        description: 'Земельный участок 150 га, кадастр № 27:05:000123:45',
        estimated_value: 7_500_000,
        insurance_expiry: '2027-06-01',
        last_inventory: '2026-04-10',
      },
      {
        id: 'c2',
        type: 'equipment',
        description: 'Трактор Кировец К-744, 2021 г.в.',
        estimated_value: 3_200_000,
        insurance_expiry: '2026-12-31',
        last_inventory: '2026-04-10',
      },
    ],
    collateral_conclusion:
      'Залоговое обеспечение достаточно. Рыночная стоимость подтверждена независимым оценщиком. Страхование действительно.',
    legal_conclusion:
      'Юридических ограничений не выявлено. Право собственности подтверждено. Обременений нет.',
    legal_risks: undefined,
    debt_load: 0.38,
    ltv: 0.62,
    analysis_conclusion: 'Долговая нагрузка в норме. LTV 62% — приемлемо. Рекомендуется одобрение.',
    committee_date: undefined,
    documents: [
      {
        id: 'd1',
        name: 'Заявление на кредит',
        category: 'Заявление',
        uploaded_at: '2026-05-10T09:15:00Z',
        uploaded_by: 'Джаксыбеков А.Б.',
      },
      {
        id: 'd2',
        name: 'Паспорт заёмщика',
        category: 'Удостоверение',
        uploaded_at: '2026-05-10T09:15:00Z',
        uploaded_by: 'Джаксыбеков А.Б.',
      },
      {
        id: 'd3',
        name: 'Акт оценки земельного участка',
        category: 'Оценка',
        uploaded_at: '2026-05-12T11:30:00Z',
        uploaded_by: 'Ахметов Р.С.',
      },
    ],
    audit: [
      {
        id: 'a1',
        from_status: null,
        to_status: 'received',
        actor: 'Система',
        comment: 'Заявка зарегистрирована через сайт',
        created_at: '2026-05-10T09:15:00Z',
      },
      {
        id: 'a2',
        from_status: 'received',
        to_status: 'primary_scoring',
        actor: 'Ахметов Р.С.',
        created_at: '2026-05-12T10:00:00Z',
      },
      {
        id: 'a3',
        from_status: 'primary_scoring',
        to_status: 'security_check',
        actor: 'Ахметов Р.С.',
        comment: 'Скоринг пройден, передаю в СБ',
        created_at: '2026-05-14T14:30:00Z',
      },
      {
        id: 'a4',
        from_status: 'security_check',
        to_status: 'collateral_expertise',
        actor: 'Нурланова А.К.',
        comment: 'СБ замечаний не имеет',
        created_at: '2026-05-18T09:00:00Z',
      },
    ],
  },
  'APP-2026-0004': {
    id: 'APP-2026-0004',
    borrower: 'Нурмагамбетова Салтанат Жаксыбековна',
    borrower_id: 'BRW-0008',
    inn: '780905401567',
    phone: '+7 705 987 65 43',
    program: 'Агробизнес 2025',
    amount: 20_000_000,
    term_months: 60,
    activity_type: 'Смешанное производство (скот + зерновые)',
    district: 'Таскалинский район',
    status: 'credit_committee',
    created_at: '2026-04-25T08:00:00Z',
    collaterals: [
      {
        id: 'c1',
        type: 'real_estate',
        description: 'Производственная база, 1200 м²',
        estimated_value: 15_000_000,
        insurance_expiry: '2027-04-01',
        last_inventory: '2026-04-01',
      },
      {
        id: 'c2',
        type: 'livestock',
        description: 'Поголовье КРС — 120 голов',
        estimated_value: 9_600_000,
        insurance_expiry: '2026-10-01',
        last_inventory: '2026-04-01',
      },
    ],
    collateral_conclusion: 'Залог достаточен. Покрытие 124%. Всё застраховано.',
    legal_conclusion: 'Юридических препятствий нет. Документы в порядке.',
    debt_load: 0.42,
    ltv: 0.81,
    analysis_conclusion:
      'Долговая нагрузка приемлема. LTV на верхней границе — рекомендован залог дополнительного имущества или снижение суммы до 18 млн.',
    committee_date: '2026-05-22',
    committee_votes_for: 4,
    committee_votes_against: 1,
    committee_result: 'approved',
    documents: [
      {
        id: 'd1',
        name: 'Заявление',
        category: 'Заявление',
        uploaded_at: '2026-04-25T08:00:00Z',
        uploaded_by: 'Нурмагамбетова С.Ж.',
      },
      {
        id: 'd2',
        name: 'Отчёт об оценке базы',
        category: 'Оценка',
        uploaded_at: '2026-04-30T11:00:00Z',
        uploaded_by: 'Ахметов Р.С.',
      },
      {
        id: 'd3',
        name: 'Протокол КК № 12',
        category: 'Протокол',
        uploaded_at: '2026-05-22T17:00:00Z',
        uploaded_by: 'Система',
      },
    ],
    audit: [
      {
        id: 'a1',
        from_status: null,
        to_status: 'received',
        actor: 'Система',
        created_at: '2026-04-25T08:00:00Z',
      },
      {
        id: 'a2',
        from_status: 'received',
        to_status: 'primary_scoring',
        actor: 'Ахметов Р.С.',
        created_at: '2026-04-26T09:00:00Z',
      },
      {
        id: 'a3',
        from_status: 'primary_scoring',
        to_status: 'security_check',
        actor: 'Ахметов Р.С.',
        created_at: '2026-04-28T10:00:00Z',
      },
      {
        id: 'a4',
        from_status: 'security_check',
        to_status: 'collateral_expertise',
        actor: 'Нурланова А.К.',
        created_at: '2026-05-02T09:00:00Z',
      },
      {
        id: 'a5',
        from_status: 'collateral_expertise',
        to_status: 'legal_check',
        actor: 'Нурланова А.К.',
        created_at: '2026-05-07T14:00:00Z',
      },
      {
        id: 'a6',
        from_status: 'legal_check',
        to_status: 'credit_analysis',
        actor: 'Нурланова А.К.',
        created_at: '2026-05-12T10:00:00Z',
      },
      {
        id: 'a7',
        from_status: 'credit_analysis',
        to_status: 'credit_committee',
        actor: 'Ахметов Р.С.',
        created_at: '2026-05-20T09:00:00Z',
      },
      {
        id: 'a8',
        from_status: 'credit_committee',
        to_status: 'approved',
        actor: 'КК',
        comment: 'Одобрено 4/5 голосов',
        created_at: '2026-05-22T17:00:00Z',
      },
    ],
  },
};

// Minimal fallback for other queue items
function buildFallback(id: string): ApplicationDetail {
  const STATUSES: Record<string, ApplicationStatus> = {
    'APP-2026-0001': 'primary_scoring',
    'APP-2026-0003': 'credit_analysis',
    'APP-2026-0005': 'documentation',
    'APP-2026-0006': 'received',
  };
  const status = STATUSES[id];
  if (!status) return null as unknown as ApplicationDetail;
  return {
    id,
    borrower: 'Заёмщик',
    borrower_id: 'BRW-0000',
    inn: '000000000000',
    phone: '+7 700 000 00 00',
    program: 'Агробизнес 2025',
    amount: 5_000_000,
    term_months: 36,
    activity_type: 'Растениеводство',
    district: 'Уральск',
    status,
    created_at: new Date().toISOString(),
    collaterals: [],
    documents: [],
    audit: [
      {
        id: 'a1',
        from_status: null,
        to_status: 'received',
        actor: 'Система',
        created_at: new Date().toISOString(),
      },
    ],
  };
}

interface PageProps {
  params: { id: string };
}

export default function ExpertiseDetailPage({ params }: PageProps) {
  const app = MOCK_APPS[params.id] ?? buildFallback(params.id);
  if (!app) notFound();

  return (
    <div className="max-w-4xl">
      {/* Header */}
      <div className="mb-5">
        <Link
          href="/expertise"
          className="text-xs text-brand-green hover:underline flex items-center gap-1 mb-3"
        >
          ← Очередь экспертизы
        </Link>
        <div className="flex items-start justify-between gap-4">
          <div>
            <p className="text-xs text-gray-400 font-mono mb-0.5">{app.id}</p>
            <h1 className="text-xl font-bold text-gray-900">{app.borrower}</h1>
            <p className="text-sm text-gray-500 mt-0.5">
              {app.amount.toLocaleString('ru-KZ')} ₸ · {app.term_months} мес. · {app.program}
            </p>
          </div>
          <ApplicationStatusBadge status={app.status} />
        </div>
      </div>

      {/* FSM transitions */}
      <div className="mb-5">
        <FsmPanel applicationId={app.id} currentStatus={app.status} />
      </div>

      {/* Tabs */}
      <ApplicationTabs app={app} />
    </div>
  );
}
