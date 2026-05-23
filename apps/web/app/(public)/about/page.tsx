import type { Metadata } from 'next';
import { getTeam } from '@/lib/directus';

export const metadata: Metadata = {
  title: 'О компании',
  description: 'История и миссия ТОО «Первое кредитное товарищество». Команда профессионалов.',
};

export default async function AboutPage() {
  const team = await getTeam().catch(() => []);

  return (
    <div className="container-page py-12 space-y-16">
      {/* Mission */}
      <section className="max-w-3xl">
        <p className="text-brand-gold text-xs font-semibold uppercase tracking-widest mb-3">
          О нас
        </p>
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-4">
          Поддерживаем агросектор ЗКО с 2010 года
        </h1>
        <div className="prose prose-sm text-gray-600 space-y-4">
          <p>
            ТОО «Первое кредитное товарищество» — микрофинансовая организация, специализирующаяся на
            кредитовании малых и средних агропредприятий Западно-Казахстанской области.
          </p>
          <p>
            За 14 лет работы мы выдали более 5 млрд тенге займов, поддержав свыше 500 фермерских
            хозяйств. Наша миссия — сделать финансирование доступным для каждого аграрника ЗКО.
          </p>
        </div>
      </section>

      {/* Values */}
      <section>
        <h2 className="text-xl font-bold text-gray-900 mb-6">Наши принципы</h2>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
          {[
            {
              icon: '🌱',
              title: 'Доступность',
              text: 'Минимальный пакет документов, выезд специалиста к заёмщику, рассмотрение за 3 дня.',
            },
            {
              icon: '🤝',
              title: 'Партнёрство',
              text: 'Сопровождаем заёмщика на всех этапах — от заявки до погашения.',
            },
            {
              icon: '📊',
              title: 'Прозрачность',
              text: 'Чёткие условия, никаких скрытых комиссий, личный кабинет с историей.',
            },
          ].map(({ icon, title, text }) => (
            <div key={title} className="rounded-xl border border-gray-200 bg-white p-6">
              <div className="text-3xl mb-3">{icon}</div>
              <h3 className="font-semibold text-gray-900 mb-2">{title}</h3>
              <p className="text-sm text-gray-500 leading-relaxed">{text}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Team */}
      {team.length > 0 && (
        <section>
          <h2 className="text-xl font-bold text-gray-900 mb-6">Команда</h2>
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-6">
            {team.map((member) => (
              <div key={member.id} className="text-center">
                <div className="mx-auto mb-3 h-20 w-20 rounded-full bg-brand-green-100 overflow-hidden">
                  {member.photo ? (
                    <img
                      src={`${process.env.DIRECTUS_URL ?? 'http://localhost:8055'}/assets/${member.photo}`}
                      alt={member.name}
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <div className="flex h-full items-center justify-center text-xl font-bold text-brand-green">
                      {member.name[0]}
                    </div>
                  )}
                </div>
                <p className="font-medium text-gray-900 text-sm">{member.name}</p>
                <p className="text-xs text-gray-400 mt-0.5">{member.position}</p>
              </div>
            ))}
          </div>
        </section>
      )}
    </div>
  );
}
