import type { Metadata } from 'next';
import Link from 'next/link';
import { Button } from '@/components/ui/button';

export const metadata: Metadata = { title: 'Подать заявку' };

export default function ApplyPage() {
  return (
    <div className="container-page py-20 max-w-xl text-center">
      <div className="text-5xl mb-6">📋</div>
      <h1 className="text-2xl font-bold text-gray-900 mb-3">Подать заявку</h1>
      <p className="text-gray-500 text-sm mb-8">
        Для подачи заявки войдите в личный кабинет. Если у вас ещё нет аккаунта — обратитесь к
        менеджеру по телефону{' '}
        <a href="tel:+77112000000" className="text-brand-green font-medium">
          8 (7112) 00-00-00
        </a>
        .
      </p>
      <div className="flex flex-col sm:flex-row gap-3 justify-center">
        <Link href="/login">
          <Button>Войти в личный кабинет</Button>
        </Link>
        <Link href="/programs">
          <Button variant="outline">Смотреть программы</Button>
        </Link>
      </div>
    </div>
  );
}
