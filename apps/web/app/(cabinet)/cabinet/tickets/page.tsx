import { Badge } from '@/components/ui/badge';
import { NewTicketForm } from '@/components/cabinet/new-ticket-form';
import type { Metadata } from 'next';

export const metadata: Metadata = { title: 'Обращения' };

type TicketStatus = 'open' | 'in_progress' | 'resolved' | 'closed';
type TicketCategory = 'payment' | 'document' | 'question' | 'complaint' | 'other';

interface Ticket {
  id: string;
  subject: string;
  category: TicketCategory;
  status: TicketStatus;
  created_at: string;
  updated_at: string;
  last_reply?: string;
}

// Mock — replace with: apiClient.get('/cabinet/tickets', { token: session.accessToken })
const MOCK_TICKETS: Ticket[] = [
  {
    id: 'TK-0001',
    subject: 'Запрос справки об остатке долга',
    category: 'document',
    status: 'resolved',
    created_at: '2026-04-10',
    updated_at: '2026-04-12',
    last_reply: 'Справка готова и отправлена на ваш email.',
  },
  {
    id: 'TK-0002',
    subject: 'Вопрос по досрочному погашению займа',
    category: 'payment',
    status: 'in_progress',
    created_at: '2026-05-18',
    updated_at: '2026-05-20',
    last_reply:
      'Ваш вопрос передан кредитному менеджеру. Ожидайте ответа в течение 2 рабочих дней.',
  },
  {
    id: 'TK-0003',
    subject: 'Изменение контактных данных',
    category: 'other',
    status: 'open',
    created_at: '2026-05-23',
    updated_at: '2026-05-23',
  },
];

const STATUS_LABELS: Record<TicketStatus, string> = {
  open: 'Открыто',
  in_progress: 'В работе',
  resolved: 'Решено',
  closed: 'Закрыто',
};

const STATUS_VARIANT: Record<TicketStatus, 'default' | 'blue' | 'green' | 'gold'> = {
  open: 'gold',
  in_progress: 'blue',
  resolved: 'green',
  closed: 'default',
};

const CATEGORY_LABELS: Record<TicketCategory, string> = {
  payment: 'Оплата',
  document: 'Документ',
  question: 'Вопрос',
  complaint: 'Жалоба',
  other: 'Прочее',
};

export default function TicketsPage() {
  const open = MOCK_TICKETS.filter((t) => t.status === 'open' || t.status === 'in_progress');

  return (
    <div className="max-w-3xl">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Обращения</h1>
        <p className="text-sm text-gray-500 mt-1">
          {open.length > 0 ? `${open.length} активных обращений` : 'Нет активных обращений'}
        </p>
      </div>

      {/* Create ticket form */}
      <NewTicketForm />

      {/* Tickets list */}
      {MOCK_TICKETS.length > 0 && (
        <div className="space-y-3 mt-6">
          <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide">
            История обращений
          </p>
          {MOCK_TICKETS.map((ticket) => (
            <div key={ticket.id} className="bg-white rounded-xl border border-gray-200 p-4">
              <div className="flex items-start justify-between gap-3">
                <div className="min-w-0">
                  <p className="text-xs text-gray-400 mb-0.5">
                    {ticket.id} · {CATEGORY_LABELS[ticket.category]}
                  </p>
                  <p className="font-medium text-gray-900">{ticket.subject}</p>
                </div>
                <Badge variant={STATUS_VARIANT[ticket.status]} className="shrink-0">
                  {STATUS_LABELS[ticket.status]}
                </Badge>
              </div>

              {ticket.last_reply && (
                <div className="mt-3 bg-gray-50 rounded-lg px-3 py-2 text-xs text-gray-600 border-l-2 border-brand-green-100">
                  {ticket.last_reply}
                </div>
              )}

              <div className="flex gap-4 mt-2 text-xs text-gray-400">
                <span>
                  Создано:{' '}
                  {new Date(ticket.created_at).toLocaleDateString('ru-KZ', {
                    day: 'numeric',
                    month: 'short',
                    year: 'numeric',
                  })}
                </span>
                {ticket.updated_at !== ticket.created_at && (
                  <span>
                    Обновлено:{' '}
                    {new Date(ticket.updated_at).toLocaleDateString('ru-KZ', {
                      day: 'numeric',
                      month: 'short',
                      year: 'numeric',
                    })}
                  </span>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
