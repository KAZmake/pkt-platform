import { signIn } from '@/auth';
import { Button } from '@/components/ui/button';
import type { Metadata } from 'next';

export const metadata: Metadata = { title: 'Вход' };

export default function LoginPage({ searchParams }: { searchParams: { callbackUrl?: string } }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-sm bg-white rounded-2xl border border-gray-200 shadow-card p-8">
        <div className="text-center mb-8">
          <span className="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-brand-green text-white font-bold text-xl mb-4">
            ПКТ
          </span>
          <h1 className="text-xl font-semibold text-gray-900">Вход в систему</h1>
          <p className="text-sm text-gray-500 mt-1">Первое кредитное товарищество</p>
        </div>

        <form
          action={async () => {
            'use server';
            await signIn('keycloak', {
              redirectTo: searchParams.callbackUrl ?? '/cabinet',
            });
          }}
        >
          <Button type="submit" className="w-full" size="lg">
            Войти через корпоративный аккаунт
          </Button>
        </form>

        <p className="mt-6 text-center text-xs text-gray-400">
          Авторизация через Keycloak. При возникновении проблем обратитесь в IT-отдел.
        </p>
      </div>
    </div>
  );
}
