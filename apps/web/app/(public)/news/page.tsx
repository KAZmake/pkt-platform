import type { Metadata } from 'next';
import Link from 'next/link';
import { getNews } from '@/lib/directus';
import { formatDate } from '@/lib/utils';

export const metadata: Metadata = {
  title: 'Новости',
  description: 'Последние новости ТОО «Первое кредитное товарищество».',
};

export default async function NewsPage() {
  const news = await getNews().catch(() => []);

  return (
    <div className="container-page py-12">
      <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-8">Новости</h1>

      {news.length === 0 ? (
        <p className="text-gray-400 text-sm">Новости появятся здесь.</p>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
          {news.map((item) => (
            <article key={item.id} className="flex flex-col gap-3">
              {item.cover_image && (
                <Link href={`/news/${item.slug}`}>
                  <div className="aspect-video rounded-xl bg-gray-100 overflow-hidden">
                    <img
                      src={`${process.env.DIRECTUS_URL ?? 'http://localhost:8055'}/assets/${item.cover_image}`}
                      alt={item.title}
                      className="h-full w-full object-cover hover:scale-105 transition-transform duration-300"
                    />
                  </div>
                </Link>
              )}
              <div className="flex-1">
                {item.category && (
                  <span className="text-xs font-medium text-brand-gold uppercase tracking-wide">
                    {item.category}
                  </span>
                )}
                <Link href={`/news/${item.slug}`}>
                  <h2 className="mt-1 font-semibold text-gray-900 hover:text-brand-green transition-colors leading-snug">
                    {item.title}
                  </h2>
                </Link>
                <p className="mt-2 text-sm text-gray-500 line-clamp-3">{item.summary}</p>
              </div>
              <div className="flex items-center justify-between text-xs text-gray-400">
                <span>{formatDate(item.published_at)}</span>
                <Link
                  href={`/news/${item.slug}`}
                  className="text-brand-green hover:underline font-medium"
                >
                  Читать →
                </Link>
              </div>
            </article>
          ))}
        </div>
      )}
    </div>
  );
}
