import type { Metadata } from 'next';
import { LoanCalculator } from '@/components/calculator/loan-calculator';

export const metadata: Metadata = {
  title: 'Калькулятор займа',
  description: 'Рассчитайте ежемесячный платёж по аннуитетной или дифференцированной схеме.',
};

export default function CalculatorPage() {
  return (
    <div className="container-page py-12">
      <div className="mb-8">
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">Калькулятор займа</h1>
        <p className="text-gray-500 text-sm">
          Введите параметры займа, чтобы рассчитать график платежей. Итоговые условия уточняются при
          оформлении заявки.
        </p>
      </div>
      <LoanCalculator />
    </div>
  );
}
