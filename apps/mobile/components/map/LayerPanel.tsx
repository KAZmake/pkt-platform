import { Pressable, Text, View } from 'react-native';
import { useState } from 'react';
import { Ionicons } from '@expo/vector-icons';

export interface MapLayer {
  id: string;
  label: string;
  icon: React.ComponentProps<typeof Ionicons>['name'];
  active: boolean;
  employeeOnly?: boolean;
}

interface LayerPanelProps {
  layers: MapLayer[];
  onToggle: (id: string) => void;
  isEmployee?: boolean;
}

export function LayerPanel({ layers, onToggle, isEmployee }: LayerPanelProps) {
  const [isOpen, setIsOpen] = useState(false);
  const visibleLayers = layers.filter((l) => !l.employeeOnly || isEmployee);

  return (
    <View
      style={{
        position: 'absolute',
        top: 16,
        right: 16,
        alignItems: 'flex-end',
      }}
    >
      {/* Toggle button */}
      <Pressable
        onPress={() => setIsOpen((v) => !v)}
        style={{
          backgroundColor: 'white',
          borderRadius: 12,
          padding: 10,
          shadowColor: '#000',
          shadowOffset: { width: 0, height: 2 },
          shadowOpacity: 0.15,
          shadowRadius: 4,
          elevation: 4,
          flexDirection: 'row',
          alignItems: 'center',
          gap: 6,
        }}
      >
        <Ionicons name="layers-outline" size={20} color="#1a5c36" />
        {isOpen && <Text style={{ color: '#1a5c36', fontSize: 12, fontWeight: '600' }}>Слои</Text>}
      </Pressable>

      {/* Layer list */}
      {isOpen && (
        <View
          style={{
            backgroundColor: 'white',
            borderRadius: 16,
            padding: 12,
            marginTop: 8,
            minWidth: 180,
            shadowColor: '#000',
            shadowOffset: { width: 0, height: 2 },
            shadowOpacity: 0.15,
            shadowRadius: 6,
            elevation: 4,
            gap: 4,
          }}
        >
          {visibleLayers.map((layer) => (
            <Pressable
              key={layer.id}
              onPress={() => onToggle(layer.id)}
              style={{
                flexDirection: 'row',
                alignItems: 'center',
                gap: 10,
                paddingVertical: 8,
                paddingHorizontal: 4,
                borderRadius: 10,
                backgroundColor: layer.active ? '#edf7f1' : 'transparent',
              }}
            >
              <View
                style={{
                  width: 32,
                  height: 32,
                  borderRadius: 8,
                  backgroundColor: layer.active ? '#1a5c36' : '#f3f4f6',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              >
                <Ionicons name={layer.icon} size={16} color={layer.active ? 'white' : '#9ca3af'} />
              </View>
              <Text
                style={{
                  fontSize: 13,
                  fontWeight: '500',
                  color: layer.active ? '#1a5c36' : '#6b7280',
                  flex: 1,
                }}
              >
                {layer.label}
              </Text>
              {layer.employeeOnly && (
                <View
                  style={{
                    backgroundColor: '#c8921a',
                    borderRadius: 4,
                    paddingHorizontal: 4,
                    paddingVertical: 1,
                  }}
                >
                  <Text style={{ color: 'white', fontSize: 9, fontWeight: '700' }}>ЭМП</Text>
                </View>
              )}
            </Pressable>
          ))}
        </View>
      )}
    </View>
  );
}
