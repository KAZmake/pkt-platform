import { Badge } from '@/components/ui/badge';
import type { Metadata } from 'next';

export const metadata: Metadata = { title: 'Документы' };

type DocCategory = 'contract' | 'schedule' | 'certificate' | 'insurance' | 'other';

interface CabinetDocument {
  id: string;
  name: string;
  category: DocCategory;
  date: string;
  size: string;
  /** In production: presigned MinIO URL from apiClient.get('/cabinet/documents') */
  url: string;
}

// Mock — replace with: apiClient.get('/cabinet/documents', { token: session.accessToken })
const MOCK_DOCS: CabinetDocument[] = [
  {
    id: 'd1',
    name: 'Договор займа № LN-2024-0042',
    category: 'contract',
    date: '2024-06-01',
    size: '384 КБ',
    url: '#',
  },
  {
    id: 'd2',
    name: 'График платежей (оригинал)',
    category: 'schedule',
    date: '2024-06-01',
    size: '120 КБ',
    url: '#',
  },
  {
    id: 'd3',
    name: 'Справка об отсутствии задолженности',
    category: 'certificate',
    date: '2024-09-15',
    size: '89 КБ',
    url: '#',
  },
  {
    id: 'd4',
    name: 'Договор страхования залога',
    category: 'insurance',
    date: '2024-06-05',
    size: '250 КБ',
    url: '#',
  },
  {
    id: 'd5',
    name: 'Согласие на обработку персональных данных',
    category: 'other',
    date: '2024-05-20',
    size: '54 КБ',
    url: '#',
  },
  {
    id: 'd6',
    name: 'График платежей (пересмотренный)',
    category: 'schedule',
    date: '2025-01-10',
    size: '130 КБ',
    url: '#',
  },
  {
    id: 'd7',
    name: 'Выписка по займу за апрель 2026',
    category: 'certificate',
    date: '2026-05-01',
    size: '97 КБ',
    url: '#',
  },
];

const CATEGORY_LABELS: Record<DocCategory, string> = {
  contract: 'Договор',
  schedule: 'График',
  certificate: 'Справка',
  insurance: 'Страхование',
  other: 'Прочее',
};

const CATEGORY_VARIANT: Record<DocCategory, 'green' | 'blue' | 'gold' | 'purple' | 'default'> = {
  contract: 'green',
  schedule: 'blue',
  certificate: 'gold',
  insurance: 'purple',
  other: 'default',
};

export default function DocumentsPage() {
  return (
    <div className="max-w-4xl">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Документы</h1>
        <p className="text-sm text-gray-500 mt-1">{MOCK_DOCS.length} документов в архиве</p>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide">
                Документ
              </th>
              <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide hidden sm:table-cell">
                Тип
              </th>
              <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide hidden md:table-cell">
                Дата
              </th>
              <th className="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide hidden md:table-cell">
                Размер
              </th>
              <th className="text-right px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide">
                Действия
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {MOCK_DOCS.map((doc) => (
              <tr key={doc.id} className="hover:bg-gray-50 transition-colors">
                <td className="px-4 py-3">
                  <div className="flex items-center gap-2">
                    <span className="text-lg">📄</span>
                    <span className="font-medium text-gray-900">{doc.name}</span>
                  </div>
                </td>
                <td className="px-4 py-3 hidden sm:table-cell">
                  <Badge variant={CATEGORY_VARIANT[doc.category]}>
                    {CATEGORY_LABELS[doc.category]}
                  </Badge>
                </td>
                <td className="px-4 py-3 text-gray-500 hidden md:table-cell">
                  {new Date(doc.date).toLocaleDateString('ru-KZ')}
                </td>
                <td className="px-4 py-3 text-gray-400 text-xs hidden md:table-cell">{doc.size}</td>
                <td className="px-4 py-3 text-right">
                  <a
                    href={doc.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1 text-sm text-brand-green font-medium hover:underline"
                  >
                    Скачать
                  </a>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <p className="text-xs text-gray-400 mt-3">
        Документы доступны для скачивания в формате PDF. Ссылки сгенерированы через MinIO presigned
        URL и действительны 15 минут.
      </p>
    </div>
  );
}
