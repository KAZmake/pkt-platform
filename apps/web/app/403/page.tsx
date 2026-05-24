import Link from 'next/link';
import type { Metadata } from 'next';

export const metadata: Metadata = { title: 'Доступ запрещён' };

export default function ForbiddenPage() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center gap-4 text-center px-4">
      <p className="text-8xl font-bold text-brand-gold">403</p>
      <h1 className="text-2xl font-semibold text-gray-800">Доступ запрещён</h1>
      <p className="text-gray-500 max-w-sm">
        У вашей учётной записи нет прав для просмотра этого раздела. Обратитесь к администратору,
        если считаете это ошибкой.
      </p>
      <div className="flex gap-3 mt-2">
        <Link
          href="/"
          className="inline-flex items-center gap-2 bg-brand-green text-white px-5 py-2.5 rounded-lg hover:bg-brand-green-dark transition-colors"
        >
          На главную
        </Link>
        <Link
          href="/cabinet"
          className="inline-flex items-center gap-2 border border-gray-300 text-gray-700 px-5 py-2.5 rounded-lg hover:bg-gray-50 transition-colors"
        >
          Личный кабинет
        </Link>
      </div>
    </div>
  );
}
