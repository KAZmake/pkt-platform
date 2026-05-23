'use client';

import { useState } from 'react';
import type { Project } from '@/lib/directus';

interface ProjectsClientProps {
  projects: Project[];
  assetUrl: (id: string) => string;
}

export function ProjectsClient({ projects, assetUrl }: ProjectsClientProps) {
  const [lightbox, setLightbox] = useState<{ src: string; title: string } | null>(null);

  if (projects.length === 0) {
    return <p className="text-gray-400 text-sm">Проекты появятся здесь.</p>;
  }

  const ACTIVITY_LABELS: Record<string, string> = {
    crop_farming: 'Растениеводство',
    livestock: 'Животноводство',
    mixed: 'Смешанное',
  };

  return (
    <>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {projects.map((project) => {
          const firstImage = project.images?.[0];
          return (
            <div
              key={project.id}
              className="rounded-xl border border-gray-200 bg-white overflow-hidden shadow-card hover:shadow-card-hover transition-shadow"
            >
              {firstImage && (
                <button
                  className="block w-full aspect-video bg-gray-100 overflow-hidden"
                  onClick={() => setLightbox({ src: assetUrl(firstImage), title: project.title })}
                >
                  <img
                    src={assetUrl(firstImage)}
                    alt={project.title}
                    className="h-full w-full object-cover hover:scale-105 transition-transform duration-300"
                  />
                </button>
              )}
              <div className="p-5">
                <div className="flex items-center gap-2 mb-2">
                  {project.activity_type && (
                    <span className="text-xs bg-brand-green-100 text-brand-green-700 rounded-full px-2 py-0.5">
                      {ACTIVITY_LABELS[project.activity_type] ?? project.activity_type}
                    </span>
                  )}
                  {project.year && <span className="text-xs text-gray-400">{project.year}</span>}
                </div>
                <h3 className="font-semibold text-gray-900 mb-1">{project.title}</h3>
                {project.borrower_name && (
                  <p className="text-xs text-gray-400 mb-2">{project.borrower_name}</p>
                )}
                <p className="text-sm text-gray-500 line-clamp-3">{project.description}</p>
                {(project.images ?? []).length > 1 && (
                  <div className="mt-3 flex gap-1.5">
                    {project.images.slice(1, 4).map((img) => (
                      <button
                        key={img}
                        onClick={() => setLightbox({ src: assetUrl(img), title: project.title })}
                        className="h-12 w-12 rounded overflow-hidden bg-gray-100 shrink-0"
                      >
                        <img src={assetUrl(img)} alt="" className="h-full w-full object-cover" />
                      </button>
                    ))}
                    {project.images.length > 4 && (
                      <span className="h-12 w-12 flex items-center justify-center rounded bg-gray-100 text-xs text-gray-400">
                        +{project.images.length - 4}
                      </span>
                    )}
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>

      {/* Lightbox */}
      {lightbox && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4"
          onClick={() => setLightbox(null)}
        >
          <div className="relative max-w-4xl w-full" onClick={(e) => e.stopPropagation()}>
            <button
              className="absolute -top-10 right-0 text-white text-2xl hover:text-gray-300"
              onClick={() => setLightbox(null)}
              aria-label="Закрыть"
            >
              ✕
            </button>
            <img
              src={lightbox.src}
              alt={lightbox.title}
              className="w-full rounded-xl max-h-[80vh] object-contain"
            />
            <p className="mt-3 text-center text-white text-sm">{lightbox.title}</p>
          </div>
        </div>
      )}
    </>
  );
}
