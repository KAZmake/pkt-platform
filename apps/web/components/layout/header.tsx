'use client';

import Link from 'next/link';
import { useState } from 'react';
import { Button } from '@/components/ui/button';

const NAV_LINKS = [
  { href: '/about', label: 'О компании' },
  { href: '/programs', label: 'Программы' },
  { href: '/map', label: 'Карта' },
  { href: '/projects', label: 'Проекты' },
  { href: '/news', label: 'Новости' },
  { href: '/contacts', label: 'Контакты' },
  { href: '/faq', label: 'FAQ' },
];

const LANGUAGES = [
  { code: 'ru', label: 'Рус' },
  { code: 'kz', label: 'Қаз' },
  { code: 'en', label: 'Eng' },
];

export function Header() {
  const [menuOpen, setMenuOpen] = useState(false);
  const [lang, setLang] = useState('ru');

  return (
    <header className="sticky top-0 z-50 bg-white border-b border-gray-200 shadow-sm">
      {/* Top bar */}
      <div className="bg-brand-green text-white text-xs">
        <div className="container-page flex items-center justify-between h-8">
          <a
            href="tel:+77112000000"
            className="hover:text-brand-gold transition-colors font-medium tracking-wide"
          >
            8 (7112) 00-00-00 — Call-центр
          </a>
          <div className="flex items-center gap-3">
            {/* Language switcher */}
            <div className="flex items-center gap-1">
              {LANGUAGES.map((l) => (
                <button
                  key={l.code}
                  onClick={() => setLang(l.code)}
                  className={`px-1.5 py-0.5 rounded text-xs transition-colors ${
                    lang === l.code ? 'bg-white/20 font-semibold' : 'hover:bg-white/10'
                  }`}
                >
                  {l.label}
                </button>
              ))}
            </div>
            {/* Accessibility */}
            <button
              aria-label="Версия для слабовидящих"
              className="hover:text-brand-gold transition-colors text-base leading-none"
              title="Версия для слабовидящих"
            >
              👁
            </button>
          </div>
        </div>
      </div>

      {/* Main nav */}
      <div className="container-page">
        <div className="flex items-center justify-between h-16 gap-4">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2 shrink-0">
            <span className="inline-flex h-9 w-9 items-center justify-center rounded-lg bg-brand-green text-white font-bold text-lg">
              ПКТ
            </span>
            <span className="hidden sm:block text-sm font-semibold text-gray-800 leading-tight">
              Первое кредитное
              <br />
              товарищество
            </span>
          </Link>

          {/* Desktop nav */}
          <nav className="hidden lg:flex items-center gap-1">
            {NAV_LINKS.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className="px-3 py-2 text-sm text-gray-700 rounded-md hover:bg-brand-green-50 hover:text-brand-green transition-colors"
              >
                {link.label}
              </Link>
            ))}
          </nav>

          {/* CTA */}
          <div className="hidden md:flex items-center gap-2 shrink-0">
            <Link href="/cabinet">
              <Button variant="outline" size="sm">
                Личный кабинет
              </Button>
            </Link>
            <Link href="/apply">
              <Button size="sm">Подать заявку</Button>
            </Link>
          </div>

          {/* Burger */}
          <button
            className="lg:hidden p-2 rounded-md hover:bg-gray-100"
            onClick={() => setMenuOpen(!menuOpen)}
            aria-label={menuOpen ? 'Закрыть меню' : 'Открыть меню'}
          >
            <span className="block w-5 h-0.5 bg-gray-700 mb-1" />
            <span className="block w-5 h-0.5 bg-gray-700 mb-1" />
            <span className="block w-5 h-0.5 bg-gray-700" />
          </button>
        </div>
      </div>

      {/* Mobile menu */}
      {menuOpen && (
        <div className="lg:hidden border-t border-gray-200 bg-white">
          <nav className="container-page py-3 flex flex-col gap-1">
            {NAV_LINKS.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className="px-3 py-2 text-sm text-gray-700 rounded-md hover:bg-brand-green-50 hover:text-brand-green"
                onClick={() => setMenuOpen(false)}
              >
                {link.label}
              </Link>
            ))}
            <div className="flex gap-2 mt-3 pt-3 border-t border-gray-100">
              <Link href="/cabinet" className="flex-1">
                <Button variant="outline" size="sm" className="w-full">
                  Личный кабинет
                </Button>
              </Link>
              <Link href="/apply" className="flex-1">
                <Button size="sm" className="w-full">
                  Подать заявку
                </Button>
              </Link>
            </div>
          </nav>
        </div>
      )}
    </header>
  );
}
