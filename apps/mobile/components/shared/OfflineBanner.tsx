import { Text, View } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useNetworkStatus } from '~/hooks/useNetworkStatus';

export function OfflineBanner() {
  const { isOnline } = useNetworkStatus();

  if (isOnline) return null;

  return (
    <View className="bg-red-500 flex-row items-center justify-center gap-2 px-4 py-2">
      <Ionicons name="cloud-offline-outline" size={14} color="white" />
      <Text className="text-white text-xs font-medium">
        Нет подключения · Отображаются кэшированные данные
      </Text>
    </View>
  );
}
