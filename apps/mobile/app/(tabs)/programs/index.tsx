import { FlatList, Pressable, Text, TextInput, View } from 'react-native';
import { useState } from 'react';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';
import type { ActivityType, LoanProgram } from '@pkt/shared';

// Mock — replace with: apiClient.get<LoanProgram[]>('/programs')
const PROGRAMS: LoanProgram[] = [
  {
    id: '1',
    name: 'Весенне-полевые работы',
    nameKz: 'Көктемгі дала жұмыстары',
    rate: 6,
    minAmount: 500_000,
    maxAmount: 50_000_000,
    minTermMonths: 6,
    maxTermMonths: 12,
    activityTypes: ['crop_farming'],
    isActive: true,
  },
  {
    id: '2',
    name: 'Животноводство и откорм КРС',
    rate: 7,
    minAmount: 1_000_000,
    maxAmount: 100_000_000,
    minTermMonths: 12,
    maxTermMonths: 60,
    activityTypes: ['livestock'],
    isActive: true,
  },
  {
    id: '3',
    name: 'Агробизнес — старт',
    rate: 8,
    minAmount: 200_000,
    maxAmount: 10_000_000,
    minTermMonths: 6,
    maxTermMonths: 36,
    activityTypes: ['mixed'],
    isActive: true,
  },
  {
    id: '4',
    name: 'Техника и оборудование',
    rate: 9,
    minAmount: 2_000_000,
    maxAmount: 200_000_000,
    minTermMonths: 12,
    maxTermMonths: 84,
    activityTypes: ['crop_farming', 'livestock'],
    isActive: true,
  },
  {
    id: '5',
    name: 'Переработка сельхозпродукции',
    rate: 8.5,
    minAmount: 5_000_000,
    maxAmount: 500_000_000,
    minTermMonths: 24,
    maxTermMonths: 120,
    activityTypes: ['mixed'],
    isActive: true,
  },
];

const FILTER_LABELS: Record<ActivityType | 'all', string> = {
  all: 'Все',
  crop_farming: 'Растениеводство',
  livestock: 'Животноводство',
  mixed: 'Смешанное',
};

function fmt(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(0)} млн ₸`;
  if (n >= 1_000) return `${(n / 1_000).toFixed(0)} тыс ₸`;
  return `${n} ₸`;
}

type Filter = ActivityType | 'all';

export default function ProgramsScreen() {
  const router = useRouter();
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState<Filter>('all');

  const filtered = PROGRAMS.filter((p) => {
    const matchSearch = p.name.toLowerCase().includes(search.toLowerCase());
    const matchFilter = filter === 'all' || p.activityTypes.includes(filter);
    return matchSearch && matchFilter;
  });

  return (
    <Screen>
      {/* Header */}
      <View className="bg-brand-green-500 px-4 pt-6 pb-4">
        <Text className="text-brand-green-100 text-xs font-medium uppercase tracking-wide">
          Финансирование
        </Text>
        <Text className="text-white text-xl font-bold mt-0.5">Программы кредитования</Text>
      </View>

      {/* Search */}
      <View className="px-4 py-3 bg-white border-b border-gray-100">
        <View className="flex-row items-center bg-gray-100 rounded-xl px-3 py-2 gap-2">
          <Ionicons name="search-outline" size={16} color="#9ca3af" />
          <TextInput
            className="flex-1 text-gray-700 text-sm"
            placeholder="Поиск программы..."
            placeholderTextColor="#9ca3af"
            value={search}
            onChangeText={setSearch}
          />
        </View>
      </View>

      {/* Filters */}
      <View className="bg-white border-b border-gray-100 pb-3">
        <FlatList
          horizontal
          data={['all', 'crop_farming', 'livestock', 'mixed'] as Filter[]}
          keyExtractor={(item) => item}
          showsHorizontalScrollIndicator={false}
          contentContainerStyle={{ paddingHorizontal: 16, gap: 8, paddingTop: 12 }}
          renderItem={({ item }) => (
            <Pressable
              onPress={() => setFilter(item)}
              className={`px-3 py-1.5 rounded-full border ${filter === item ? 'bg-brand-green-500 border-brand-green-500' : 'bg-white border-gray-200'}`}
            >
              <Text
                className={`text-xs font-medium ${filter === item ? 'text-white' : 'text-gray-600'}`}
              >
                {FILTER_LABELS[item]}
              </Text>
            </Pressable>
          )}
        />
      </View>

      {/* List */}
      <FlatList
        data={filtered}
        keyExtractor={(item) => item.id}
        contentContainerStyle={{ padding: 16, gap: 12 }}
        showsVerticalScrollIndicator={false}
        ListEmptyComponent={
          <View className="items-center py-16">
            <Ionicons name="search-outline" size={40} color="#d1d5db" />
            <Text className="text-gray-400 mt-3">Программы не найдены</Text>
          </View>
        }
        renderItem={({ item }) => (
          <Pressable
            className="bg-white rounded-2xl p-4 border border-gray-100 shadow-sm active:opacity-80"
            onPress={() => router.push(`/(tabs)/programs/${item.id}`)}
          >
            <View className="flex-row items-start justify-between mb-2">
              <View className="bg-brand-green-50 rounded-xl px-2.5 py-1">
                <Text className="text-brand-green-600 text-xs font-bold">{item.rate}% годовых</Text>
              </View>
              <Ionicons name="chevron-forward-outline" size={16} color="#d1d5db" />
            </View>
            <Text className="text-gray-800 text-base font-bold mb-2">{item.name}</Text>
            <View className="flex-row gap-4">
              <View>
                <Text className="text-gray-400 text-xs">Сумма</Text>
                <Text className="text-gray-700 text-sm font-medium">
                  {fmt(item.minAmount)} – {fmt(item.maxAmount)}
                </Text>
              </View>
              <View>
                <Text className="text-gray-400 text-xs">Срок</Text>
                <Text className="text-gray-700 text-sm font-medium">
                  {item.minTermMonths} – {item.maxTermMonths} мес.
                </Text>
              </View>
            </View>
          </Pressable>
        )}
      />
    </Screen>
  );
}
