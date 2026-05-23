'use client';

import { useState } from 'react';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import type { ApplicationStatus, CollateralType } from '@pkt/shared';

// ─── Types ────────────────────────────────────────────────────────────────────

interface AuditEntry {
  id: string;
  from_status: ApplicationStatus | null;
  to_status: ApplicationStatus;
  actor: string;
  comment?: string;
  created_at: string;
}

interface CollateralItem {
  id: string;
  type: CollateralType;
  description: string;
  estimated_value: number;
  insurance_expiry?: string;
  last_inventory?: string;
}

interface ApplicationDocument {
  id: string;
  name: string;
  category: string;
  uploaded_at: string;
  uploaded_by: string;
}

export interface ApplicationDetail {
  id: string;
  borrower: string;
  borrower_id: string;
  inn: string;
  phone: string;
  program: string;
  amount: number;
  term_months: number;
  activity_type: string;
  district: string;
  status: ApplicationStatus;
  created_at: string;
  // Collateral
  collaterals: CollateralItem[];
  collateral_conclusion?: string;
  // Legal
  legal_conclusion?: string;
  legal_risks?: string;
  // Analysis
  debt_load?: number;
  ltv?: number;
  analysis_conclusion?: string;
  // Committee
  committee_date?: string;
  committee_result?: 'approved' | 'rejected' | 'revision';
  committee_votes_for?: number;
  committee_votes_against?: number;
  // Documents
  documents: ApplicationDocument[];
  // Audit log
  audit: AuditEntry[];
}

// ─── Tab definitions ──────────────────────────────────────────────────────────

const TABS = [
  { id: 'summary', label: 'Сводка' },
  { id: 'collateral', label: 'Залог' },
  { id: 'analysis', label: 'Анализ' },
  { id: 'committee', label: 'КК' },
  { id: 'documents', label: 'Документы' },
  { id: 'audit', label: 'Журнал' },
] as const;

type TabId = (typeof TABS)[number]['id'];

const COLLATERAL_LABELS: Record<CollateralType, string> = {
  land: 'Земля',
  equipment: 'Оборудование',
  livestock: 'Скот',
  real_estate: 'Недвижимость',
  other: 'Прочее',
};

// ─── Component ────────────────────────────────────────────────────────────────

