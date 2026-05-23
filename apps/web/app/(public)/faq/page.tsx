import type { Metadata } from 'next';
import { getFaq, type FaqItem } from '@/lib/directus';
import { Accordion } from '@/components/ui/accordion';

export const metadata: Metadata = {
  title: 'Вопросы и ответы',
  description: 'Ответы на часто задаваемые вопросы о кредитовании в ПКТ.',
};

export default async function FaqPage() {
  const items: FaqItem[] = await getFaq().catch(() => []);

  // Group by category
  const grouped: Record<string, FaqItem[]> = {};
  for (const item of items) {
    const cat = item.category ?? 'Общие вопросы';
    if (!grouped[cat]) grouped[cat] = [];
    grouped[cat].push(item);
  }

  const categories = Object.entries(grouped);

  return (
    <div className="container-page py-12 max-w-3xl">
      <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">Вопросы и ответы</h1>
      <p className="text-gray-500 text-sm mb-10">
        Не нашли ответ?{' '}
        <a href="/contacts" className="text-brand-green hover:underline">
          Свяжитесь с нами
        </a>{' '}
        или напишите в{' '}
        <a href="/cabinet" className="text-brand-green hover:underline">
          личном кабинете
        </a>
        .
      </p>

      {categories.length === 0 ? (
        <p className="text-gray-400 text-sm">FAQ появится здесь.</p>
      ) : (
        <div className="space-y-10">
          {categories.map(([category, faqItems]) => (
            <section key={category}>
              <h2 className="text-base font-semibold text-gray-900 mb-4">{category}</h2>
              <Accordion
                items={faqItems.map((i) => ({
                  id: i.id,
                  question: i.question,
                  answer: i.answer,
                }))}
              />
            </section>
          ))}
        </div>
      )}
    </div>
  );
}
