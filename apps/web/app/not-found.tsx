import Link from 'next/link';

export default function NotFound() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center gap-4 text-center px-4">
      <p className="text-8xl font-bold text-brand-green">404</p>
      <h1 className="text-2xl font-semibold text-gray-800">Страница не найдена</h1>
      <p className="text-gray-500 max-w-sm">
        Запрашиваемая страница не существует или была перемещена.
      </p>
      <Link
        href="/"
        className="mt-2 inline-flex items-center gap-2 bg-brand-green text-white px-5 py-2.5 rounded-lg hover:bg-brand-green-dark transition-colors"
      >
        На главную
      </Link>
    </div>
  );
}
