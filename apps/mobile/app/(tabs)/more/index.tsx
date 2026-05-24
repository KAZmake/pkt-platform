import { Pressable, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';
import { useAuth } from '~/lib/auth/context';

type IoniconName = React.ComponentProps<typeof Ionicons>['name'];

interface MenuItem {
  icon: IoniconName;
  label: string;
  sub?: string;
  route: string;
  color?: string;
}

const MENU_SECTIONS: { title: string; items: MenuItem[] }[] = [
  {
    title: 'Компания',
    items: [
      { icon: 'information-circle-outline', label: 'О компании', route: '/(tabs)/more/about' },
      { icon: 'newspaper-outline', label: 'Новости', route: '/(tabs)/more/news' },
      {
        icon: 'help-circle-outline',
        label: 'FAQ',
        sub: 'Частые вопросы',
        route: '/(tabs)/more/faq',
      },
      {
        icon: 'location-outline',
        label: 'Контакты',
        sub: 'Адреса и телефоны',
        route: '/(tabs)/more/contacts',
      },
    ],
  },
  {
    title: 'Сервисы',
    items: [
      {
        icon: 'calculator-outline',
        label: 'Калькулятор займа',
        route: '/(tabs)/programs',
        color: '#c8921a',
      },
      { icon: 'map-outline', label: 'Карта хозяйств', route: '/(tabs)/map' },
      {
        icon: 'document-text-outline',
        label: 'Подать заявку',
        route: '/(tabs)/programs',
        color: '#1a5c36',
      },
    ],
  },
];

export default function MoreScreen() {
  const router = useRouter();
  const { session, logout } = useAuth();

  return (
    <Screen scroll>
      <View className="bg-brand-green-500 px-4 pt-6 pb-5">
        <Text className="text-white text-xl font-bold">Меню</Text>
        <Text className="text-brand-green-100 text-sm mt-0.5">
          ТОО «Первое кредитное товарищество»
        </Text>
      </View>

      <View className="px-4 pt-5 gap-6 pb-8">
        {MENU_SECTIONS.map((section) => (
          <View key={section.title}>
            <Text className="text-gray-400 text-xs font-semibold uppercase tracking-wider mb-2 px-1">
              {section.title}
            </Text>
            <View className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
              {section.items.map((item, i) => (
                <Pressable
                  key={item.route + item.label}
                  className={`flex-row items-center gap-3 px-4 py-3.5 active:bg-gray-50 ${i < section.items.length - 1 ? 'border-b border-gray-100' : ''}`}
                  onPress={() => router.push(item.route as never)}
                >
                  <View className="w-9 h-9 bg-brand-green-50 rounded-xl items-center justify-center">
                    <Ionicons name={item.icon} size={18} color={item.color ?? '#1a5c36'} />
                  </View>
                  <View className="flex-1">
                    <Text className="text-gray-800 text-sm font-semibold">{item.label}</Text>
                    {item.sub && <Text className="text-gray-400 text-xs mt-0.5">{item.sub}</Text>}
                  </View>
                  <Ionicons name="chevron-forward-outline" size={16} color="#d1d5db" />
                </Pressable>
              ))}
            </View>
          </View>
        ))}

        {/* Account section */}
        <View>
          <Text className="text-gray-400 text-xs font-semibold uppercase tracking-wider mb-2 px-1">
            Аккаунт
          </Text>
          <View className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
            {session ? (
              <>
                <View className="flex-row items-center gap-3 px-4 py-3.5 border-b border-gray-100">
                  <View className="w-9 h-9 bg-brand-green-50 rounded-xl items-center justify-center">
                    <Ionicons name="person-outline" size={18} color="#1a5c36" />
                  </View>
                  <View className="flex-1">
                    <Text className="text-gray-800 text-sm font-semibold">
                      {session.firstName} {session.lastName}
                    </Text>
                    <Text className="text-gray-400 text-xs">{session.email}</Text>
                  </View>
                </View>
                <Pressable
                  className="flex-row items-center gap-3 px-4 py-3.5 active:bg-gray-50"
                  onPress={logout}
                >
                  <View className="w-9 h-9 bg-red-50 rounded-xl items-center justify-center">
                    <Ionicons name="log-out-outline" size={18} color="#ef4444" />
                  </View>
                  <Text className="text-red-500 text-sm font-semibold flex-1">Выйти</Text>
                </Pressable>
              </>
            ) : (
              <Pressable
                className="flex-row items-center gap-3 px-4 py-3.5 active:bg-gray-50"
                onPress={() => router.push('/(auth)/login')}
              >
                <View className="w-9 h-9 bg-brand-green-50 rounded-xl items-center justify-center">
                  <Ionicons name="log-in-outline" size={18} color="#1a5c36" />
                </View>
                <Text className="text-brand-green-500 text-sm font-semibold flex-1">
                  Войти в личный кабинет
                </Text>
                <Ionicons name="chevron-forward-outline" size={16} color="#d1d5db" />
              </Pressable>
            )}
          </View>
        </View>

        <Text className="text-gray-300 text-xs text-center">ТОО «ПКТ» · Уральск, ЗКО · v1.0.0</Text>
      </View>
    </Screen>
  );
}
