import Link from 'next/link';
import { Button } from '@/components/ui/button';

export function Hero() {
  return (
    <section className="relative overflow-hidden bg-brand-green text-white">
      {/* Background pattern */}
      <div
        className="absolute inset-0 opacity-10"
        style={{
          backgroundImage:
            'radial-gradient(circle at 20% 50%, #2d7a50 0%, transparent 50%), radial-gradient(circle at 80% 20%, #c8921a 0%, transparent 40%)',
        }}
      />

      <div className="container-page relative py-20 lg:py-28">
        <div className="max-w-2xl">
          <p className="text-brand-gold text-sm font-semibold uppercase tracking-widest mb-4">
            Западно-Казахстанская область
          </p>
          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold leading-tight mb-6">
            Кредитование
            <br />
            <span className="text-brand-gold">агросектора</span> ЗКО
          </h1>
          <p className="text-white/80 text-lg leading-relaxed mb-8 max-w-xl">
            ТОО «Первое кредитное товарищество» — льготные займы для фермеров и сельскохозяйственных
            предприятий Западного Казахстана с 2010 года.
          </p>
          <div className="flex flex-wrap gap-3">
            <Link href="/apply">
              <Button size="lg" variant="secondary">
                Подать заявку
              </Button>
            </Link>
            <Link href="/programs">
              <Button
                size="lg"
                variant="outline"
                className="border-white text-white hover:bg-white/10"
              >
                Программы кредитования
              </Button>
            </Link>
          </div>
        </div>
      </div>

      {/* Decorative SVG wave */}
      <svg
        className="absolute bottom-0 left-0 right-0 text-white"
        viewBox="0 0 1440 40"
        fill="currentColor"
        preserveAspectRatio="none"
      >
        <path d="M0,40 C360,0 1080,0 1440,40 L1440,40 L0,40 Z" />
      </svg>
    </section>
  );
}
