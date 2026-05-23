/** ZKO region center (Uralsk/Oral) */
export const MAP_CENTER: [number, number] = [51.367, 51.233]; // [lng, lat]
export const MAP_ZOOM = 7;
export const MAP_MIN_ZOOM = 5;
export const MAP_MAX_ZOOM = 16;

/** Free OSM raster base-map style — no API key required. */
export const MAP_STYLE = {
  version: 8 as const,
  sources: {
    osm: {
      type: 'raster' as const,
      tiles: ['https://tile.openstreetmap.org/{z}/{x}/{y}.png'],
      tileSize: 256,
      attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
    },
  },
  layers: [
    {
      id: 'osm-tiles',
      type: 'raster' as const,
      source: 'osm',
      minzoom: 0,
      maxzoom: 22,
    },
  ],
};

// ─── Layer IDs ────────────────────────────────────────────────────────────────

export const LAYER = {
  DISTRICTS_FILL: 'districts-fill',
  DISTRICTS_LINE: 'districts-line',
  DISTRICTS_LABEL: 'districts-label',
  FARMS_CLUSTER: 'farms-cluster',
  FARMS_CLUSTER_COUNT: 'farms-cluster-count',
  FARMS_UNCLUSTERED: 'farms-unclustered',
  ACTIVITY_HEAT: 'activity-heat',
  COVERAGE_FILL: 'coverage-fill',
  PORTFOLIO_HEAT: 'portfolio-heat',
  LOAN_STATUS: 'loan-status',
} as const;

export const SOURCE = {
  DISTRICTS: 'districts',
  FARMS: 'farms',
  PORTFOLIO: 'portfolio',
} as const;

// ─── Colour mapping ───────────────────────────────────────────────────────────

export const ACTIVITY_COLORS: Record<string, string> = {
  crop_farming: '#2d7a50',
  livestock: '#c8921a',
  mixed: '#3b82f6',
};

export const LOAN_STATUS_COLORS: Record<string, string> = {
  active: '#10b981',
  overdue: '#ef4444',
  closed: '#6b7280',
};
