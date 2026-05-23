'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import type { ApplicationStatus } from '@pkt/shared';

interface Transition {
  to: ApplicationStatus;
  label: string;
  variant: 'primary' | 'secondary' | 'destructive';
}

const FSM: Partial<Record<ApplicationStatus, Transition[]>> = {
  received: [
    { to: 'primary_scoring', label: 'Начать скоринг', variant: 'primary' },
    { to: 'rejected', label: 'Отказать', variant: 'destructive' },
  ],
  primary_scoring: [
    { to: 'security_check', label: 'Передать в СБ', variant: 'primary' },
    { to: 'revision', label: 'На доработку', variant: 'secondary' },
    { to: 'rejected', label: 'Отказать', variant: 'destructive' },
  ],
  security_check: [
    { to: 'collateral_expertise', label: 'Передать на экспертизу залога', variant: 'primary' },
    { to: 'revision', label: 'На доработку', variant: 'secondary' },
    { to: 'rejected', label: 'Отказать', variant: 'destructive' },
  ],
  collateral_expertise: [
    { to: 'legal_check', label: 'Передать на юр. проверку', variant: 'primary' },
    { to: 'revision', label: 'На доработку', variant: 'secondary' },
    { to: 'rejected', label: 'Отказать', variant: 'destructive' },
  ],
  legal_check: [
    { to: 'credit_analysis', label: 'Передать на кред. анализ', variant: 'primary' },
    { to: 'revision', label: 'На доработку', variant: 'secondary' },
    { to: 'rejected', label: 'Отказать', variant: 'destructive' },
  ],
  credit_analysis: [
    { to: 'credit_committee', label: 'Вынести на КК', variant: 'primary' },
    { to: 'revision', label: 'На доработку', variant: 'secondary' },
    { to: 'rejected', label: 'Отказать', variant: 'destructive' },
  ],
  credit_committee: [
    { to: 'approved', label: 'Одобрить', variant: 'primary' },
    { to: 'revision', label: 'На доработку', variant: 'secondary' },
    { to: 'rejected', label: 'Отказать', variant: 'destructive' },
  ],
  approved: [{ to: 'documentation', label: 'Начать оформление', variant: 'primary' }],
  documentation: [{ to: 'issued', label: 'Выдать займ', variant: 'primary' }],
  revision: [{ to: 'primary_scoring', label: 'Вернуть в работу', variant: 'primary' }],
};

const TERMINAL: ApplicationStatus[] = ['issued', 'rejected'];

interface FsmPanelProps {
  applicationId: string;
  currentStatus: ApplicationStatus;
}

export function FsmPanel({ applicationId: _applicationId, currentStatus }: FsmPanelProps) {
  const [status, setStatus] = useState(currentStatus);
  const [loading, setLoading] = useState<ApplicationStatus | null>(null);
  const [comment, setComment] = useState('');
  const [pendingTransition, setPendingTransition] = useState<Transition | null>(null);

  const transitions = FSM[status] ?? [];
  const isTerminal = TERMINAL.includes(status);

  async function confirmTransition() {
    if (!pendingTransition) return;
    setLoading(pendingTransition.to);
    // In production:
    // await apiClient.post(`/expertise/applications/${applicationId}/transition`, {
    //   to: pendingTransition.to, comment,
    // }, { token: session.accessToken });
    await new Promise((r) => setTimeout(r, 700));
    setStatus(pendingTransition.to);
    setLoading(null);
    setPendingTransition(null);
    setComment('');
  }

  if (isTerminal) {
    return (
      <div className="bg-gray-50 rounded-xl border border-gray-200 px-4 py-3 text-sm text-gray-500">
        Заявка находится в терминальном статусе — дальнейшие переходы недоступны.
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl border border-gray-200 p-4">
      <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">
        Переход по статусу
      </p>

      {pendingTransition ? (
        <div className="space-y-3">
          <p className="text-sm text-gray-700">
            Подтвердите переход: <strong>{pendingTransition.label}</strong>
          </p>
          <textarea
            rows={2}
            value={comment}
            onChange={(e) => setComment(e.target.value)}
            placeholder="Комментарий (необязательно)..."
            className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-green resize-none"
          />
          <div className="flex gap-2">
            <Button
              size="sm"
              variant={pendingTransition.variant === 'destructive' ? 'destructive' : 'primary'}
              isLoading={loading === pendingTransition.to}
              onClick={confirmTransition}
            >
              Подтвердить
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={() => setPendingTransition(null)}
              disabled={loading !== null}
            >
              Отмена
            </Button>
          </div>
        </div>
      ) : (
        <div className="flex flex-wrap gap-2">
          {transitions.map((t) => (
            <Button
              key={t.to}
              size="sm"
              variant={t.variant}
              onClick={() => setPendingTransition(t)}
            >
              {t.label}
            </Button>
          ))}
        </div>
      )}
    </div>
  );
}
