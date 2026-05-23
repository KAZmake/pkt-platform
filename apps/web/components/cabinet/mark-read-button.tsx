'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';

export function MarkReadButton() {
  const [done, setDone] = useState(false);

  function handleClick() {
    // In production: await apiClient.post('/cabinet/notifications/read-all', {}, { token })
    setDone(true);
  }

  return (
    <Button variant="outline" size="sm" onClick={handleClick} disabled={done}>
      {done ? 'Все прочитаны' : 'Отметить все прочитанными'}
    </Button>
  );
}
