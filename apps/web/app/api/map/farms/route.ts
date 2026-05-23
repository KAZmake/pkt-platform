import { NextResponse } from 'next/server';

/**
 * GeoJSON FeatureCollection of farm/enterprise point locations.
 * Production: fetched from core-api PostGIS (farms table with geometry).
 */

const ACTIVITY = ['crop_farming', 'livestock', 'mixed'] as const;
const DISTRICTS = [
  'Зеленовский',
  'Акжаикский',
  'Бурлинский',
  'Таскалинский',
  'Теректинский',
  'Сырымский',
  'Чингирлауский',
  'Каратобинский',
];

// Deterministic pseudo-random spread around Uralsk (51.23°N, 51.37°E)
const FARMS = Array.from({ length: 60 }, (_, i) => {
  const seed = (i * 2654435761) >>> 0;
  const lng = 49.5 + (seed % 1000) / 200;
  const lat = 49.0 + ((seed >> 10) % 1000) / 300;
  return {
    id: `farm-${i + 1}`,
    name: `Хозяйство №${i + 1}`,
    borrower_id: `borrower-${(i % 20) + 1}`,
    activity_type: ACTIVITY[i % 3],
    district: DISTRICTS[i % DISTRICTS.length],
    area_ha: 50 + ((i * 37) % 950),
    loan_status: ['active', 'active', 'active', 'overdue', 'closed'][i % 5] as string,
    loan_amount: 500_000 + ((i * 123_456) % 4_500_000),
    lng,
    lat,
  };
});

const geojson = {
  type: 'FeatureCollection',
  features: FARMS.map((f) => ({
    type: 'Feature',
    properties: {
      id: f.id,
      name: f.name,
      borrower_id: f.borrower_id,
      activity_type: f.activity_type,
      district: f.district,
      area_ha: f.area_ha,
      loan_status: f.loan_status,
      loan_amount: f.loan_amount,
    },
    geometry: { type: 'Point', coordinates: [f.lng, f.lat] },
  })),
};

export function GET() {
  return NextResponse.json(geojson, {
    headers: { 'Cache-Control': 'public, max-age=60' },
  });
}
