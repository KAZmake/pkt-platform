import { FlatList, Pressable, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';

// Mock — replace with: apiClient.get<Document[]>('/cabinet/documents', token)
const MOCK_DOCS = [
  {
    id: '1',
    name: 'Договор займа № 2024-0142',
    type: 'PDF',
    date: '14 мар 2024',
    size: '1.2 MB',
    category: 'Договор',
  },
  {
    id: '2',
    name: 'График платежей',
    type: 'PDF',
    date: '14 мар 2024',
    size: '0.4 MB',
    category: 'График',
  },
  {
    id: '3',
    name: 'Акт приёма-передачи залога',
    type: 'PDF',
    date: '15 мар 2024',
    size: '0.8 MB',
    category: 'Залог',
  },
  {
    id: '4',
    name: 'Справка об остатке долга',
    type: 'PDF',
    date: '1 янв 2025',
    size: '0.2 MB',
    category: 'Справка',
  },
  {
    id: '5',
    name: 'Страховой полис',
    type: 'PDF',
    date: '16 мар 2024',
    size: '1.5 MB',
    category: 'Залог',
  },
];

const CATEGORY_COLORS: Record<string, string> = {
  Договор: 'bg-brand-green-50 text-brand-green-600',
  График: 'bg-blue-50 text-blue-600',
  Залог: 'bg-amber-50 text-amber-600',
  Справка: 'bg-purple-50 text-purple-600',
};

export default function DocumentsScreen() {
  const router = useRouter();

  const handleOpen = (_doc: (typeof MOCK_DOCS)[0]) => {
    // Mock — replace with: presigned MinIO URL from apiClient
  };

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
        <Text className="text-white text-xl font-bold">Документы</Text>
        <Text className="text-brand-green-100 text-sm mt-1">Договоры, справки и приложения</Text>
      </View>

      <FlatList
        data={MOCK_DOCS}
        keyExtractor={(item) => item.id}
        showsVerticalScrollIndicator={false}
        contentContainerStyle={{ padding: 16, gap: 10, paddingBottom: 32 }}
        renderItem={({ item }) => {
          const colorClass = CATEGORY_COLORS[item.category] ?? 'bg-gray-100 text-gray-600';
          const [bg, text] = colorClass.split(' ');
          return (
            <Pressable
              className="bg-white rounded-2xl border border-gray-100 shadow-sm p-4 flex-row items-center gap-3 active:opacity-80"
              onPress={() => handleOpen(item)}
            >
              <View className="w-11 h-11 bg-red-50 rounded-xl items-center justify-center">
                <Ionicons name="document-outline" size={22} color="#ef4444" />
              </View>
              <View className="flex-1">
                <Text className="text-gray-800 text-sm font-semibold" numberOfLines={1}>
                  {item.name}
                </Text>
                <View className="flex-row items-center gap-2 mt-1">
                  <View className={`px-2 py-0.5 rounded-full ${bg}`}>
                    <Text className={`text-xs font-medium ${text}`}>{item.category}</Text>
                  </View>
                  <Text className="text-gray-400 text-xs">{item.date}</Text>
                  <Text className="text-gray-300 text-xs">{item.size}</Text>
                </View>
              </View>
              <Ionicons name="download-outline" size={18} color="#9ca3af" />
            </Pressable>
          );
        }}
      />
    </Screen>
  );
}
