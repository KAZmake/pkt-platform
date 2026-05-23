'use client';

import { useState } from 'react';
import { cn } from '@/lib/utils';

interface AccordionItem {
  id: string;
  question: string;
  answer: string;
}

interface AccordionProps {
  items: AccordionItem[];
  className?: string;
}

export function Accordion({ items, className }: AccordionProps) {
  const [open, setOpen] = useState<string | null>(null);

  return (
    <div className={cn('divide-y divide-gray-200 border-y border-gray-200', className)}>
      {items.map((item) => {
        const isOpen = open === item.id;
        return (
          <div key={item.id}>
            <button
              className="flex w-full items-center justify-between gap-4 py-4 text-left text-sm font-medium text-gray-900 hover:text-brand-green transition-colors"
              onClick={() => setOpen(isOpen ? null : item.id)}
              aria-expanded={isOpen}
            >
              <span>{item.question}</span>
              <span
                className={cn(
                  'shrink-0 text-brand-green transition-transform duration-200',
                  isOpen && 'rotate-45',
                )}
              >
                ✕
              </span>
            </button>
            {isOpen && (
              <div className="pb-4 text-sm text-gray-600 leading-relaxed">{item.answer}</div>
            )}
          </div>
        );
      })}
    </div>
  );
}
