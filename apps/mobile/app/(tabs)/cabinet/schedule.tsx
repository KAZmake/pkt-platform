import { FlatList, Pressable, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';

// Mock — replace with: apiClient.get<PaymentSchedule>('/cabinet/schedule', token)
const MOCK_SCHEDULE = {
  loanAmount: 5_000_000,
  remainingBalance: 3_200_000,
  totalPayments: 12,
  paidPayments: 4,
  payments: [
    {
      id: '1',
      date: '2024-04-15',
      principal: 350_000,
      interest: 25_000,
      total: 375_000,
      status: 'paid' as const,
    },
    {
      id: '2',
      date: '2024-05-15',
      principal: 350_000,
      interest: 23_250,
      total: 373_250,
      status: 'paid' as const,
    },
    {
      id: '3',
      date: '2024-06-15',
      principal: 350_000,
      interest: 21_500,
      total: 371_500,
      status: 'paid' as const,
    },
    {
      id: '4',
      date: '2024-07-15',
      principal: 350_000,
      interest: 19_750,
      total: 369_750,
      status: 'paid' as const,
    },
    {
      id: '5',
      date: '2025-02-15',
      principal: 350_000,
      interest: 16_000,
      total: 366_000,
      status: 'upcoming' as const,
    },
    {
      id: '6',
      date: '2025-03-15',
      principal: 350_000,
      interest: 14_250,
      total: 364_250,
      status: 'pending' as const,
    },
    {
      id: '7',
      date: '2025-04-15',
      principal: 350_000,
      interest: 12_500,
      total: 362_500,
      status: 'pending' as const,
    },
    {
      id: '8',
      date: '2025-05-15',
      principal: 350_000,
      interest: 10_750,
      total: 360_750,
      status: 'pending' as const,
    },
  ],
};

const STATUS_CONFIG = {
  paid: {
    label: 'Оплачен',
    bg: 'bg-green-100',
    text: 'text-green-700',
    icon: 'checkmark-circle' as const,
  },
  upcoming: {
    label: 'Ближайший',
    bg: 'bg-amber-100',
    text: 'text-amber-700',
    icon: 'time' as const,
  },
  pending: {
    label: 'Предстоит',
    bg: 'bg-gray-100',
    text: 'text-gray-500',
    icon: 'ellipse-outline' as const,
  },
};

function fmt(n: number): string {
  return n.toLocaleString('ru-KZ') + ' ₸';
}

export default function ScheduleScreen() {
  const router = useRouter();
  const s = MOCK_SCHEDULE;

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
        <Text className="text-white text-xl font-bold">График платежей</Text>
        <Text className="text-brand-green-100 text-sm mt-1">Расписание выплат по займу</Text>
      </View>

      <FlatList
        data={s.payments}
        keyExtractor={(item) => item.id}
        showsVerticalScrollIndicator={false}
        ListHeaderComponent={
          <View className="gap-3 mb-4">
            {/* Summary */}
            <View className="bg-white rounded-2xl border border-gray-100 shadow-sm p-4 flex-row gap-4">
              <View className="flex-1 items-center">
                <Text className="text-gray-400 text-xs mb-0.5">Всего платежей</Text>
                <Text className="text-gray-800 font-bold text-base">{s.totalPayments}</Text>
              </View>
              <View className="w-px bg-gray-100" />
              <View className="flex-1 items-center">
                <Text className="text-gray-400 text-xs mb-0.5">Оплачено</Text>
                <Text className="text-brand-green-500 font-bold text-base">{s.paidPayments}</Text>
              </View>
              <View className="w-px bg-gray-100" />
              <View className="flex-1 items-center">
                <Text className="text-gray-400 text-xs mb-0.5">Остаток</Text>
                <Text className="text-gray-800 font-bold text-base">
                  {s.totalPayments - s.paidPayments}
                </Text>
              </View>
            </View>

            {/* Progress bar */}
            <View className="bg-white rounded-2xl border border-gray-100 shadow-sm p-4">
              <View className="flex-row justify-between mb-2">
                <Text className="text-gray-500 text-sm">Погашение займа</Text>
                <Text className="text-brand-green-500 font-semibold text-sm">
                  {Math.round((s.paidPayments / s.totalPayments) * 100)}%
                </Text>
              </View>
              <View className="h-2.5 bg-gray-100 rounded-full overflow-hidden">
                <View
                  className="h-full bg-brand-green-500 rounded-full"
                  style={{ width: `${(s.paidPayments / s.totalPayments) * 100}%` }}
                />
              </View>
              <Text className="text-gray-400 text-xs mt-2">
                Остаток долга: {fmt(s.remainingBalance)}
              </Text>
            </View>
          </View>
        }
        contentContainerStyle={{ padding: 16, paddingBottom: 32 }}
        renderItem={({ item }) => {
          const cfg = STATUS_CONFIG[item.status];
          return (
            <View
              className={`bg-white rounded-2xl border shadow-sm p-4 mb-3 ${item.status === 'upcoming' ? 'border-amber-200' : 'border-gray-100'}`}
            >
              <View className="flex-row items-center justify-between mb-3">
                <View className="flex-row items-center gap-2">
                  <Ionicons
                    name={cfg.icon}
                    size={16}
                    color={
                      item.status === 'paid'
                        ? '#10b981'
                        : item.status === 'upcoming'
                          ? '#f59e0b'
                          : '#9ca3af'
                    }
                  />
                  <Text className="text-gray-600 text-sm font-medium">{item.date}</Text>
                </View>
                <View className={`px-2.5 py-1 rounded-full ${cfg.bg}`}>
                  <Text className={`text-xs font-semibold ${cfg.text}`}>{cfg.label}</Text>
                </View>
              </View>
              <View className="flex-row gap-4">
                <View className="flex-1">
                  <Text className="text-gray-400 text-xs mb-0.5">Основной долг</Text>
                  <Text className="text-gray-700 text-sm font-medium">{fmt(item.principal)}</Text>
                </View>
                <View className="flex-1">
                  <Text className="text-gray-400 text-xs mb-0.5">Вознаграждение</Text>
                  <Text className="text-gray-700 text-sm font-medium">{fmt(item.interest)}</Text>
                </View>
                <View className="flex-1">
                  <Text className="text-gray-400 text-xs mb-0.5">Итого</Text>
                  <Text className="text-gray-800 text-sm font-bold">{fmt(item.total)}</Text>
                </View>
              </View>
            </View>
          );
        }}
      />
    </Screen>
  );
}
