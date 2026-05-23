import { auth } from '@/auth';
import { redirect } from 'next/navigation';
import { EmployeeSidebar } from '@/components/employee/employee-sidebar';

const ALLOWED_ROLES = ['employee', 'expert', 'admin'];

export default async function EmployeeLayout({ children }: { children: React.ReactNode }) {
  const session = await auth();
  if (!session) redirect('/login');
  if (!ALLOWED_ROLES.includes(session.role)) redirect('/403');

  return (
    <div className="min-h-screen bg-gray-50 flex">
      <EmployeeSidebar userName={session.user?.name ?? 'Сотрудник'} role={session.role} />
      <main className="flex-1 p-6 lg:p-8 overflow-auto min-w-0">{children}</main>
    </div>
  );
}
