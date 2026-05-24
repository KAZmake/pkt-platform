import { FlatList, Pressable, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';
import { StatusBadge } from '~/components/ui/Badge';
import type { ApplicationStatus } from '@pkt/shared';

// Mock — replace with: apiClient.get<Application[]>('/cabinet/applications', token)
const MOCK_APPS: {
  id: string;
  number: string;
  program: string;
  amount: number;
  status: ApplicationStatus;
  createdAt: string;
  updatedAt: string;
}[] = [
  {
    id: '1',
    number: '2024-0142',
    program: 'Весенне-полевые работы',
    amount: 5_000_000,
    status: 'issued',
    createdAt: '10 мар 2024',
    updatedAt: '14 мар 2024',
  },
  {
    id: '2',
    number: '2023-0089',
    program: 'Животноводство и откорм',
    amount: 8_000_000,
    status: 'issued',
    createdAt: '5 янв 2023',
    updatedAt: '12 янв 2023',
  },
  {
    id: '3',
    number: '2025-0015',
    program: 'Техника и оборудование',
    amount: 15_000_000,
    status: 'credit_committee',
    createdAt: '20 янв 2025',
    updatedAt: '22 янв 2025',
  },
];

function fmt(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)} млн ₸`;
  return `${(n / 1_000).toFixed(0)} тыс ₸`;
}

export default function ApplicationsScreen() {
  const router = useRouter();

  return (
    <Screen>
      <View className="bg-brand-green-500 px-4 pt-6 pb-6">
        <Pressable
          onPress={() => router.back()}
          className="flex-row items-center gap-1 mb-4 self-start"
        >
          <Ionicons name="chevron-back" size={18} color="rgba(255,255,255,0.8)" />
          <Text className="text-white/80 text-sm">Назад</Text>
        </Pressable>
        <Text className="text-white text-xl font-bold">Мои заявки</Text>
        <Text className="text-brand-green-100 text-sm mt-1">История обращений</Text>
      </View>

      <FlatList
        data={MOCK_APPS}
        keyExtractor={(item) => item.id}
        showsVerticalScrollIndicator={false}
        contentContainerStyle={{ padding: 16, gap: 12, paddingBottom: 32 }}
        ListEmptyComponent={
          <View className="items-center py-16 gap-3">
            <Ionicons name="document-outline" size={40} color="#d1d5db" />
            <Text className="text-gray-400">Заявок пока нет</Text>
          </View>
        }
        renderItem={({ item }) => (
          <View className="bg-white rounded-2xl border border-gray-100 shadow-sm p-4">
            <View className="flex-row items-start justify-between mb-3">
              <View>
                <Text className="text-gray-400 text-xs mb-0.5">Заявка № {item.number}</Text>
                <Text className="text-gray-800 text-sm font-bold" numberOfLines={2}>
                  {item.program}
                </Text>
              </View>
              <StatusBadge status={item.status} />
            </View>
            <View className="flex-row gap-4 border-t border-gray-100 pt-3">
              <View className="flex-1">
                <Text className="text-gray-400 text-xs">Сумма</Text>
                <Text className="text-gray-700 text-sm font-semibold">{fmt(item.amount)}</Text>
              </View>
              <View className="flex-1">
                <Text className="text-gray-400 text-xs">Подана</Text>
                <Text className="text-gray-700 text-sm">{item.createdAt}</Text>
              </View>
              <View className="flex-1">
                <Text className="text-gray-400 text-xs">Обновлено</Text>
                <Text className="text-gray-700 text-sm">{item.updatedAt}</Text>
              </View>
            </View>
          </View>
        )}
      />
    </Screen>
  );
}
