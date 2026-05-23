'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { signOut } from 'next-auth/react';
import { cn } from '@/lib/utils';

const NAV_ITEMS = [
  { href: '/expertise', label: 'Экспертиза', icon: '📋', exact: false },
  { href: '/analytics', label: 'Аналитика', icon: '📊', exact: false },
  { href: '/content', label: 'Контент', icon: '✏️', exact: false },
  { href: '/mail', label: 'Почта и диск', icon: '📧', exact: false },
];

const ROLE_LABELS: Record<string, string> = {
  employee: 'Сотрудник',
  expert: 'Эксперт / КК',
  admin: 'Администратор',
};

interface EmployeeSidebarProps {
  userName: string;
  role: string;
}

export function EmployeeSidebar({ userName, role }: EmployeeSidebarProps) {
  const pathname = usePathname();

  return (
    <aside className="w-60 shrink-0 bg-white border-r border-gray-200 flex-col hidden md:flex">
      {/* Logo */}
      <div className="px-4 py-4 border-b border-gray-200">
        <Link href="/" className="flex items-center gap-2">
          <div className="w-8 h-8 bg-brand-green rounded-lg flex items-center justify-center text-white font-bold text-xs">
            ПКТ
          </div>
          <span className="font-semibold text-gray-900 text-sm">Рабочее место</span>
        </Link>
      </div>

      {/* User */}
      <div className="px-4 py-3 border-b border-gray-100">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-full bg-brand-gold-50 text-brand-gold-500 flex items-center justify-center font-semibold text-sm shrink-0">
            {userName.charAt(0).toUpperCase()}
          </div>
          <div className="min-w-0">
            <p className="text-sm font-medium text-gray-900 truncate">{userName}</p>
            <p className="text-xs text-brand-gold">{ROLE_LABELS[role] ?? role}</p>
          </div>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-3 space-y-0.5">
        {NAV_ITEMS.map((item) => {
          const active = item.exact ? pathname === item.href : pathname.startsWith(item.href);
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                'flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors',
                active
                  ? 'bg-brand-green text-white font-medium'
                  : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900',
              )}
            >
              <span className="text-base">{item.icon}</span>
              {item.label}
            </Link>
          );
        })}
      </nav>

      {/* Bottom links */}
      <div className="px-3 py-3 border-t border-gray-100 space-y-0.5">
        <Link
          href="/cabinet"
          className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm text-gray-500 hover:bg-gray-50 hover:text-gray-900 transition-colors"
        >
          <span className="text-base">🏠</span>
          Личный кабинет
        </Link>
        <button
          onClick={() => signOut({ callbackUrl: '/' })}
          className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm text-gray-500 hover:bg-gray-50 hover:text-gray-900 w-full transition-colors"
        >
          <span className="text-base">🚪</span>
          Выйти
        </button>
      </div>
    </aside>
  );
}
