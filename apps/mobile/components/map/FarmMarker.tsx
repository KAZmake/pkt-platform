import { Text, View } from 'react-native';
import { Ionicons } from '@expo/vector-icons';

interface FarmMarkerProps {
  hasLoan: boolean;
}

export function FarmMarker({ hasLoan }: FarmMarkerProps) {
  return (
    <View
      style={{
        width: 36,
        height: 36,
        borderRadius: 18,
        backgroundColor: hasLoan ? '#1a5c36' : '#6b7280',
        borderWidth: 2.5,
        borderColor: 'white',
        alignItems: 'center',
        justifyContent: 'center',
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.25,
        shadowRadius: 3,
        elevation: 4,
      }}
    >
      <Ionicons name="business" size={14} color="white" />
    </View>
  );
}

interface ClusterMarkerProps {
  count: number;
}

export function ClusterMarker({ count }: ClusterMarkerProps) {
  const size = count > 50 ? 56 : count > 20 ? 48 : 42;
  const fontSize = count > 99 ? 11 : 13;

  return (
    <View
      style={{
        width: size,
        height: size,
        borderRadius: size / 2,
        backgroundColor: '#1a5c36',
        borderWidth: 3,
        borderColor: 'white',
        alignItems: 'center',
        justifyContent: 'center',
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.3,
        shadowRadius: 4,
        elevation: 5,
      }}
    >
      <Text style={{ color: 'white', fontWeight: '700', fontSize }}>{count}</Text>
    </View>
  );
}
