import { Text, View } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';

// Task 4.5: react-native-maps + clusters + farm popup
export default function MapScreen() {
  return (
    <Screen>
      <View className="bg-brand-green-500 px-4 pt-6 pb-4">
        <Text className="text-brand-green-100 text-xs font-medium uppercase tracking-wide">
          Интерактивная карта
        </Text>
        <Text className="text-white text-xl font-bold mt-0.5">Карта хозяйств ЗКО</Text>
      </View>
      <View className="flex-1 items-center justify-center px-8 gap-4">
        <View className="w-20 h-20 bg-brand-green-50 rounded-full items-center justify-center">
          <Ionicons name="map-outline" size={40} color="#1a5c36" />
        </View>
        <Text className="text-gray-700 text-lg font-bold text-center">Карта хозяйств</Text>
        <Text className="text-gray-400 text-sm text-center">
          Интерактивная карта с хозяйствами ЗКО, кластерами и детальным попапом будет доступна в
          следующем обновлении (задача 4.5).
        </Text>
      </View>
    </Screen>
  );
}
