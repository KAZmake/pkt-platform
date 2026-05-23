import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Главная',
};

export default function HomePage() {
  return (
    <div className="container-page py-16 text-center">
      <h1 className="text-3xl font-bold text-brand-green mb-4">Первое кредитное товарищество</h1>
      <p className="text-gray-600 max-w-xl mx-auto">
        Кредитование агросектора Западно-Казахстанской области.
      </p>
    </div>
  );
}
