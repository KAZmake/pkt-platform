import { ContentEditor } from '@/components/content/content-editor';
import type { Metadata } from 'next';

export const metadata: Metadata = { title: 'Редактор контента' };

export default function ContentPage() {
  const directusAdminUrl = process.env.NEXT_PUBLIC_DIRECTUS_URL
    ? `${process.env.NEXT_PUBLIC_DIRECTUS_URL}/admin`
    : 'http://localhost:8055/admin';

  return (
    <div className="max-w-4xl">
      <div className="mb-6 flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Редактор контента</h1>
          <p className="text-sm text-gray-500 mt-1">Управление публикациями через Directus API</p>
        </div>
        <a
          href={directusAdminUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="text-sm text-brand-green font-medium hover:underline whitespace-nowrap shrink-0"
        >
          Directus Admin →
        </a>
      </div>

      <ContentEditor />
    </div>
  );
}
