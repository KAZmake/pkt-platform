import { Pressable, ScrollView, Text, View } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';

// Mock — replace with: directusClient.getNewsById(id)
const NEWS_MAP: Record<
  string,
  { id: string; title: string; date: string; category: string; body: string }
> = {
  '1': {
    id: '1',
    title: 'Новые условия кредитования на 2025 год',
    date: '15 января 2025',
    category: 'Программы',
    body: `ТОО «Первое кредитное товарищество» рады сообщить об улучшении условий кредитования для фермеров Западно-Казахстанской области в 2025 году.

Ключевые изменения:

• Ставка по программе «Весенне-полевые работы» снижена с 7% до 6% годовых
• Максимальный лимит займа по программе «Животноводство» увеличен до 100 млн тенге
• Введена новая программа для малых фермерских хозяйств с упрощённым пакетом документов
• Расширен список принимаемых залогов

Новые условия вступают в силу с 1 февраля 2025 года. Подробности уточняйте у менеджеров в офисах ТОО «ПКТ» или по телефону Call-центра 8 800 080 1700.`,
  },
  '2': {
    id: '2',
    title: 'ПКТ расширяет присутствие в районах ЗКО',
    date: '10 января 2025',
    category: 'Компания',
    body: `В рамках стратегии развития и повышения доступности финансовых услуг для фермеров ТОО «ПКТ» открывает новые представительства в районах Западно-Казахстанской области.

В первом квартале 2025 года откроются офисы в:
• Зеленовском районе (с. Переметное)
• Каратобинском районе (с. Каратобе)
• Жанибекском районе (с. Жанибек)

Часы работы: понедельник–пятница, 9:00–17:00.

Это позволит фермерам из отдалённых районов получать консультации и оформлять займы, не выезжая в Уральск.`,
  },
};

export default function NewsDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();

  const article = NEWS_MAP[id] ?? {
    id,
    title: 'Новость',
    date: '',
    category: 'Новости',
    body: 'Полный текст новости загружается...',
  };

  return (
    <Screen>
      {/* Header image placeholder */}
      <View className="h-48 bg-brand-green-500 px-4 pt-6 pb-4 justify-between">
        <Pressable
          onPress={() => router.back()}
          className="flex-row items-center gap-1 self-start active:opacity-70"
        >
          <Ionicons name="chevron-back" size={18} color="rgba(255,255,255,0.8)" />
          <Text className="text-white/80 text-sm">Назад</Text>
        </Pressable>
        <View>
          <View className="bg-white/20 rounded-full px-2.5 py-0.5 self-start mb-2">
            <Text className="text-white text-xs font-medium">{article.category}</Text>
          </View>
          <Text className="text-brand-green-100 text-xs">{article.date}</Text>
        </View>
      </View>

      <ScrollView
        showsVerticalScrollIndicator={false}
        contentContainerStyle={{ padding: 20, paddingBottom: 40 }}
      >
        <Text className="text-gray-800 text-xl font-bold mb-4 leading-7">{article.title}</Text>
        <Text className="text-gray-600 text-sm leading-7">{article.body}</Text>

        {/* Share / back */}
        <Pressable
          className="mt-8 flex-row items-center gap-2 self-start active:opacity-70"
          onPress={() => router.back()}
        >
          <Ionicons name="arrow-back-outline" size={16} color="#1a5c36" />
          <Text className="text-brand-green-500 text-sm font-medium">Все новости</Text>
        </Pressable>
      </ScrollView>
    </Screen>
  );
}
