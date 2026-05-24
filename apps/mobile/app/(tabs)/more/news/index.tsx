import { FlatList, Pressable, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';

// Mock — replace with: directusClient.getNews()
const NEWS = [
  {
    id: '1',
    title: 'Новые условия кредитования на 2025 год',
    date: '15 янв 2025',
    category: 'Программы',
    excerpt:
      'ТОО «ПКТ» объявляет о снижении процентных ставок по ряду программ кредитования для фермеров ЗКО в новом году.',
  },
  {
    id: '2',
    title: 'ПКТ расширяет присутствие в районах ЗКО',
    date: '10 янв 2025',
    category: 'Компания',
    excerpt:
      'Открытие новых представительств в районах Западно-Казахстанской области для улучшения доступности финансовых услуг.',
  },
  {
    id: '3',
    title: 'Запуск мобильного приложения для заёмщиков',
    date: '5 янв 2025',
    category: 'Технологии',
    excerpt:
      'Теперь управлять займом, просматривать документы и отправлять обращения можно прямо со смартфона.',
  },
  {
    id: '4',
    title: 'Итоги 2024 года: рекордные показатели',
    date: '28 дек 2024',
    category: 'Компания',
    excerpt:
      'В 2024 году ТОО «ПКТ» выдало займов на сумму свыше 1,2 млрд тенге, поддержав более 150 фермерских хозяйств.',
  },
  {
    id: '5',
    title: 'Семинар для фермеров в Уральске',
    date: '15 дек 2024',
    category: 'События',
    excerpt:
      'Приглашаем фермеров ЗКО на бесплатный семинар по агрофинансированию и льготным программам.',
  },
];

export default function NewsScreen() {
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
        <Text className="text-white text-xl font-bold">Новости</Text>
        <Text className="text-brand-green-100 text-sm mt-1">Последние события ТОО «ПКТ»</Text>
      </View>

      <FlatList
        data={NEWS}
        keyExtractor={(item) => item.id}
        showsVerticalScrollIndicator={false}
        contentContainerStyle={{ padding: 16, gap: 12, paddingBottom: 32 }}
        renderItem={({ item }) => (
          <Pressable
            className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden active:opacity-80"
            onPress={() => router.push(`/(tabs)/more/news/${item.id}`)}
          >
            {/* Placeholder image */}
            <View className="h-32 bg-brand-green-50 items-center justify-center">
              <Ionicons name="newspaper-outline" size={32} color="#c8e9d4" />
            </View>
            <View className="p-4">
              <View className="flex-row items-center gap-2 mb-2">
                <View className="bg-brand-green-50 rounded-full px-2.5 py-0.5">
                  <Text className="text-brand-green-600 text-xs font-medium">{item.category}</Text>
                </View>
                <Text className="text-gray-400 text-xs">{item.date}</Text>
              </View>
              <Text className="text-gray-800 text-sm font-bold mb-1.5" numberOfLines={2}>
                {item.title}
              </Text>
              <Text className="text-gray-500 text-xs leading-5" numberOfLines={2}>
                {item.excerpt}
              </Text>
            </View>
          </Pressable>
        )}
      />
    </Screen>
  );
}
