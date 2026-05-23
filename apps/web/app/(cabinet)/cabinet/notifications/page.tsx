import { Badge } from '@/components/ui/badge';
import { MarkReadButton } from '@/components/cabinet/mark-read-button';
import type { Metadata } from 'next';

export const metadata: Metadata = { title: 'Уведомления' };

type NotifType = 'payment' | 'document' | 'system' | 'ticket';

interface Notification {
  id: string;
  type: NotifType;
  title: string;
  body: string;
  date: string;
  read: boolean;
}

// Mock — replace with: apiClient.get('/cabinet/notifications', { token: session.accessToken })
const MOCK_NOTIFICATIONS: Notification[] = [
  {
    id: 'n1',
    type: 'payment',
    title: 'Напоминание о предстоящем платеже',
    body: 'Ваш ежемесячный платёж 312 500 ₸ запланирован на 1 июня 2026 года. Убедитесь, что на счету достаточно средств.',
    date: '2026-05-20',
    read: false,
  },
  {
    id: 'n2',
    type: 'document',
    title: 'Новый документ доступен',
    body: 'Выписка по займу за апрель 2026 года готова к скачиванию в разделе «Документы».',
    date: '2026-05-15',
    read: false,
  },
  {
    id: 'n3',
    type: 'ticket',
    title: 'Ответ по обращению TK-0002',
    body: 'Ваш вопрос по досрочному погашению передан кредитному менеджеру. Ожидайте ответа в течение 2 рабочих дней.',
    date: '2026-05-20',
    read: false,
  },
  {
    id: 'n4',
    type: 'payment',
    title: 'Платёж за апрель подтверждён',
    body: 'Ваш платёж за апрель 2026 года в размере 312 500 ₸ успешно зачислен. Остаток основного долга: 11 250 000 ₸.',
    date: '2026-05-01',
    read: true,
  },
  {
    id: 'n5',
    type: 'system',
    title: 'Плановые технические работы',
    body: 'В ночь с 25 на 26 мая с 01:00 до 04:00 возможны кратковременные перебои в работе личного кабинета. Приносим извинения.',
    date: '2026-05-19',
    read: true,
  },
  {
    id: 'n6',
    type: 'document',
    title: 'Обновлённый график платежей',
    body: 'График платежей (пересмотренный) от 10.01.2025 доступен для скачивания в разделе «Документы».',
    date: '2025-01-10',
    read: true,
  },
];

const TYPE_ICON: Record<NotifType, string> = {
  payment: '💳',
  document: '📄',
  system: '⚙️',
  ticket: '💬',
};

const TYPE_VARIANT: Record<NotifType, 'blue' | 'gold' | 'default' | 'purple'> = {
  payment: 'blue',
  document: 'gold',
  system: 'default',
  ticket: 'purple',
};

const TYPE_LABELS: Record<NotifType, string> = {
  payment: 'Платёж',
  document: 'Документ',
  system: 'Система',
  ticket: 'Обращение',
};

export default function NotificationsPage() {
  const unread = MOCK_NOTIFICATIONS.filter((n) => !n.read);

  return (
    <div className="max-w-3xl">
      <div className="mb-6 flex items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Уведомления</h1>
          <p className="text-sm text-gray-500 mt-1">
            {unread.length > 0 ? `${unread.length} непрочитанных` : 'Все уведомления прочитаны'}
          </p>
        </div>
        {unread.length > 0 && <MarkReadButton />}
      </div>

      <div className="space-y-3">
        {MOCK_NOTIFICATIONS.map((n) => (
          <div
            key={n.id}
            className={`rounded-xl border p-4 ${
              n.read ? 'bg-white border-gray-200' : 'bg-brand-green-50 border-brand-green-100'
            }`}
          >
            <div className="flex items-start gap-3">
              <span className="text-xl shrink-0 mt-0.5">{TYPE_ICON[n.type]}</span>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1 flex-wrap">
                  <p
                    className={`text-sm font-medium ${n.read ? 'text-gray-700' : 'text-gray-900'}`}
                  >
                    {n.title}
                  </p>
                  <Badge variant={TYPE_VARIANT[n.type]}>{TYPE_LABELS[n.type]}</Badge>
                  {!n.read && <span className="w-2 h-2 bg-brand-green rounded-full shrink-0" />}
                </div>
                <p className="text-sm text-gray-500">{n.body}</p>
                <p className="text-xs text-gray-400 mt-2">
                  {new Date(n.date).toLocaleDateString('ru-KZ', {
                    day: 'numeric',
                    month: 'long',
                    year: 'numeric',
                  })}
                </p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