export function ApplicationTabs({ app }: { app: ApplicationDetail }) {
  const [activeTab, setActiveTab] = useState<TabId>('summary');

  const totalCollateral = app.collaterals.reduce((s, c) => s + c.estimated_value, 0);

  return (
    <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
      {/* Tab bar */}
      <div className="flex border-b border-gray-200 overflow-x-auto">
        {TABS.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={cn(
              'px-4 py-3 text-sm font-medium whitespace-nowrap transition-colors border-b-2 -mb-px',
              activeTab === tab.id
                ? 'border-brand-green text-brand-green'
                : 'border-transparent text-gray-500 hover:text-gray-900',
            )}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <div className="p-5">
        {/* ── Сводка ── */}
        {activeTab === 'summary' && (
          <div className="grid md:grid-cols-2 gap-6">
            <section>
              <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-3">
                Заёмщик
              </p>
              <dl className="space-y-2">
                <Row label="ФИО" value={app.borrower} />
                <Row label="ИИН/БИН" value={app.inn} />
                <Row label="Телефон" value={app.phone} />
                <Row label="Район" value={app.district} />
                <Row label="Вид деятельности" value={app.activity_type} />
              </dl>
            </section>
            <section>
              <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-3">
                Параметры займа
              </p>
              <dl className="space-y-2">
                <Row label="Программа" value={app.program} />
                <Row label="Сумма" value={app.amount.toLocaleString('ru-KZ') + ' ₸'} />
                <Row label="Срок" value={`${app.term_months} мес.`} />
                <Row label="Подана" value={new Date(app.created_at).toLocaleDateString('ru-KZ')} />
              </dl>
            </section>
          </div>
        )}

        {/* ── Залог ── */}
        {activeTab === 'collateral' && (
          <div>
            <div className="flex items-center justify-between mb-4">
              <p className="text-sm font-medium text-gray-700">
                Залоговое обеспечение · итого:{' '}
                <span className="text-brand-green font-semibold">
                  {totalCollateral.toLocaleString('ru-KZ')} ₸
                </span>
              </p>
            </div>
            <div className="overflow-x-auto mb-4">
              <table className="w-full text-sm">
                <thead className="bg-gray-50 border-y border-gray-200">
                  <tr>
                    <th className="text-left px-3 py-2 text-xs font-semibold text-gray-500">Вид</th>
                    <th className="text-left px-3 py-2 text-xs font-semibold text-gray-500">
                      Описание
                    </th>
                    <th className="text-right px-3 py-2 text-xs font-semibold text-gray-500">
                      Оценка
                    </th>
                    <th className="text-left px-3 py-2 text-xs font-semibold text-gray-500 hidden md:table-cell">
                      Страховка до
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {app.collaterals.map((c) => (
                    <tr key={c.id}>
                      <td className="px-3 py-2">
                        <Badge variant="default">{COLLATERAL_LABELS[c.type]}</Badge>
                      </td>
                      <td className="px-3 py-2 text-gray-700">{c.description}</td>
                      <td className="px-3 py-2 text-right tabular-nums font-medium">
                        {c.estimated_value.toLocaleString('ru-KZ')} ₸
                      </td>
                      <td className="px-3 py-2 text-gray-500 hidden md:table-cell">
                        {c.insurance_expiry
                          ? new Date(c.insurance_expiry).toLocaleDateString('ru-KZ')
                          : '—'}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            {app.collateral_conclusion && (
              <div className="bg-gray-50 rounded-lg p-4">
                <p className="text-xs font-semibold text-gray-500 mb-1">Заключение по залогу</p>
                <p className="text-sm text-gray-700">{app.collateral_conclusion}</p>
              </div>
            )}
          </div>
        )}

        {/* ── Анализ ── */}
        {activeTab === 'analysis' && (
          <div className="space-y-5">
            <section>
              <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-3">
                Юридическая проверка
              </p>
              {app.legal_conclusion ? (
                <div className="space-y-3">
                  <div className="bg-gray-50 rounded-lg p-4">
                    <p className="text-xs font-semibold text-gray-500 mb-1">Заключение</p>
                    <p className="text-sm text-gray-700">{app.legal_conclusion}</p>
                  </div>
                  {app.legal_risks && (
                    <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                      <p className="text-xs font-semibold text-yellow-700 mb-1">Риски</p>
                      <p className="text-sm text-yellow-800">{app.legal_risks}</p>
                    </div>
                  )}
                </div>
              ) : (
                <p className="text-sm text-gray-400">Юридическая проверка ещё не проведена.</p>
              )}
            </section>

            <section>
              <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-3">
                Кредитный анализ
              </p>
              {app.debt_load !== undefined ? (
                <div className="space-y-3">
                  <div className="grid grid-cols-2 gap-3">
                    <div className="bg-gray-50 rounded-lg p-3">
                      <p className="text-xs text-gray-500">Долговая нагрузка</p>
                      <p
                        className={`text-lg font-bold mt-0.5 ${app.debt_load > 0.5 ? 'text-red-600' : 'text-brand-green'}`}
                      >
                        {(app.debt_load * 100).toFixed(0)}%
                      </p>
                    </div>
                    {app.ltv !== undefined && (
                      <div className="bg-gray-50 rounded-lg p-3">
                        <p className="text-xs text-gray-500">LTV (займ / залог)</p>
                        <p
                          className={`text-lg font-bold mt-0.5 ${app.ltv > 0.8 ? 'text-red-600' : 'text-brand-green'}`}
                        >
                          {(app.ltv * 100).toFixed(0)}%
                        </p>
                      </div>
                    )}
                  </div>
                  {app.analysis_conclusion && (
                    <div className="bg-gray-50 rounded-lg p-4">
                      <p className="text-xs font-semibold text-gray-500 mb-1">
                        Заключение аналитика
                      </p>
                      <p className="text-sm text-gray-700">{app.analysis_conclusion}</p>
                    </div>
                  )}
                </div>
              ) : (
                <p className="text-sm text-gray-400">Кредитный анализ ещё не проведён.</p>
              )}
            </section>
          </div>
        )}

        {/* ── КК ── */}
        {activeTab === 'committee' && (
          <div>
            {app.committee_date ? (
              <div className="space-y-4">
                <div className="flex items-center gap-4">
                  <div className="bg-gray-50 rounded-lg p-3 text-center flex-1">
                    <p className="text-xs text-gray-500">Дата заседания</p>
                    <p className="font-semibold text-gray-900 mt-0.5">
                      {new Date(app.committee_date).toLocaleDateString('ru-KZ', {
                        day: 'numeric',
                        month: 'long',
                        year: 'numeric',
                      })}
                    </p>
                  </div>
                  {app.committee_votes_for !== undefined && (
                    <>
                      <div className="bg-emerald-50 rounded-lg p-3 text-center flex-1">
                        <p className="text-xs text-emerald-600">За</p>
                        <p className="text-2xl font-bold text-emerald-600 mt-0.5">
                          {app.committee_votes_for}
                        </p>
                      </div>
                      <div className="bg-red-50 rounded-lg p-3 text-center flex-1">
                        <p className="text-xs text-red-500">Против</p>
                        <p className="text-2xl font-bold text-red-500 mt-0.5">
                          {app.committee_votes_against ?? 0}
                        </p>
                      </div>
                    </>
                  )}
                  {app.committee_result && (
                    <div className="flex-1 text-center">
                      <p className="text-xs text-gray-500 mb-1">Решение</p>
                      <Badge
                        variant={
                          app.committee_result === 'approved'
                            ? 'green'
                            : app.committee_result === 'rejected'
                              ? 'red'
                              : 'yellow'
                        }
                      >
                        {app.committee_result === 'approved'
                          ? 'Одобрено'
                          : app.committee_result === 'rejected'
                            ? 'Отказано'
                            : 'На доработку'}
                      </Badge>
                    </div>
                  )}
                </div>
              </div>
            ) : (
              <p className="text-sm text-gray-400">Заседание КК ещё не назначено.</p>
            )}
          </div>
        )}

        {/* ── Документы ── */}
        {activeTab === 'documents' && (
          <div>
            {app.documents.length > 0 ? (
              <div className="space-y-2">
                {app.documents.map((doc) => (
                  <div
                    key={doc.id}
                    className="flex items-center gap-3 p-3 rounded-lg border border-gray-100 hover:bg-gray-50"
                  >
                    <span className="text-xl">📄</span>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">{doc.name}</p>
                      <p className="text-xs text-gray-400">
                        {doc.category} · {doc.uploaded_by} ·{' '}
                        {new Date(doc.uploaded_at).toLocaleDateString('ru-KZ')}
                      </p>
                    </div>
                    <a
                      href="#"
                      className="text-sm text-brand-green font-medium hover:underline shrink-0"
                    >
                      Открыть
                    </a>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-sm text-gray-400">Документы не прикреплены.</p>
            )}
          </div>
        )}

        {/* ── Журнал ── */}
        {activeTab === 'audit' && (
          <div className="space-y-3">
            {app.audit.map((entry, i) => (
              <div key={entry.id} className="flex gap-3">
                <div className="flex flex-col items-center">
                  <div className="w-2 h-2 rounded-full bg-brand-green mt-1.5 shrink-0" />
                  {i < app.audit.length - 1 && <div className="w-px flex-1 bg-gray-200 my-1" />}
                </div>
                <div className="pb-3 min-w-0">
                  <div className="flex items-center gap-2 flex-wrap">
                    {entry.from_status && (
                      <>
                        <Badge variant="default">{entry.from_status}</Badge>
                        <span className="text-gray-400 text-xs">→</span>
                      </>
                    )}
                    <Badge variant="green">{entry.to_status}</Badge>
                  </div>
                  <p className="text-xs text-gray-500 mt-1">
                    {entry.actor} ·{' '}
                    {new Date(entry.created_at).toLocaleString('ru-KZ', {
                      day: 'numeric',
                      month: 'short',
                      year: 'numeric',
                      hour: '2-digit',
                      minute: '2-digit',
                    })}
                  </p>
                  {entry.comment && (
                    <p className="text-xs text-gray-600 mt-0.5 italic">{entry.comment}</p>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex gap-2">
      <dt className="text-xs text-gray-400 w-36 shrink-0 pt-0.5">{label}</dt>
      <dd className="text-sm text-gray-900">{value}</dd>
    </div>
  );
}
