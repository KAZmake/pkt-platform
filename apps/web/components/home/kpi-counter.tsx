'use client';

import { useEffect, useRef, useState } from 'react';

const KPI_ITEMS = [
  { value: 14, suffix: '+', label: 'лет на рынке' },
  { value: 500, suffix: '+', label: 'заёмщиков' },
  { value: 5, suffix: ' млрд ₸', label: 'выдано займов' },
  { value: 12, suffix: '', label: 'программ кредитования' },
];

function useCountUp(target: number, duration = 1500, active: boolean) {
  const [count, setCount] = useState(0);
  useEffect(() => {
    if (!active) return;
    let start = 0;
    const step = target / (duration / 16);
    const timer = setInterval(() => {
      start += step;
      if (start >= target) {
        setCount(target);
        clearInterval(timer);
      } else {
        setCount(Math.floor(start));
      }
    }, 16);
    return () => clearInterval(timer);
  }, [target, duration, active]);
  return count;
}

function KpiCard({ value, suffix, label }: (typeof KPI_ITEMS)[0]) {
  const ref = useRef<HTMLDivElement>(null);
  const [active, setActive] = useState(false);
  const count = useCountUp(value, 1200, active);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) setActive(true);
      },
      { threshold: 0.4 },
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  return (
    <div ref={ref} className="text-center">
      <p className="text-4xl sm:text-5xl font-bold text-brand-green">
        {count}
        <span className="text-brand-gold">{suffix}</span>
      </p>
      <p className="mt-2 text-sm text-gray-600">{label}</p>
    </div>
  );
}

export function KpiSection() {
  return (
    <section className="bg-white py-14 border-b border-gray-100">
      <div className="container-page grid grid-cols-2 lg:grid-cols-4 gap-8">
        {KPI_ITEMS.map((item) => (
          <KpiCard key={item.label} {...item} />
        ))}
      </div>
    </section>
  );
}
