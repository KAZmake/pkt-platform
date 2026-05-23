'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

type Collection = 'news' | 'programs' | 'faq';

interface ContentItem {
  id: string;
  title: string;
  status: 'published' | 'draft';
  updated_at: string;
}

// Mock content — replace with Directus REST:
// GET /items/{collection}?fields=id,title,status,date_updated&sort=-date_updated
const MOCK_CONTENT: Record<Collection, ContentItem[]> = {
  news: [
    {
      id: '1',
      title: 'ПКТ снижает ставки по программе «Агробизнес 2025»',
      status: 'published',
      updated_at: '2026-05-20',
    },
    {
      id: '2',
      title: 'Открытие нового офиса в Аксае',
      status: 'published',
      updated_at: '2026-05-10',
    },
    { id: '3', title: 'Итоги работы за I квартал 2026', status: 'draft', updated_at: '2026-05-22' },
  ],
  programs: [
    { id: '1', title: 'Агробизнес 2025', status: 'published', updated_at: '2026-04-01' },
    { id: '2', title: 'Растениеводство', status: 'published', updated_at: '2026-04-01' },
    { id: '3', title: 'Животноводство 2025', status: 'published', updated_at: '2026-04-01' },
    { id: '4', title: 'Молодой фермер', status: 'draft', updated_at: '2026-05-15' },
  ],
  faq: [
    {
      id: '1',
      title: 'Какие документы нужны для подачи заявки?',
      status: 'published',
      updated_at: '2026-03-01',
    },
    {
      id: '2',
      title: 'Каков минимальный возраст заёмщика?',
      status: 'published',
      updated_at: '2026-03-01',
    },
    {
      id: '3',
      title: 'Можно ли погасить займ досрочно?',
      status: 'published',
      updated_at: '2026-03-01',
    },
  ],
};

const COLLECTION_LABELS: Record<Collection, string> = {
  news: 'Новости',
  programs: 'Программы',
  faq: 'FAQ',
};

interface EditModalState {
  collection: Collection;
  item: ContentItem;
  title: string;
  status: 'published' | 'draft';
}

export function ContentEditor() {
  const [activeTab, setActiveTab] = useState<Collection>('news');
  const [edit, setEdit] = useState<EditModalState | null>(null);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState<string | null>(null);

  async function handleSave() {
    if (!edit) return;
    setSaving(true);
    // In production: PATCH /items/{collection}/{id} with Directus token
    // await directusFetch(`/items/${edit.collection}/${edit.item.id}`, {
    //   method: 'PATCH',
    //   body: JSON.stringify({ title: edit.title, status: edit.status }),
    // });
    await new Promise((r) => setTimeout(r, 600));
    setSaving(false);
    setSaved(edit.item.id);
    setEdit(null);
    setTimeout(() => setSaved(null), 3000);
  }

  const items = MOCK_CONTENT[activeTab];

  return (
    <div>
      {saved && (
        <div className="mb-4 bg-emerald-50 border border-emerald-200 rounded-lg px-4 py-3 text-sm text-emerald-700">
          Изменения сохранены в Directus.
        </div>
      )}

      {/* Collection tabs */}
      <div className="flex gap-1 bg-gray-100 rounded-lg p-1 mb-5 w-fit">
        {(Object.keys(COLLECTION_LABELS) as Collection[]).map((col) => (
          <button
            key={col}
            onClick={() => setActiveTab(col)}
            className={`px-4 py-1.5 rounded-md text-sm font-medium transition-colors ${
              col === activeTab
                ? 'bg-white text-gray-900 shadow-sm'
                : 'text-gray-500 hover:text-gray-900'
            }`}
          >
            {COLLECTION_LABELS[col]}
          </button>
        ))}
      </div>

      {/* Items table */}
      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide">
                Заголовок
              </th>
              <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide hidden sm:table-cell">
                Статус
              </th>
              <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide hidden md:table-cell">
                Обновлено
              </th>
              <th className="text-right px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide">
                Действия
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {items.map((item) => (
              <tr key={item.id} className="hover:bg-gray-50 transition-colors">
                <td className="px-4 py-3 font-medium text-gray-900 max-w-xs truncate">
                  {item.title}
                </td>
                <td className="px-4 py-3 hidden sm:table-cell">
                  <Badge variant={item.status === 'published' ? 'green' : 'default'}>
                    {item.status === 'published' ? 'Опубликован' : 'Черновик'}
                  </Badge>
                </td>
                <td className="px-4 py-3 text-gray-400 text-xs hidden md:table-cell">
                  {new Date(item.updated_at).toLocaleDateString('ru-KZ')}
                </td>
                <td className="px-4 py-3 text-right">
                  <button
                    onClick={() =>
                      setEdit({
                        collection: activeTab,
                        item,
                        title: item.title,
                        status: item.status,
                      })
                    }
                    className="text-sm text-brand-green font-medium hover:underline"
                  >
                    Редактировать
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Edit modal */}
      {edit && (
        <div className="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl shadow-xl w-full max-w-lg p-6">
            <h2 className="text-base font-semibold text-gray-900 mb-4">
              Редактировать — {COLLECTION_LABELS[edit.collection]}
            </h2>

            <div className="space-y-4">
              <div>
                <label className="block text-xs font-medium text-gray-600 mb-1">Заголовок</label>
                <input
                  type="text"
                  value={edit.title}
                  onChange={(e) => setEdit((s) => s && { ...s, title: e.target.value })}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-green"
                />
              </div>

              <div>
                <label className="block text-xs font-medium text-gray-600 mb-1">Статус</label>
                <select
                  value={edit.status}
                  onChange={(e) =>
                    setEdit((s) => s && { ...s, status: e.target.value as 'published' | 'draft' })
                  }
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-green bg-white"
                >
                  <option value="published">Опубликован</option>
                  <option value="draft">Черновик</option>
                </select>
              </div>

              <p className="text-xs text-gray-400">
                Полный редактор содержимого доступен в{' '}
                <a
                  href={`${process.env.NEXT_PUBLIC_DIRECTUS_URL ?? 'http://localhost:8055'}/admin/content/${edit.collection}/${edit.item.id}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-brand-green hover:underline"
                >
                  Directus Admin →
                </a>
              </p>
            </div>

            <div className="flex gap-2 justify-end mt-5">
              <Button variant="ghost" size="sm" onClick={() => setEdit(null)} disabled={saving}>
                Отмена
              </Button>
              <Button size="sm" isLoading={saving} onClick={handleSave}>
                Сохранить
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
