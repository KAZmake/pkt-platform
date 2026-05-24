import { Pressable, ScrollView, Text, View } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import type { FarmProperties } from '~/lib/map/data';

interface FarmPopupProps {
  farm: FarmProperties;
  onClose: () => void;
  isEmployee?: boolean;
}

export function FarmPopup({ farm, onClose, isEmployee }: FarmPopupProps) {
  return (
    <View
      style={{
        position: 'absolute',
        bottom: 88,
        left: 16,
        right: 16,
        backgroundColor: 'white',
        borderRadius: 24,
        padding: 20,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: -2 },
        shadowOpacity: 0.15,
        shadowRadius: 12,
        elevation: 8,
      }}
    >
      {/* Header */}
      <View className="flex-row items-start justify-between mb-3">
        <View className="flex-1 mr-3">
          <Text className="text-gray-800 text-base font-bold leading-5">{farm.name}</Text>
          <Text className="text-gray-400 text-xs mt-0.5">{farm.district} р-н</Text>
        </View>
        <Pressable
          onPress={onClose}
          className="w-7 h-7 bg-gray-100 rounded-full items-center justify-center"
        >
          <Ionicons name="close" size={14} color="#6b7280" />
        </Pressable>
      </View>

      {/* Crops / activities */}
      <ScrollView horizontal showsHorizontalScrollIndicator={false} className="mb-3">
        <View className="flex-row gap-2">
          {farm.crops.map((crop) => (
            <View key={crop} className="bg-brand-green-50 rounded-full px-2.5 py-1">
              <Text className="text-brand-green-600 text-xs font-medium">{crop}</Text>
            </View>
          ))}
        </View>
      </ScrollView>

      {/* Stats row */}
      <View className="flex-row gap-4 mb-4">
        <View>
          <Text className="text-gray-400 text-xs">Площадь</Text>
          <Text className="text-gray-800 text-sm font-semibold">{farm.area} га</Text>
        </View>
        <View>
          <Text className="text-gray-400 text-xs">Владелец</Text>
          <Text className="text-gray-800 text-sm font-semibold" numberOfLines={1}>
            {farm.owner}
          </Text>
        </View>
        <View>
          <Text className="text-gray-400 text-xs">Займ</Text>
          <Text
            className={`text-sm font-semibold ${farm.hasActiveLoan ? 'text-brand-green-500' : 'text-gray-400'}`}
          >
            {farm.hasActiveLoan ? 'Активен' : 'Нет'}
          </Text>
        </View>
      </View>

      {/* Employee-only: loan amount + dossier link */}
      {isEmployee && farm.hasActiveLoan && farm.loanAmount && (
        <View className="bg-brand-green-50 rounded-xl px-4 py-3 mb-4 flex-row items-center justify-between">
          <View>
            <Text className="text-gray-400 text-xs">Сумма займа</Text>
            <Text className="text-brand-green-700 text-sm font-bold">
              {(farm.loanAmount / 1_000_000).toFixed(1)} млн ₸
            </Text>
          </View>
          <Pressable className="bg-brand-green-500 rounded-xl px-3 py-1.5 active:opacity-80">
            <Text className="text-white text-xs font-semibold">Досье →</Text>
          </Pressable>
        </View>
      )}
    </View>
  );
}
