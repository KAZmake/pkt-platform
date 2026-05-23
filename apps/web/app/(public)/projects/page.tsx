import type { Metadata } from 'next';
import { getProjects, assetUrl } from '@/lib/directus';
import { ProjectsClient } from '@/components/projects/projects-client';

export const metadata: Metadata = {
  title: 'Проекты',
  description: 'Реализованные проекты заёмщиков ПКТ в агросекторе ЗКО.',
};

export default async function ProjectsPage() {
  const projects = await getProjects().catch(() => []);

  return (
    <div className="container-page py-12">
      <div className="mb-8">
        <p className="text-brand-gold text-xs font-semibold uppercase tracking-widest mb-2">
          Истории успеха
        </p>
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Проекты наших заёмщиков</h1>
      </div>
      <ProjectsClient projects={projects} assetUrl={assetUrl} />
    </div>
  );
}
