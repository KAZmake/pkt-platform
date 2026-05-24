import { FlatList, Pressable, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';
import { Card } from '~/components/ui/Card';
import type { LoanProgram } from '@pkt/shared';

const KPI_ITEMS = [
  { value: '14+', label: 'лет на рынке' },
  { value: '500+', label: 'заёмщиков' },
  { value: '5 млрд ₸', label: 'выдано займов' },
  { value: '12', label: 'программ' },
];

// Mock — replace with: apiClient.get<LoanProgram[]>('/programs?featured=true')
const FEATURED_PROGRAMS: LoanProgram[] = [
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
    name: 'Животноводство и откорм',
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
];

// Mock — replace with: directusClient.getNews(3)
const LATEST_NEWS = [
  {
    id: '1',
    title: 'Новые условия кредитования на 2025 год',
    date: '15 янв 2025',
    category: 'Программы',
  },
  {
    id: '2',
    title: 'ПКТ расширяет присутствие в районах ЗКО',
    date: '10 янв 2025',
    category: 'Новости',
  },
  {
    id: '3',
    title: 'Запуск мобильного приложения для заёмщиков',
    date: '5 янв 2025',
    category: 'Технологии',
  },
];

function fmt(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(0)} млн ₸`;
  if (n >= 1_000) return `${(n / 1_000).toFixed(0)} тыс ₸`;
  return `${n} ₸`;
}

export default function HomeScreen() {
  const router = useRouter();

  return (
    <Screen scroll>
      {/* Brand header */}
      <View className="bg-brand-green-500 px-4 pt-6 pb-10">
        <View className="flex-row items-center justify-between mb-5">
          <View>
            <Text className="text-brand-green-100 text-xs font-medium tracking-wide uppercase">
              ТОО «ПКТ»
            </Text>
            <Text className="text-white text-xl font-bold mt-0.5">Первое кредитное</Text>
          </View>
          <Pressable
            onPress={() => router.push('/(tabs)/cabinet')}
            className="w-10 h-10 bg-white/20 rounded-full items-center justify-center active:bg-white/30"
          >
            <Ionicons name="person-outline" size={20} color="white" />
          </Pressable>
        </View>
        <Pressable className="bg-white/15 rounded-2xl px-4 py-3 flex-row items-center gap-3 active:bg-white/25">
          <Ionicons name="call-outline" size={18} color="white" />
          <Text className="text-white font-semibold text-sm">8 800 080 1700</Text>
          <Text className="text-brand-green-100 text-xs ml-auto">Бесплатно</Text>
        </Pressable>
      </View>

      {/* KPI card */}
      <View className="px-4 -mt-5">
        <Card className="flex-row flex-wrap p-0 overflow-hidden">
          {KPI_ITEMS.map((item, i) => (
            <View
              key={i}
              className={`w-1/2 items-center py-4 px-2 ${i < 2 ? 'border-b border-gray-100' : ''} ${i % 2 === 0 ? 'border-r border-gray-100' : ''}`}
            >
              <Text className="text-brand-green-500 text-lg font-bold">{item.value}</Text>
              <Text className="text-gray-400 text-xs text-center mt-0.5">{item.label}</Text>
            </View>
          ))}
        </Card>
      </View>

      {/* Quick actions */}
      <View className="px-4 mt-5 flex-row gap-3">
        <Pressable
          className="flex-1 bg-brand-gold-500 rounded-2xl py-4 items-center gap-1.5 active:opacity-80"
          onPress={() => router.push('/(tabs)/programs')}
        >
          <Ionicons name="calculator-outline" size={22} color="white" />
          <Text className="text-white text-xs font-semibold">Калькулятор</Text>
        </Pressable>
        <Pressable
          className="flex-1 bg-brand-green-500 rounded-2xl py-4 items-center gap-1.5 active:opacity-80"
          onPress={() => router.push('/(tabs)/programs')}
        >
          <Ionicons name="document-text-outline" size={22} color="white" />
          <Text className="text-white text-xs font-semibold">Подать заявку</Text>
        </Pressable>
        <Pressable
          className="flex-1 bg-gray-100 rounded-2xl py-4 items-center gap-1.5 active:opacity-80"
          onPress={() => router.push('/(tabs)/cabinet')}
        >
          <Ionicons name="person-circle-outline" size={22} color="#1a5c36" />
          <Text className="text-brand-green-500 text-xs font-semibold">Кабинет</Text>
        </Pressable>
      </View>

      {/* Programs */}
      <View className="mt-7">
        <View className="px-4 flex-row items-center justify-between mb-3">
          <Text className="text-gray-800 text-base font-bold">Программы займов</Text>
          <Pressable onPress={() => router.push('/(tabs)/programs')}>
            <Text className="text-brand-green-500 text-sm font-medium">Все →</Text>
          </Pressable>
        </View>
        <FlatList
          horizontal
          data={FEATURED_PROGRAMS}
          keyExtractor={(item) => item.id}
          showsHorizontalScrollIndicator={false}
          contentContainerStyle={{ paddingHorizontal: 16, gap: 12 }}
          renderItem={({ item }) => (
            <Pressable
              className="w-52 bg-white rounded-2xl p-4 border border-gray-100 shadow-sm active:opacity-80"
              onPress={() => router.push(`/(tabs)/programs/${item.id}`)}
            >
              <View className="bg-brand-green-50 rounded-xl px-2.5 py-1 self-start mb-3">
                <Text className="text-brand-green-600 text-xs font-bold">{item.rate}% год.</Text>
              </View>
              <Text className="text-gray-800 text-sm font-bold mb-1.5" numberOfLines={2}>
                {item.name}
              </Text>
              <Text className="text-gray-400 text-xs">
                до {fmt(item.maxAmount)} · {item.maxTermMonths} мес.
              </Text>
            </Pressable>
          )}
        />
      </View>

      {/* News */}
      <View className="px-4 mt-7 mb-8">
        <View className="flex-row items-center justify-between mb-3">
          <Text className="text-gray-800 text-base font-bold">Новости</Text>
          <Pressable onPress={() => router.push('/(tabs)/more/news')}>
            <Text className="text-brand-green-500 text-sm font-medium">Все →</Text>
          </Pressable>
        </View>
        <View className="gap-3">
          {LATEST_NEWS.map((item) => (
            <Pressable
              key={item.id}
              className="bg-white border border-gray-100 rounded-2xl p-4 shadow-sm active:opacity-80 flex-row items-center gap-3"
              onPress={() => router.push(`/(tabs)/more/news/${item.id}`)}
            >
              <View className="flex-1">
                <View className="bg-brand-green-50 rounded-full px-2.5 py-0.5 self-start mb-2">
                  <Text className="text-brand-green-600 text-xs font-medium">{item.category}</Text>
                </View>
                <Text className="text-gray-800 text-sm font-semibold" numberOfLines={2}>
                  {item.title}
                </Text>
                <Text className="text-gray-400 text-xs mt-1">{item.date}</Text>
              </View>
              <Ionicons name="chevron-forward-outline" size={16} color="#d1d5db" />
            </Pressable>
          ))}
        </View>
      </View>
    </Screen>
  );
}
