'use client';

import 'maplibre-gl/dist/maplibre-gl.css';
import maplibregl from 'maplibre-gl';
import { useEffect, useRef, useState, useCallback } from 'react';
import ReactDOM from 'react-dom/client';
import {
  MAP_CENTER,
  MAP_ZOOM,
  MAP_MIN_ZOOM,
  MAP_MAX_ZOOM,
  MAP_STYLE,
  LAYER,
  SOURCE,
  ACTIVITY_COLORS,
  LOAN_STATUS_COLORS,
} from '@/lib/map/config';
import type { LayerToggle, FarmProperties } from '@/lib/map/types';
import { LayerPanel } from './layer-panel';
import { FarmPopup } from './farm-popup';

const DEFAULT_LAYERS: LayerToggle[] = [
  { id: 'districts', label: 'Районы ЗКО', icon: '🗺', employeeOnly: false, enabled: true },
  { id: 'farms', label: 'Хозяйства', icon: '🌾', employeeOnly: false, enabled: true },
  { id: 'activity', label: 'Вид деятельности', icon: '🏷', employeeOnly: false, enabled: false },
  { id: 'coverage', label: 'Зона охвата', icon: '📍', employeeOnly: false, enabled: false },
  {
    id: 'portfolio',
    label: 'Тепловая карта портфеля',
    icon: '🔥',
    employeeOnly: true,
    enabled: false,
  },
  { id: 'loan-status', label: 'Статус займов', icon: '📊', employeeOnly: true, enabled: false },
];

interface MapViewProps {
  role?: string;
}

