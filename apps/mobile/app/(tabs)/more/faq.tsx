import { Pressable, ScrollView, Text, View } from 'react-native';
import { useState } from 'react';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';

// Mock — replace with: directusClient.getFaq()
const FAQ_ITEMS = [
  {
    id: '1',
    category: 'Займы',
    question: 'Какие документы нужны для получения займа?',
    answer:
      'Для оформления займа потребуются: удостоверение личности, ИИН, документы на залоговое имущество, документы подтверждающие доход или деятельность хозяйства (договоры, акты сверки, накладные).',
  },
  {
    id: '2',
    category: 'Займы',
    question: 'Какой минимальный срок рассмотрения заявки?',
    answer:
      'Стандартный срок рассмотрения — от 5 до 14 рабочих дней с момента подачи полного пакета документов. При наличии всех документов срок может быть сокращён.',
  },
  {
    id: '3',
    category: 'Займы',
    question: 'Можно ли досрочно погасить займ?',
    answer:
      'Да, досрочное погашение займа возможно без штрафных санкций. Свяжитесь с менеджером или оставьте заявку в личном кабинете для уточнения остатка и условий погашения.',
  },
  {
    id: '4',
    category: 'Залог',
    question: 'Что может быть залогом при оформлении займа?',
    answer:
      'В качестве залога принимаются: земельные участки, сельскохозяйственная техника, животноводческие объекты, товары в обороте, недвижимость. Оценку проводит аккредитованный оценщик.',
  },
  {
    id: '5',
    category: 'Залог',
    question: 'Обязательно ли страховать залог?',
    answer:
      'Страхование залогового имущества является обязательным условием для большинства программ кредитования. ТОО «ПКТ» сотрудничает с аккредитованными страховыми компаниями.',
  },
  {
    id: '6',
    category: 'Личный кабинет',
    question: 'Как подключить личный кабинет?',
    answer:
      'Для доступа к личному кабинету обратитесь в офис ТОО «ПКТ» или позвоните на Call-центр 8 800 080 1700. Менеджер создаст учётную запись и выдаст данные для входа.',
  },
  {
    id: '7',
    category: 'Личный кабинет',
    question: 'Что доступно в личном кабинете заёмщика?',
    answer:
      'В личном кабинете вы можете: просматривать документы по займу, отслеживать график платежей, отправлять обращения и запросы, получать уведомления о важных событиях.',
  },
];

const categories = [...new Set(FAQ_ITEMS.map((f) => f.category))];

export default function FaqScreen() {
  const router = useRouter();
  const [openId, setOpenId] = useState<string | null>(null);
  const [activeCategory, setActiveCategory] = useState('Все');

  const displayed =
    activeCategory === 'Все' ? FAQ_ITEMS : FAQ_ITEMS.filter((f) => f.category === activeCategory);

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
        <Text className="text-white text-xl font-bold">FAQ</Text>
        <Text className="text-brand-green-100 text-sm mt-1">Часто задаваемые вопросы</Text>
      </View>

      {/* Category tabs */}
      <View className="bg-white border-b border-gray-100">
        <ScrollView
          horizontal
          showsHorizontalScrollIndicator={false}
          contentContainerStyle={{ paddingHorizontal: 16, paddingVertical: 12, gap: 8 }}
        >
          {['Все', ...categories].map((cat) => (
            <Pressable
              key={cat}
              onPress={() => setActiveCategory(cat)}
              className={`px-3 py-1.5 rounded-full border ${activeCategory === cat ? 'bg-brand-green-500 border-brand-green-500' : 'bg-white border-gray-200'}`}
            >
              <Text
                className={`text-xs font-medium ${activeCategory === cat ? 'text-white' : 'text-gray-600'}`}
              >
                {cat}
              </Text>
            </Pressable>
          ))}
        </ScrollView>
      </View>

      <ScrollView
        showsVerticalScrollIndicator={false}
        contentContainerStyle={{ padding: 16, gap: 8, paddingBottom: 32 }}
      >
        {displayed.map((item) => {
          const isOpen = openId === item.id;
          return (
            <Pressable
              key={item.id}
              className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden active:opacity-95"
              onPress={() => setOpenId(isOpen ? null : item.id)}
            >
              <View className="flex-row items-center gap-3 px-4 py-3.5">
                <Text className="text-gray-800 text-sm font-semibold flex-1 leading-5">
                  {item.question}
                </Text>
                <Ionicons
                  name={isOpen ? 'chevron-up-outline' : 'chevron-down-outline'}
                  size={16}
                  color="#9ca3af"
                />
              </View>
              {isOpen && (
                <View className="px-4 pb-4 border-t border-gray-100">
                  <Text className="text-gray-500 text-sm leading-6 pt-3">{item.answer}</Text>
                </View>
              )}
            </Pressable>
          );
        })}
      </ScrollView>
    </Screen>
  );
}
