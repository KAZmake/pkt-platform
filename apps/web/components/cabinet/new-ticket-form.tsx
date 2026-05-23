'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';

const CATEGORIES = [
  { value: 'payment', label: 'Вопрос по оплате' },
  { value: 'document', label: 'Запрос документа' },
  { value: 'question', label: 'Общий вопрос' },
  { value: 'complaint', label: 'Жалоба' },
  { value: 'other', label: 'Прочее' },
];

interface FormState {
  category: string;
  subject: string;
  message: string;
}

export function NewTicketForm() {
  const [open, setOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [success, setSuccess] = useState(false);
  const [form, setForm] = useState<FormState>({
    category: 'question',
    subject: '',
    message: '',
  });

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSubmitting(true);
    // In production: await apiClient.post('/cabinet/tickets', form, { token: session.accessToken })
    await new Promise((r) => setTimeout(r, 800));
    setSubmitting(false);
    setSuccess(true);
    setOpen(false);
    setForm({ category: 'question', subject: '', message: '' });
    setTimeout(() => setSuccess(false), 5000);
  }

  return (
    <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
      {success && (
        <div className="bg-emerald-50 border-b border-emerald-200 px-4 py-3 text-sm text-emerald-700">
          Обращение создано. Мы свяжемся с вами в рабочее время (пн–пт, 9:00–18:00).
        </div>
      )}
      <button
        onClick={() => setOpen((v) => !v)}
        className="w-full flex items-center justify-between px-4 py-3 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
        type="button"
      >
        <span className="flex items-center gap-2">
          <span>✏️</span>
          Создать новое обращение
        </span>
        <span className="text-gray-400 text-xs">{open ? '▲' : '▼'}</span>
      </button>

      {open && (
        <form onSubmit={handleSubmit} className="px-4 pb-4 pt-3 border-t border-gray-100 space-y-3">
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">Категория</label>
            <select
              value={form.category}
              onChange={(e) => setForm((f) => ({ ...f, category: e.target.value }))}
              className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-green bg-white"
            >
              {CATEGORIES.map((c) => (
                <option key={c.value} value={c.value}>
                  {c.label}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">Тема</label>
            <input
              type="text"
              required
              maxLength={200}
              value={form.subject}
              onChange={(e) => setForm((f) => ({ ...f, subject: e.target.value }))}
              placeholder="Кратко опишите проблему или вопрос"
              className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-green"
            />
          </div>

          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">Сообщение</label>
            <textarea
              required
              rows={4}
              maxLength={2000}
              value={form.message}
              onChange={(e) => setForm((f) => ({ ...f, message: e.target.value }))}
              placeholder="Подробно опишите ваш вопрос или проблему..."
              className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-green resize-none"
            />
          </div>

          <div className="flex gap-2 justify-end pt-1">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() => setOpen(false)}
              disabled={submitting}
            >
              Отмена
            </Button>
            <Button type="submit" size="sm" isLoading={submitting}>
              Отправить
            </Button>
          </div>
        </form>
      )}
    </div>
  );
}
