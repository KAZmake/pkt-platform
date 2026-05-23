import type { Metadata } from 'next';
import Link from 'next/link';
import { notFound } from 'next/navigation';
import { getNewsItem, getNews } from '@/lib/directus';
import { formatDate } from '@/lib/utils';

interface Props {
  params: { slug: string };
}

export async function generateStaticParams() {
  const news = await getNews(100).catch(() => []);
  return news.map((item) => ({ slug: item.slug }));
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const item = await getNewsItem(params.slug).catch(() => null);
  return {
    title: item?.title ?? 'Новость',
    description: item?.summary,
  };
}

export default async function NewsItemPage({ params }: Props) {
  const item = await getNewsItem(params.slug).catch(() => null);
  if (!item) notFound();

  const DIRECTUS_URL = process.env.DIRECTUS_URL ?? 'http://localhost:8055';

  return (
    <div className="container-page py-12 max-w-3xl">
      <nav className="text-sm text-gray-400 mb-6 flex items-center gap-2">
        <Link href="/news" className="hover:text-brand-green">
          Новости
        </Link>
        <span>/</span>
        <span className="text-gray-700 truncate">{item.title}</span>
      </nav>

      {item.cover_image && (
        <div className="aspect-video rounded-xl bg-gray-100 overflow-hidden mb-8">
          <img
            src={`${DIRECTUS_URL}/assets/${item.cover_image}`}
            alt={item.title}
            className="h-full w-full object-cover"
          />
        </div>
      )}

      <div className="flex items-center gap-3 mb-4 text-xs text-gray-400">
        {item.category && (
          <span className="text-brand-gold font-semibold uppercase tracking-wide">
            {item.category}
          </span>
        )}
        <span>{formatDate(item.published_at)}</span>
      </div>

      <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-6">{item.title}</h1>

      <div
        className="prose prose-sm max-w-none text-gray-600 leading-relaxed"
        dangerouslySetInnerHTML={{ __html: item.content }}
      />

      <div className="mt-10 pt-6 border-t border-gray-100">
        <Link href="/news" className="text-sm text-brand-green hover:underline">
          ← Все новости
        </Link>
      </div>
    </div>
  );
}
