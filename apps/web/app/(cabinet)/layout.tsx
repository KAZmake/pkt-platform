import { auth } from '@/auth';
import { redirect } from 'next/navigation';
import { CabinetSidebar } from '@/components/cabinet/cabinet-sidebar';

export default async function CabinetLayout({ children }: { children: React.ReactNode }) {
  const session = await auth();
  if (!session) redirect('/login?callbackUrl=/cabinet');

  return (
    <div className="min-h-screen bg-gray-50 flex">
      <CabinetSidebar userName={session.user?.name ?? 'Пользователь'} role={session.role} />
      <main className="flex-1 p-6 lg:p-8 overflow-auto min-w-0">{children}</main>
    </div>
  );
}
