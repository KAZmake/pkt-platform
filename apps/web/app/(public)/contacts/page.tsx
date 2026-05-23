import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Контакты',
  description: 'Адрес, телефоны и email ТОО «Первое кредитное товарищество». Как к нам добраться.',
};

const CONTACTS = [
  { label: 'Адрес', value: 'г. Уральск, ул. Дина Нурпейисовой, д. 1, офис 201', icon: '📍' },
  { label: 'Call-центр', value: '8 (7112) 00-00-00', href: 'tel:+77112000000', icon: '📞' },
  { label: 'Email', value: 'info@pkt.kz', href: 'mailto:info@pkt.kz', icon: '✉️' },
  { label: 'Режим работы', value: 'Пн–Пт: 09:00 – 18:00', icon: '🕐' },
];

export default function ContactsPage() {
  return (
    <div className="container-page py-12">
      <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-8">Контакты</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
        {/* Contact details */}
        <div className="space-y-6">
          {CONTACTS.map(({ label, value, href, icon }) => (
            <div key={label} className="flex items-start gap-4">
              <span className="text-2xl">{icon}</span>
              <div>
                <p className="text-xs text-gray-400 mb-0.5">{label}</p>
                {href ? (
                  <a href={href} className="font-medium text-brand-green hover:underline">
                    {value}
                  </a>
                ) : (
                  <p className="font-medium text-gray-900">{value}</p>
                )}
              </div>
            </div>
          ))}

          <div className="mt-8 rounded-xl bg-brand-green-50 border border-brand-green-100 p-6">
            <h2 className="font-semibold text-gray-900 mb-2">Написать нам</h2>
            <p className="text-sm text-gray-500 mb-4">
              Для вопросов по кредитованию воспользуйтесь чатом на сайте или оформите обращение в
              личном кабинете.
            </p>
            <a
              href="/cabinet"
              className="inline-flex items-center gap-2 bg-brand-green text-white text-sm font-medium px-4 py-2 rounded-lg hover:bg-brand-green-dark transition-colors"
            >
              Личный кабинет
            </a>
          </div>
        </div>

        {/* Map (OpenStreetMap embed — Uralsk, WKO) */}
        <div className="rounded-xl overflow-hidden border border-gray-200 min-h-72">
          <iframe
            title="Карта расположения офиса ПКТ"
            width="100%"
            height="100%"
            style={{ minHeight: '360px', border: 0 }}
            loading="lazy"
            referrerPolicy="no-referrer-when-downgrade"
            src="https://www.openstreetmap.org/export/embed.html?bbox=51.15%2C51.19%2C51.45%2C51.39&layer=mapnik"
          />
        </div>
      </div>
    </div>
  );
}
