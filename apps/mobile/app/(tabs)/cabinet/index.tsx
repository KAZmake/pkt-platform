import { Pressable, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';
import { Card } from '~/components/ui/Card';
import { useAuth } from '~/lib/auth/context';

// Mock — replace with: apiClient.get<CabinetSummary>('/cabinet/summary', token)
const MOCK_SUMMARY = {
  activeLoan: {
    id: 'loan-001',
    program: 'Весенне-полевые работы',
    amount: 5_000_000,
    balance: 3_200_000,
    nextPayment: { date: '2025-02-15', amount: 450_000 },
    progress: 36,
  },
  unreadNotifications: 2,
  pendingApplications: 1,
};

type IoniconName = React.ComponentProps<typeof Ionicons>['name'];

const MENU_ITEMS: {
  icon: IoniconName;
  label: string;
  sub: string;
  route: string;
  badge?: number;
}[] = [
  {
    icon: 'folder-outline',
    label: 'Документы',
    sub: 'Договоры и справки',
    route: '/(tabs)/cabinet/documents',
  },
  {
    icon: 'calendar-outline',
    label: 'График платежей',
    sub: 'Расписание выплат',
    route: '/(tabs)/cabinet/schedule',
  },
  {
    icon: 'document-text-outline',
    label: 'Мои заявки',
    sub: 'История обращений',
    route: '/(tabs)/cabinet/applications',
    badge: MOCK_SUMMARY.pendingApplications,
  },
  {
    icon: 'notifications-outline',
    label: 'Уведомления',
    sub: 'Сообщения от ПКТ',
    route: '/(tabs)/cabinet/notifications',
    badge: MOCK_SUMMARY.unreadNotifications,
  },
];

function fmt(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)} млн ₸`;
  if (n >= 1_000) return `${(n / 1_000).toFixed(0)} тыс ₸`;
  return `${n} ₸`;
}

export default function CabinetDashboardScreen() {
  const router = useRouter();
  const { session, logout } = useAuth();
  const loan = MOCK_SUMMARY.activeLoan;

  return (
    <Screen scroll>
      {/* Header */}
      <View className="bg-brand-green-500 px-4 pt-6 pb-10">
        <View className="flex-row items-center justify-between mb-1">
          <View>
            <Text className="text-brand-green-100 text-xs">Личный кабинет</Text>
            <Text className="text-white text-lg font-bold">
              {session?.firstName} {session?.lastName}
            </Text>
          </View>
          <Pressable
            onPress={logout}
            className="w-9 h-9 bg-white/15 rounded-full items-center justify-center active:bg-white/25"
          >
            <Ionicons name="log-out-outline" size={18} color="white" />
          </Pressable>
        </View>
        <Text className="text-brand-green-100 text-xs">{session?.email}</Text>
      </View>

      <View className="px-4 -mt-5 pb-8 gap-4">
        {/* Active loan card */}
        <Card className="border-0 bg-white shadow-md">
          <View className="flex-row items-start justify-between mb-4">
            <View>
              <Text className="text-gray-400 text-xs mb-0.5">Активный займ</Text>
              <Text className="text-gray-800 text-sm font-bold" numberOfLines={1}>
                {loan.program}
              </Text>
            </View>
            <View className="bg-brand-green-50 rounded-xl px-2.5 py-1">
              <Text className="text-brand-green-600 text-xs font-semibold">Активен</Text>
            </View>
          </View>

          {/* Progress bar */}
          <View className="mb-4">
            <View className="flex-row justify-between mb-1.5">
              <Text className="text-gray-400 text-xs">Погашено</Text>
              <Text className="text-gray-600 text-xs font-medium">{loan.progress}%</Text>
            </View>
            <View className="h-2 bg-gray-100 rounded-full overflow-hidden">
              <View
                className="h-full bg-brand-green-500 rounded-full"
                style={{ width: `${loan.progress}%` }}
              />
            </View>
          </View>

          <View className="flex-row gap-4">
            <View className="flex-1">
              <Text className="text-gray-400 text-xs mb-0.5">Остаток долга</Text>
              <Text className="text-gray-800 font-bold">{fmt(loan.balance)}</Text>
            </View>
            <View className="flex-1">
              <Text className="text-gray-400 text-xs mb-0.5">Следующий платёж</Text>
              <Text className="text-gray-800 font-bold">{fmt(loan.nextPayment.amount)}</Text>
              <Text className="text-gray-400 text-xs">{loan.nextPayment.date}</Text>
            </View>
          </View>
        </Card>

        {/* Menu */}
        <View className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          {MENU_ITEMS.map((item, i) => (
            <Pressable
              key={item.route}
              className={`flex-row items-center gap-3 px-4 py-3.5 active:bg-gray-50 ${i < MENU_ITEMS.length - 1 ? 'border-b border-gray-100' : ''}`}
              onPress={() => router.push(item.route as never)}
            >
              <View className="w-9 h-9 bg-brand-green-50 rounded-xl items-center justify-center">
                <Ionicons name={item.icon} size={18} color="#1a5c36" />
              </View>
              <View className="flex-1">
                <Text className="text-gray-800 text-sm font-semibold">{item.label}</Text>
                <Text className="text-gray-400 text-xs mt-0.5">{item.sub}</Text>
              </View>
              {item.badge ? (
                <View className="bg-brand-gold-500 w-5 h-5 rounded-full items-center justify-center">
                  <Text className="text-white text-xs font-bold">{item.badge}</Text>
                </View>
              ) : (
                <Ionicons name="chevron-forward-outline" size={16} color="#d1d5db" />
              )}
            </Pressable>
          ))}
        </View>
      </View>
    </Screen>
  );
}
