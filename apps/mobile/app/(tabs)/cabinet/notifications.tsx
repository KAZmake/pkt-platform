import { FlatList, Pressable, Text, View } from 'react-native';
import { useState } from 'react';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';

// Mock — replace with: apiClient.get<Notification[]>('/cabinet/notifications', token)
const MOCK_NOTIFICATIONS = [
  {
    id: '1',
    title: 'Платёж принят',
    body: 'Ваш платёж на сумму 375 000 ₸ от 15 января 2025 успешно зачислен.',
    time: '2 ч назад',
    type: 'success' as const,
    read: false,
  },
  {
    id: '2',
    title: 'Напоминание о платеже',
    body: 'Через 5 дней (15 февраля 2025) предстоит платёж на сумму 366 000 ₸.',
    time: '1 день назад',
    type: 'warning' as const,
    read: false,
  },
  {
    id: '3',
    title: 'Статус заявки обновлён',
    body: 'Ваша заявка № 2025-0015 передана в кредитный комитет.',
    time: '3 дня назад',
    type: 'info' as const,
    read: true,
  },
  {
    id: '4',
    title: 'Документ доступен',
    body: 'Справка об остатке долга сформирована и доступна в разделе «Документы».',
    time: '5 дней назад',
    type: 'info' as const,
    read: true,
  },
  {
    id: '5',
    title: 'Срок страховки истекает',
    body: 'Страховой полис на залоговое имущество истекает через 30 дней. Необходимо продление.',
    time: '1 нед назад',
    type: 'warning' as const,
    read: true,
  },
];

const TYPE_CONFIG = {
  success: { icon: 'checkmark-circle' as const, color: '#10b981', bg: 'bg-green-50' },
  warning: { icon: 'warning' as const, color: '#f59e0b', bg: 'bg-amber-50' },
  info: { icon: 'information-circle' as const, color: '#3b82f6', bg: 'bg-blue-50' },
};

export default function NotificationsScreen() {
  const router = useRouter();
  const [notifications, setNotifications] = useState(MOCK_NOTIFICATIONS);

  const markAllRead = () => {
    setNotifications((prev) => prev.map((n) => ({ ...n, read: true })));
  };

  const unreadCount = notifications.filter((n) => !n.read).length;

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
        <View className="flex-row items-center justify-between">
          <View>
            <Text className="text-white text-xl font-bold">Уведомления</Text>
            <Text className="text-brand-green-100 text-sm mt-1">
              {unreadCount > 0 ? `${unreadCount} непрочитанных` : 'Все прочитаны'}
            </Text>
          </View>
          {unreadCount > 0 && (
            <Pressable
              onPress={markAllRead}
              className="bg-white/20 rounded-xl px-3 py-1.5 active:bg-white/30"
            >
              <Text className="text-white text-xs font-medium">Прочитать все</Text>
            </Pressable>
          )}
        </View>
      </View>

      <FlatList
        data={notifications}
        keyExtractor={(item) => item.id}
        showsVerticalScrollIndicator={false}
        contentContainerStyle={{ padding: 16, gap: 8, paddingBottom: 32 }}
        renderItem={({ item }) => {
          const cfg = TYPE_CONFIG[item.type];
          return (
            <Pressable
              className={`bg-white rounded-2xl border shadow-sm p-4 flex-row gap-3 active:opacity-80 ${!item.read ? 'border-brand-green-100' : 'border-gray-100'}`}
              onPress={() =>
                setNotifications((prev) =>
                  prev.map((n) => (n.id === item.id ? { ...n, read: true } : n)),
                )
              }
            >
              <View
                className={`w-9 h-9 ${cfg.bg} rounded-xl items-center justify-center shrink-0 mt-0.5`}
              >
                <Ionicons name={cfg.icon} size={20} color={cfg.color} />
              </View>
              <View className="flex-1">
                <View className="flex-row items-center gap-2 mb-0.5">
                  <Text className="text-gray-800 text-sm font-semibold flex-1">{item.title}</Text>
                  {!item.read && <View className="w-2 h-2 bg-brand-green-500 rounded-full" />}
                </View>
                <Text className="text-gray-500 text-xs leading-5">{item.body}</Text>
                <Text className="text-gray-300 text-xs mt-1.5">{item.time}</Text>
              </View>
            </Pressable>
          );
        }}
      />
    </Screen>
  );
}
