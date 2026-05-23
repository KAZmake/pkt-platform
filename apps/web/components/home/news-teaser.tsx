import Link from 'next/link';
import type { NewsItem } from '@/lib/directus';
import { Button } from '@/components/ui/button';
import { formatDate } from '@/lib/utils';

function NewsCard({ item }: { item: NewsItem }) {
  return (
    <article className="flex flex-col gap-3">
      {item.cover_image && (
        <div className="aspect-video rounded-lg bg-gray-100 overflow-hidden">
          <img
            src={`${process.env.DIRECTUS_URL ?? 'http://localhost:8055'}/assets/${item.cover_image}`}
            alt={item.title}
            className="h-full w-full object-cover"
          />
        </div>
      )}
      <div className="flex-1">
        {item.category && (
          <span className="text-xs font-medium text-brand-gold uppercase tracking-wide">
            {item.category}
          </span>
        )}
        <Link href={`/news/${item.slug}`}>
          <h3 className="mt-1 font-semibold text-gray-900 hover:text-brand-green transition-colors leading-snug">
            {item.title}
          </h3>
        </Link>
        <p className="mt-1 text-sm text-gray-500 line-clamp-2">{item.summary}</p>
      </div>
      <p className="text-xs text-gray-400">{formatDate(item.published_at)}</p>
    </article>
  );
}

export function NewsTeaser({ news }: { news: NewsItem[] }) {
  const latest = news.slice(0, 3);
  return (
    <section className="py-16 bg-gray-50">
      <div className="container-page">
        <div className="flex items-end justify-between mb-8">
          <div>
            <p className="text-brand-gold text-xs font-semibold uppercase tracking-widest mb-2">
              Актуально
            </p>
            <h2 className="text-2xl sm:text-3xl font-bold text-gray-900">Новости</h2>
          </div>
          <Link href="/news" className="hidden sm:block">
            <Button variant="ghost" size="sm">
              Все новости →
            </Button>
          </Link>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
          {latest.length > 0 ? (
            latest.map((item) => <NewsCard key={item.id} item={item} />)
          ) : (
            <p className="text-gray-400 text-sm col-span-3">Новости появятся здесь.</p>
          )}
        </div>
      </div>
    </section>
  );
}
