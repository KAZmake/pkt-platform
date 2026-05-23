import type { Metadata } from 'next';
import { Hero } from '@/components/home/hero';
import { KpiSection } from '@/components/home/kpi-counter';
import { ZkoMapPreview } from '@/components/home/zko-map';
import { ProgramsTeaser } from '@/components/home/programs-teaser';
import { NewsTeaser } from '@/components/home/news-teaser';
import { getPrograms } from '@/lib/api/programs';
import { getNews } from '@/lib/directus';

export const metadata: Metadata = { title: 'Главная' };

export default async function HomePage() {
  const [programs, news] = await Promise.all([
    getPrograms().catch(() => []),
    getNews(3).catch(() => []),
  ]);

  return (
    <>
      <Hero />
      <KpiSection />
      <ProgramsTeaser programs={programs} />
      <ZkoMapPreview />
      <NewsTeaser news={news} />
    </>
  );
}
