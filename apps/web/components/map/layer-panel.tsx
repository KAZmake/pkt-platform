'use client';

import { useState } from 'react';
import type { LayerToggle } from '@/lib/map/types';
import { cn } from '@/lib/utils';

interface LayerPanelProps {
  layers: LayerToggle[];
  onToggle: (id: LayerToggle['id']) => void;
  isEmployee: boolean;
}

export function LayerPanel({ layers, onToggle, isEmployee }: LayerPanelProps) {
  const [open, setOpen] = useState(true);

  const publicLayers = layers.filter((l) => !l.employeeOnly);
  const employeeLayers = layers.filter((l) => l.employeeOnly);

  return (
    <div className="absolute top-3 right-3 z-10">
      {/* Toggle button */}
      <button
        onClick={() => setOpen((v) => !v)}
        className="flex items-center gap-2 bg-white rounded-lg shadow-md px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 border border-gray-200 mb-2"
        aria-label="Слои карты"
      >
        <span>🗂</span>
        <span className="hidden sm:inline">Слои</span>
        <span className="text-gray-400 text-xs">{open ? '▲' : '▼'}</span>
      </button>

      {open && (
        <div className="bg-white rounded-xl shadow-lg border border-gray-200 p-4 w-56 space-y-4">
          {/* Public layers */}
          <div>
            <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2">
              Публичные слои
            </p>
            <div className="space-y-1">
              {publicLayers.map((layer) => (
                <LayerRow key={layer.id} layer={layer} onToggle={onToggle} />
              ))}
            </div>
          </div>

          {/* Employee layers */}
          {isEmployee && employeeLayers.length > 0 && (
            <div>
              <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2">
                Для сотрудников
              </p>
              <div className="space-y-1">
                {employeeLayers.map((layer) => (
                  <LayerRow key={layer.id} layer={layer} onToggle={onToggle} />
                ))}
              </div>
            </div>
          )}

          {!isEmployee && (
            <p className="text-xs text-gray-400 border-t border-gray-100 pt-3">
              Войдите как сотрудник для доступа к слоям портфеля.
            </p>
          )}
        </div>
      )}
    </div>
  );
}

function LayerRow({
  layer,
  onToggle,
}: {
  layer: LayerToggle;
  onToggle: (id: LayerToggle['id']) => void;
}) {
  return (
    <label className="flex items-center gap-2 cursor-pointer group py-0.5">
      <div
        className={cn(
          'relative h-4 w-7 rounded-full transition-colors',
          layer.enabled ? 'bg-brand-green' : 'bg-gray-200',
        )}
        onClick={() => onToggle(layer.id)}
      >
        <span
          className={cn(
            'absolute top-0.5 h-3 w-3 rounded-full bg-white shadow transition-transform',
            layer.enabled ? 'translate-x-3.5' : 'translate-x-0.5',
          )}
        />
      </div>
      <span className="text-sm text-gray-700 group-hover:text-gray-900">
        {layer.icon} {layer.label}
      </span>
    </label>
  );
}
