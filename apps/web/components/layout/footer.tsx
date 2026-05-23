import Link from 'next/link';

const FOOTER_LINKS = {
  company: [
    { href: '/about', label: 'О компании' },
    { href: '/projects', label: 'Проекты' },
    { href: '/news', label: 'Новости' },
    { href: '/contacts', label: 'Контакты' },
  ],
  services: [
    { href: '/programs', label: 'Программы кредитования' },
    { href: '/map', label: 'Интерактивная карта' },
    { href: '/apply', label: 'Подать заявку' },
    { href: '/calculator', label: 'Калькулятор' },
  ],
  support: [
    { href: '/faq', label: 'FAQ' },
    { href: '/cabinet', label: 'Личный кабинет' },
    { href: '/tickets', label: 'Обращения' },
  ],
};

export function Footer() {
  return (
    <footer className="bg-gray-900 text-gray-400 mt-auto">
      <div className="container-page py-12">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
          {/* Brand */}
          <div>
            <div className="flex items-center gap-2 mb-3">
              <span className="inline-flex h-8 w-8 items-center justify-center rounded-lg bg-brand-green text-white font-bold text-sm">
                ПКТ
              </span>
              <span className="text-white font-semibold text-sm">
                Первое кредитное
                <br />
                товарищество
              </span>
            </div>
            <p className="text-sm leading-relaxed">
              Кредитование агросектора ЗКО. Поддержка фермеров с 2010 года.
            </p>
            <p className="mt-3 text-sm font-medium text-white">
              <a href="tel:+77112000000" className="hover:text-brand-gold transition-colors">
                8 (7112) 00-00-00
              </a>
            </p>
          </div>

          {/* Links */}
          <div>
            <h4 className="text-white text-sm font-semibold mb-3">Компания</h4>
            <ul className="space-y-2">
              {FOOTER_LINKS.company.map((l) => (
                <li key={l.href}>
                  <Link href={l.href} className="text-sm hover:text-white transition-colors">
                    {l.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h4 className="text-white text-sm font-semibold mb-3">Услуги</h4>
            <ul className="space-y-2">
              {FOOTER_LINKS.services.map((l) => (
                <li key={l.href}>
                  <Link href={l.href} className="text-sm hover:text-white transition-colors">
                    {l.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h4 className="text-white text-sm font-semibold mb-3">Поддержка</h4>
            <ul className="space-y-2">
              {FOOTER_LINKS.support.map((l) => (
                <li key={l.href}>
                  <Link href={l.href} className="text-sm hover:text-white transition-colors">
                    {l.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <div className="mt-10 pt-6 border-t border-gray-800 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs">
          <p>
            © {new Date().getFullYear()} ТОО «Первое кредитное товарищество». Все права защищены.
          </p>
          <p>ЗКО, г. Уральск</p>
        </div>
      </div>
    </footer>
  );
}
