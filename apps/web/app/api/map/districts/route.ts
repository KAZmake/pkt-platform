import { NextResponse } from 'next/server';

/**
 * GeoJSON FeatureCollection of ZKO district boundaries.
 * In production this is served from PostGIS via core-api.
 * Coordinates are approximate rectangles for demo purposes.
 */
const DISTRICTS = [
  { id: 'zelenovsky', name: 'Зеленовский', center: [50.9, 51.5], sz: 0.7 },
  { id: 'akzhaiksky', name: 'Акжаикский', center: [50.2, 51.2], sz: 0.7 },
  { id: 'burlinsky', name: 'Бурлинский', center: [53.0, 51.3], sz: 0.7 },
  { id: 'kaztalovsky', name: 'Казталовский', center: [49.6, 49.5], sz: 0.8 },
  { id: 'karatobinsky', name: 'Каратобинский', center: [52.2, 50.2], sz: 0.7 },
  { id: 'syrymsky', name: 'Сырымский', center: [51.8, 50.9], sz: 0.7 },
  { id: 'taskalinsky', name: 'Таскалинский', center: [49.6, 51.0], sz: 0.7 },
  { id: 'teректinsky', name: 'Теректинский', center: [52.7, 51.2], sz: 0.7 },
  { id: 'chingirlausky', name: 'Чингирлауский', center: [50.5, 51.0], sz: 0.7 },
  { id: 'zhangalinsky', name: 'Жангалинский', center: [52.5, 49.3], sz: 0.7 },
  { id: 'zhanybекsky', name: 'Жаныбекский', center: [50.6, 49.3], sz: 0.7 },
  { id: 'bokeiordinsky', name: 'Бокейординский', center: [49.3, 49.2], sz: 0.8 },
];

function rect(lng: number, lat: number, sz: number) {
  const h = sz / 2;
  return [
    [lng - h, lat - h],
    [lng + h, lat - h],
    [lng + h, lat + h],
    [lng - h, lat + h],
    [lng - h, lat - h],
  ];
}

const geojson = {
  type: 'FeatureCollection',
  features: DISTRICTS.map((d) => ({
    type: 'Feature',
    properties: { id: d.id, name: d.name, farms_count: Math.floor(Math.random() * 40 + 5) },
    geometry: {
      type: 'Polygon',
      coordinates: [rect(d.center[0], d.center[1], d.sz)],
    },
  })),
};

export function GET() {
  return NextResponse.json(geojson, {
    headers: { 'Cache-Control': 'public, max-age=3600' },
  });
}
