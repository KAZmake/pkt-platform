import { Pressable, ScrollView, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';
import { Card } from '~/components/ui/Card';

const STATS = [
  { value: '2010', label: 'год основания' },
  { value: '500+', label: 'заёмщиков' },
  { value: '5 млрд ₸', label: 'выдано займов' },
  { value: '12', label: 'программ' },
];

const VALUES = [
  {
    icon: 'shield-checkmark-outline' as const,
    title: 'Надёжность',
    text: 'Работаем с 2010 года и выполняем все взятые на себя обязательства перед заёмщиками и партнёрами.',
  },
  {
    icon: 'people-outline' as const,
    title: 'Поддержка АПК',
    text: 'Специализируемся на кредитовании аграрного сектора ЗКО — понимаем сезонность и специфику вашего бизнеса.',
  },
  {
    icon: 'trending-up-outline' as const,
    title: 'Развитие',
    text: 'Помогаем фермерам ЗКО развиваться, расширять производство и повышать устойчивость хозяйств.',
  },
];

export default function AboutScreen() {
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
        <Text className="text-white text-xl font-bold">О компании</Text>
        <Text className="text-brand-green-100 text-sm mt-1">
          ТОО «Первое кредитное товарищество»
        </Text>
      </View>

      <ScrollView
        showsVerticalScrollIndicator={false}
        contentContainerStyle={{ padding: 16, gap: 16, paddingBottom: 32 }}
      >
        {/* About text */}
        <Card>
          <Text className="text-gray-800 font-bold mb-3">Кто мы</Text>
          <Text className="text-gray-600 text-sm leading-6">
            ТОО «Первое кредитное товарищество» — ведущая финансовая организация
            Западно-Казахстанской области, специализирующаяся на льготном кредитовании
            агропромышленного комплекса.
          </Text>
          <Text className="text-gray-600 text-sm leading-6 mt-3">
            С 2010 года мы поддерживаем фермеров ЗКО, предоставляя доступные займы для проведения
            полевых работ, развития животноводства, приобретения техники и расширения производства.
          </Text>
        </Card>

        {/* Stats */}
        <View className="flex-row flex-wrap gap-3">
          {STATS.map((s) => (
            <View
              key={s.label}
              className="bg-brand-green-500 rounded-2xl p-4 flex-1 min-w-[40%] items-center"
            >
              <Text className="text-white text-xl font-black">{s.value}</Text>
              <Text className="text-brand-green-100 text-xs text-center mt-0.5">{s.label}</Text>
            </View>
          ))}
        </View>

        {/* Values */}
        <Card>
          <Text className="text-gray-800 font-bold mb-4">Наши принципы</Text>
          <View className="gap-4">
            {VALUES.map((v) => (
              <View key={v.title} className="flex-row gap-3">
                <View className="w-10 h-10 bg-brand-green-50 rounded-xl items-center justify-center shrink-0">
                  <Ionicons name={v.icon} size={20} color="#1a5c36" />
                </View>
                <View className="flex-1">
                  <Text className="text-gray-800 text-sm font-bold mb-1">{v.title}</Text>
                  <Text className="text-gray-500 text-sm leading-5">{v.text}</Text>
                </View>
              </View>
            ))}
          </View>
        </Card>

        {/* Requisites */}
        <Card>
          <Text className="text-gray-800 font-bold mb-3">Реквизиты</Text>
          <View className="gap-2">
            {[
              ['Полное название', 'ТОО «Первое кредитное товарищество»'],
              ['БИН', '100 240 042 714'],
              ['Адрес', 'г. Уральск, ЗКО, Казахстан'],
              ['Лицензия', 'МФО №001 / 2010'],
            ].map(([label, value]) => (
              <View key={label} className="flex-row justify-between gap-2">
                <Text className="text-gray-400 text-sm">{label}</Text>
                <Text
                  className="text-gray-700 text-sm font-medium text-right flex-1"
                  numberOfLines={2}
                >
                  {value}
                </Text>
              </View>
            ))}
          </View>
        </Card>
      </ScrollView>
    </Screen>
  );
}
