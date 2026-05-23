import { NextResponse } from 'next/server';

/**
 * GeoJSON points for the employee portfolio heatmap.
 * Each point has a `weight` property (loan amount normalised 0–1).
 * Production: joins loans table with borrower coordinates.
 */

const LOANS = Array.from({ length: 80 }, (_, i) => {
  const seed = (i * 1664525 + 1013904223) >>> 0;
  const lng = 49.3 + (seed % 1200) / 240;
  const lat = 48.8 + ((seed >> 12) % 900) / 250;
  const amount = 200_000 + (seed % 4_800_000);
  return { lng, lat, weight: Math.min(amount / 5_000_000, 1) };
});

const geojson = {
  type: 'FeatureCollection',
  features: LOANS.map((l, i) => ({
    type: 'Feature',
    id: i,
    properties: { weight: l.weight },
    geometry: { type: 'Point', coordinates: [l.lng, l.lat] },
  })),
};

export function GET() {
  return NextResponse.json(geojson, {
    headers: { 'Cache-Control': 'private, max-age=300' },
  });
}
