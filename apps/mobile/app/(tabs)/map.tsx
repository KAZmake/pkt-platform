import { useCallback, useEffect, useRef, useState } from 'react';
import { Text, View } from 'react-native';
import MapView, { Circle, Marker, type Region } from 'react-native-maps';
import Supercluster from 'supercluster';
import { useAuth } from '~/lib/auth/context';
import { FARM_FEATURES, type FarmProperties } from '~/lib/map/data';
import { ClusterMarker, FarmMarker } from '~/components/map/FarmMarker';
import { FarmPopup } from '~/components/map/FarmPopup';
import { LayerPanel, type MapLayer } from '~/components/map/LayerPanel';

const INITIAL_REGION: Region = {
  latitude: 50.8,
  longitude: 51.5,
  latitudeDelta: 4.5,
  longitudeDelta: 4.5,
};

const INITIAL_LAYERS: MapLayer[] = [
  { id: 'farms', label: 'Хозяйства', icon: 'business-outline', active: true },
  { id: 'coverage', label: 'Покрытие КТ', icon: 'radio-outline', active: false },
  {
    id: 'portfolio',
    label: 'Портфель займов',
    icon: 'stats-chart-outline',
    active: false,
    employeeOnly: true,
  },
  {
    id: 'status',
    label: 'Статус займов',
    icon: 'checkmark-circle-outline',
    active: false,
    employeeOnly: true,
  },
];

function getZoomLevel(longitudeDelta: number): number {
  return Math.min(Math.round(Math.log(360 / longitudeDelta) / Math.LN2), 20);
}

function getBBox(region: Region): [number, number, number, number] {
  return [
    region.longitude - region.longitudeDelta / 2,
    region.latitude - region.latitudeDelta / 2,
    region.longitude + region.longitudeDelta / 2,
    region.latitude + region.latitudeDelta / 2,
  ];
}

type AnyCluster = ReturnType<Supercluster<FarmProperties>['getClusters']>[number];

export default function MapScreen() {
  const { session } = useAuth();
  const mapRef = useRef<MapView>(null);
  const scRef = useRef<Supercluster<FarmProperties> | null>(null);

  const [clusters, setClusters] = useState<AnyCluster[]>([]);
  const [selectedFarm, setSelectedFarm] = useState<FarmProperties | null>(null);
  const [layers, setLayers] = useState<MapLayer[]>(INITIAL_LAYERS);
  const isEmployee = session?.role ? ['employee', 'expert', 'admin'].includes(session.role) : false;

  const farmsLayer = layers.find((l) => l.id === 'farms');
  const portfolioLayer = layers.find((l) => l.id === 'portfolio');

  useEffect(() => {
    const sc = new Supercluster<FarmProperties>({ radius: 50, maxZoom: 16 });
    sc.load(FARM_FEATURES);
    scRef.current = sc;
    const zoom = getZoomLevel(INITIAL_REGION.longitudeDelta);
    const bbox = getBBox(INITIAL_REGION);
    setClusters(sc.getClusters(bbox, zoom));
  }, []);

  const onRegionChangeComplete = useCallback(
    (region: Region) => {
      if (!scRef.current || !farmsLayer?.active) return;
      const zoom = getZoomLevel(region.longitudeDelta);
      const bbox = getBBox(region);
      setClusters(scRef.current.getClusters(bbox, zoom));
    },
    [farmsLayer?.active],
  );

  const handleClusterPress = useCallback((cluster: AnyCluster) => {
    if (!scRef.current || !cluster.properties.cluster) return;
    const clusterId = cluster.properties.cluster_id as number;
    const expansionZoom = Math.min(scRef.current.getClusterExpansionZoom(clusterId), 16);
    const [lng, lat] = cluster.geometry.coordinates;
    mapRef.current?.animateToRegion(
      {
        latitude: lat,
        longitude: lng,
        latitudeDelta: 360 / Math.pow(2, expansionZoom),
        longitudeDelta: 360 / Math.pow(2, expansionZoom),
      },
      400,
    );
  }, []);

  const toggleLayer = useCallback((id: string) => {
    setLayers((prev) => prev.map((l) => (l.id === id ? { ...l, active: !l.active } : l)));
  }, []);

  return (
    <View style={{ flex: 1 }}>
      {/* Title chip */}
      <View
        style={{
          position: 'absolute',
          top: 16,
          left: 16,
          backgroundColor: '#1a5c36',
          borderRadius: 12,
          paddingHorizontal: 14,
          paddingVertical: 8,
          zIndex: 10,
          shadowColor: '#000',
          shadowOffset: { width: 0, height: 2 },
          shadowOpacity: 0.2,
          shadowRadius: 4,
          elevation: 4,
        }}
      >
        <Text style={{ color: 'white', fontWeight: '700', fontSize: 13 }}>Карта хозяйств ЗКО</Text>
        <Text style={{ color: 'rgba(255,255,255,0.7)', fontSize: 10 }}>
          {FARM_FEATURES.length} хозяйств
        </Text>
      </View>

      <MapView
        ref={mapRef}
        style={{ flex: 1 }}
        initialRegion={INITIAL_REGION}
        onRegionChangeComplete={onRegionChangeComplete}
        onPress={() => setSelectedFarm(null)}
        showsUserLocation
        showsCompass
        showsScale
      >
        {/* Farm markers / clusters */}
        {farmsLayer?.active &&
          clusters.map((cluster) => {
            const [lng, lat] = cluster.geometry.coordinates;
            const isCluster = Boolean(cluster.properties.cluster);

            if (isCluster) {
              return (
                <Marker
                  key={`cluster-${cluster.properties.cluster_id}`}
                  coordinate={{ latitude: lat, longitude: lng }}
                  onPress={() => handleClusterPress(cluster)}
                  tracksViewChanges={false}
                >
                  <ClusterMarker count={cluster.properties.point_count as number} />
                </Marker>
              );
            }

            const props = cluster.properties as FarmProperties;
            return (
              <Marker
                key={`farm-${props.id}`}
                coordinate={{ latitude: lat, longitude: lng }}
                onPress={() => setSelectedFarm(props)}
                tracksViewChanges={false}
              >
                <FarmMarker hasLoan={props.hasActiveLoan} />
              </Marker>
            );
          })}

        {/* Employee+: portfolio circles proportional to loan amount */}
        {isEmployee &&
          portfolioLayer?.active &&
          FARM_FEATURES.filter((f) => f.properties.hasActiveLoan && f.properties.loanAmount).map(
            (f) => {
              const [lng, lat] = f.geometry.coordinates;
              const radius = Math.sqrt((f.properties.loanAmount ?? 0) / 1_000_000) * 500;
              return (
                <Circle
                  key={`portfolio-${f.properties.id}`}
                  center={{ latitude: lat, longitude: lng }}
                  radius={radius}
                  fillColor="rgba(26,92,54,0.15)"
                  strokeColor="rgba(26,92,54,0.4)"
                  strokeWidth={1}
                />
              );
            },
          )}
      </MapView>

      <LayerPanel layers={layers} onToggle={toggleLayer} isEmployee={isEmployee} />

      {selectedFarm && (
        <FarmPopup
          farm={selectedFarm}
          onClose={() => setSelectedFarm(null)}
          isEmployee={isEmployee}
        />
      )}
    </View>
  );
}