export function MapView({ role = 'public' }: MapViewProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const mapRef = useRef<maplibregl.Map | null>(null);
  const popupRef = useRef<maplibregl.Popup | null>(null);
  const popupRootRef = useRef<ReactDOM.Root | null>(null);
  const [layers, setLayers] = useState<LayerToggle[]>(DEFAULT_LAYERS);
  const [mapReady, setMapReady] = useState(false);

  const isEmployee = ['employee', 'expert', 'admin'].includes(role);

  // ── Initialise map ──────────────────────────────────────────────────────────
  useEffect(() => {
    if (!containerRef.current || mapRef.current) return;

    const map = new maplibregl.Map({
      container: containerRef.current,
      style: MAP_STYLE as maplibregl.StyleSpecification,
      center: MAP_CENTER,
      zoom: MAP_ZOOM,
      minZoom: MAP_MIN_ZOOM,
      maxZoom: MAP_MAX_ZOOM,
    });

    map.addControl(new maplibregl.NavigationControl(), 'bottom-right');
    map.addControl(new maplibregl.ScaleControl({ unit: 'metric' }), 'bottom-left');

    map.on('load', async () => {
      // ── Districts ──────────────────────────────────────────────────────────
      const [districtsData, farmsData, portfolioData] = await Promise.all([
        fetch('/api/map/districts').then((r) => r.json()),
        fetch('/api/map/farms').then((r) => r.json()),
        isEmployee ? fetch('/api/map/portfolio').then((r) => r.json()) : Promise.resolve(null),
      ]);

      // Districts source
      map.addSource(SOURCE.DISTRICTS, { type: 'geojson', data: districtsData });

      map.addLayer({
        id: LAYER.DISTRICTS_FILL,
        type: 'fill',
        source: SOURCE.DISTRICTS,
        paint: {
          'fill-color': '#1a5c36',
          'fill-opacity': 0.08,
        },
      });
      map.addLayer({
        id: LAYER.DISTRICTS_LINE,
        type: 'line',
        source: SOURCE.DISTRICTS,
        paint: { 'line-color': '#1a5c36', 'line-width': 1.5, 'line-opacity': 0.6 },
      });
      map.addLayer({
        id: LAYER.DISTRICTS_LABEL,
        type: 'symbol',
        source: SOURCE.DISTRICTS,
        layout: {
          'text-field': ['get', 'name'],
          'text-size': 11,
          'text-font': ['Open Sans Regular'],
        },
        paint: { 'text-color': '#0f3b22', 'text-halo-color': '#fff', 'text-halo-width': 1 },
      });

      // Farms source with clustering
      map.addSource(SOURCE.FARMS, {
        type: 'geojson',
        data: farmsData,
        cluster: true,
        clusterMaxZoom: 11,
        clusterRadius: 50,
      });

      // Cluster circle
      map.addLayer({
        id: LAYER.FARMS_CLUSTER,
        type: 'circle',
        source: SOURCE.FARMS,
        filter: ['has', 'point_count'],
        paint: {
          'circle-color': ['step', ['get', 'point_count'], '#2d7a50', 10, '#1a5c36', 30, '#0f3b22'],
          'circle-radius': ['step', ['get', 'point_count'], 18, 10, 24, 30, 30],
          'circle-stroke-width': 2,
          'circle-stroke-color': '#fff',
        },
      });
      map.addLayer({
        id: LAYER.FARMS_CLUSTER_COUNT,
        type: 'symbol',
        source: SOURCE.FARMS,
        filter: ['has', 'point_count'],
        layout: {
          'text-field': '{point_count_abbreviated}',
          'text-size': 12,
          'text-font': ['Open Sans Bold'],
        },
        paint: { 'text-color': '#fff' },
      });

      // Individual farm circles
      map.addLayer({
        id: LAYER.FARMS_UNCLUSTERED,
        type: 'circle',
        source: SOURCE.FARMS,
        filter: ['!', ['has', 'point_count']],
        paint: {
          'circle-color': [
            'match',
            ['get', 'activity_type'],
            'crop_farming',
            ACTIVITY_COLORS.crop_farming,
            'livestock',
            ACTIVITY_COLORS.livestock,
            ACTIVITY_COLORS.mixed,
          ],
          'circle-radius': 7,
          'circle-stroke-width': 2,
          'circle-stroke-color': '#fff',
        },
      });

      // Portfolio heatmap (employee-only)
      if (portfolioData) {
        map.addSource(SOURCE.PORTFOLIO, { type: 'geojson', data: portfolioData });
        map.addLayer({
          id: LAYER.PORTFOLIO_HEAT,
          type: 'heatmap',
          source: SOURCE.PORTFOLIO,
          paint: {
            'heatmap-weight': ['interpolate', ['linear'], ['get', 'weight'], 0, 0, 1, 1],
            'heatmap-intensity': ['interpolate', ['linear'], ['zoom'], 5, 1, 12, 3],
            'heatmap-color': [
              'interpolate',
              ['linear'],
              ['heatmap-density'],
              0,
              'rgba(0,0,0,0)',
              0.2,
              'rgba(200,146,26,0.3)',
              0.6,
              'rgba(200,146,26,0.7)',
              1,
              'rgba(200,146,26,1)',
            ],
            'heatmap-radius': ['interpolate', ['linear'], ['zoom'], 5, 20, 12, 40],
            'heatmap-opacity': 0.7,
          },
          layout: { visibility: 'none' },
        });

        // Loan-status coloured circles
        map.addLayer({
          id: LAYER.LOAN_STATUS,
          type: 'circle',
          source: SOURCE.FARMS,
          filter: ['!', ['has', 'point_count']],
          paint: {
            'circle-color': [
              'match',
              ['get', 'loan_status'],
              'active',
              LOAN_STATUS_COLORS.active,
              'overdue',
              LOAN_STATUS_COLORS.overdue,
              LOAN_STATUS_COLORS.closed,
            ],
            'circle-radius': 8,
            'circle-stroke-width': 2,
            'circle-stroke-color': '#fff',
          },
          layout: { visibility: 'none' },
        });
      }

      // ── Interactions ──────────────────────────────────────────────────────
      // Zoom into cluster on click
      map.on('click', LAYER.FARMS_CLUSTER, async (e) => {
        const features = map.queryRenderedFeatures(e.point, {
          layers: [LAYER.FARMS_CLUSTER],
        });
        if (!features[0]) return;
        const clusterId = features[0].properties?.cluster_id as number;
        const src = map.getSource(SOURCE.FARMS) as maplibregl.GeoJSONSource;
        const zoom = await src.getClusterExpansionZoom(clusterId);
        map.easeTo({
          center: (features[0].geometry as GeoJSON.Point).coordinates as [number, number],
          zoom,
        });
      });

      // Farm popup
      map.on('click', LAYER.FARMS_UNCLUSTERED, (e) => {
        const feature = e.features?.[0];
        if (!feature) return;
        const props = feature.properties as FarmProperties;
        const coords = (feature.geometry as GeoJSON.Point).coordinates as [number, number];
        showPopup(map, coords, props, isEmployee);
      });

      map.on('mouseenter', LAYER.FARMS_CLUSTER, () => {
        map.getCanvas().style.cursor = 'pointer';
      });
      map.on('mouseleave', LAYER.FARMS_CLUSTER, () => {
        map.getCanvas().style.cursor = '';
      });
      map.on('mouseenter', LAYER.FARMS_UNCLUSTERED, () => {
        map.getCanvas().style.cursor = 'pointer';
      });
      map.on('mouseleave', LAYER.FARMS_UNCLUSTERED, () => {
        map.getCanvas().style.cursor = '';
      });

      setMapReady(true);
    });

    mapRef.current = map;
    return () => {
      map.remove();
      mapRef.current = null;
    };
  }, []); // intentionally empty — map init runs once

  // ── Layer visibility ────────────────────────────────────────────────────────
  useEffect(() => {
    const map = mapRef.current;
    if (!map || !mapReady) return;

    function vis(layerId: string, enabled: boolean) {
      if (map!.getLayer(layerId)) {
        map!.setLayoutProperty(layerId, 'visibility', enabled ? 'visible' : 'none');
      }
    }

    const state = Object.fromEntries(layers.map((l) => [l.id, l.enabled]));

    vis(LAYER.DISTRICTS_FILL, state.districts);
    vis(LAYER.DISTRICTS_LINE, state.districts);
    vis(LAYER.DISTRICTS_LABEL, state.districts);
    vis(LAYER.FARMS_CLUSTER, state.farms);
    vis(LAYER.FARMS_CLUSTER_COUNT, state.farms);
    vis(LAYER.FARMS_UNCLUSTERED, state.farms && !state.activity && !state['loan-status']);
    vis(LAYER.ACTIVITY_HEAT, state.activity);
    vis(LAYER.PORTFOLIO_HEAT, state.portfolio);
    vis(LAYER.LOAN_STATUS, state['loan-status']);
  }, [layers, mapReady]);

  const handleToggle = useCallback((id: LayerToggle['id']) => {
    setLayers((prev) => prev.map((l) => (l.id === id ? { ...l, enabled: !l.enabled } : l)));
  }, []);

  function showPopup(
    map: maplibregl.Map,
    coords: [number, number],
    props: FarmProperties,
    employee: boolean,
  ) {
    popupRef.current?.remove();
    const el = document.createElement('div');
    const root = ReactDOM.createRoot(el);
    popupRootRef.current = root;

    const popup = new maplibregl.Popup({ closeButton: false, maxWidth: '280px' })
      .setLngLat(coords)
      .setDOMContent(el)
      .addTo(map);

    popupRef.current = popup;

    root.render(<FarmPopup farm={props} isEmployee={employee} onClose={() => popup.remove()} />);
  }

  return (
    <div className="relative h-full w-full">
      <div ref={containerRef} className="h-full w-full" />
      <LayerPanel layers={layers} onToggle={handleToggle} isEmployee={isEmployee} />
      {!mapReady && (
        <div className="absolute inset-0 flex items-center justify-center bg-gray-100">
          <span className="h-8 w-8 animate-spin rounded-full border-2 border-brand-green border-t-transparent" />
        </div>
      )}
    </div>
  );
}
