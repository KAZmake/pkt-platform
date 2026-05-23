/**
 * Decorative SVG outline of Западно-Казахстанская область (ZKO).
 * Districts are stylised — not geographically precise.
 */
export function ZkoMapPreview() {
  return (
    <section className="bg-brand-green-50 py-16">
      <div className="container-page flex flex-col lg:flex-row items-center gap-12">
        {/* Text */}
        <div className="flex-1 max-w-lg">
          <p className="text-brand-gold text-xs font-semibold uppercase tracking-widest mb-3">
            География присутствия
          </p>
          <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-4">
            Работаем во всех районах ЗКО
          </h2>
          <p className="text-gray-600 text-sm leading-relaxed mb-6">
            Охватываем 12 районов Западно-Казахстанской области. Наши специалисты выезжают к
            заёмщикам для оценки залогового имущества и сопровождения проектов.
          </p>
          <ul className="grid grid-cols-2 gap-y-1.5 text-sm text-gray-700">
            {[
              'Зеленовский',
              'Акжаикский',
              'Бурлинский',
              'Казталовский',
              'Каратобинский',
              'Сырымский',
              'Таскалинский',
              'Теректинский',
              'Чингирлауский',
              'Жангалинский',
              'Жаныбекский',
              'Бокейординский',
            ].map((d) => (
              <li key={d} className="flex items-center gap-1.5">
                <span className="h-1.5 w-1.5 rounded-full bg-brand-green" />
                {d}
              </li>
            ))}
          </ul>
        </div>

        {/* Stylised SVG map */}
        <div className="flex-1 flex justify-center">
          <svg
            viewBox="0 0 300 260"
            className="w-full max-w-xs lg:max-w-sm drop-shadow-md"
            aria-label="Карта ЗКО"
          >
            {/* Outer region boundary */}
            <path
              d="M80,30 L200,20 L270,80 L260,180 L200,240 L100,250 L30,190 L20,100 Z"
              fill="#edf7f1"
              stroke="#1a5c36"
              strokeWidth="2"
            />
            {/* District dividers (stylised) */}
            {[
              'M80,30 L150,130',
              'M200,20 L150,130',
              'M270,80 L150,130',
              'M260,180 L150,130',
              'M200,240 L150,130',
              'M100,250 L150,130',
              'M30,190 L150,130',
              'M20,100 L150,130',
            ].map((d, i) => (
              <path
                key={i}
                d={d}
                stroke="#2d7a50"
                strokeWidth="0.8"
                strokeDasharray="4 3"
                fill="none"
              />
            ))}
            {/* Uralsk city dot */}
            <circle cx="150" cy="130" r="7" fill="#c8921a" stroke="white" strokeWidth="2" />
            <text x="158" y="124" fontSize="9" fill="#0f3b22" fontWeight="600">
              Уральск
            </text>
          </svg>
        </div>
      </div>
    </section>
  );
}
