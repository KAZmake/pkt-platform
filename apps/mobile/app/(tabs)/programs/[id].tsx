import { Pressable, ScrollView, Text, View } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';
import { Card } from '~/components/ui/Card';
import { Button } from '~/components/ui/Button';
import type { LoanProgram } from '@pkt/shared';

// Mock — replace with: apiClient.get<LoanProgram>(`/programs/${id}`)
const PROGRAMS_MAP: Record<
  string,
  LoanProgram & { description?: string; requirements?: string[] }
> = {
  '1': {
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
    description:
      'Льготный займ для проведения весенне-полевых работ: покупка семян, удобрений, ГСМ, аренда техники. Целевой займ с погашением после сбора урожая.',
    requirements: [
      'Регистрация в ЗКО',
      'Площадь пашни от 50 га',
      'Опыт растениеводства от 1 года',
      'Залоговое обеспечение (земля, техника)',
    ],
  },
  '2': {
    id: '2',
    name: 'Животноводство и откорм КРС',
    rate: 7,
    minAmount: 1_000_000,
    maxAmount: 100_000_000,
    minTermMonths: 12,
    maxTermMonths: 60,
    activityTypes: ['livestock'],
    isActive: true,
    description:
      'Финансирование для развития животноводства: закуп скота, строительство животноводческих объектов, приобретение кормов и ветеринарных препаратов.',
    requirements: [
      'Наличие животноводческой базы',
      'Поголовье КРС от 50 голов',
      'Ветеринарное заключение',
      'Залог: скот, недвижимость или техника',
    ],
  },
};

function fmt(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(0)} млн ₸`;
  if (n >= 1_000) return `${(n / 1_000).toFixed(0)} тыс ₸`;
  return `${n} ₸`;
}

const ACTIVITY_LABELS: Record<string, string> = {
  crop_farming: 'Растениеводство',
  livestock: 'Животноводство',
  mixed: 'Смешанное',
};

export default function ProgramDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();

  // Mock fallback for IDs without detail data
  const program = PROGRAMS_MAP[id] ?? {
    id,
    name: 'Программа кредитования',
    rate: 8,
    minAmount: 500_000,
    maxAmount: 50_000_000,
    minTermMonths: 6,
    maxTermMonths: 36,
    activityTypes: ['mixed' as const],
    isActive: true,
    description: 'Описание программы загружается...',
    requirements: [],
  };

  return (
    <Screen>
      {/* Header */}
      <View className="bg-brand-green-500 px-4 pt-6 pb-6">
        <Pressable
          onPress={() => router.back()}
          className="flex-row items-center gap-1 mb-4 self-start active:opacity-70"
        >
          <Ionicons name="chevron-back" size={18} color="rgba(255,255,255,0.8)" />
          <Text className="text-white/80 text-sm">Назад</Text>
        </Pressable>
        <View className="bg-white/15 rounded-xl px-3 py-1.5 self-start mb-3">
          <Text className="text-white text-sm font-bold">{program.rate}% годовых</Text>
        </View>
        <Text className="text-white text-xl font-bold leading-7">{program.name}</Text>
        {program.nameKz && (
          <Text className="text-brand-green-100 text-sm mt-1">{program.nameKz}</Text>
        )}
      </View>

      <ScrollView
        showsVerticalScrollIndicator={false}
        contentContainerStyle={{ padding: 16, gap: 16, paddingBottom: 32 }}
      >
        {/* Key terms */}
        <Card>
          <Text className="text-gray-800 font-bold mb-4">Условия займа</Text>
          <View className="gap-3">
            <View className="flex-row justify-between items-center">
              <Text className="text-gray-500 text-sm">Ставка</Text>
              <Text className="text-brand-green-500 font-bold">{program.rate}% годовых</Text>
            </View>
            <View className="h-px bg-gray-100" />
            <View className="flex-row justify-between items-center">
              <Text className="text-gray-500 text-sm">Минимальная сумма</Text>
              <Text className="text-gray-800 font-semibold text-sm">{fmt(program.minAmount)}</Text>
            </View>
            <View className="h-px bg-gray-100" />
            <View className="flex-row justify-between items-center">
              <Text className="text-gray-500 text-sm">Максимальная сумма</Text>
              <Text className="text-gray-800 font-semibold text-sm">{fmt(program.maxAmount)}</Text>
            </View>
            <View className="h-px bg-gray-100" />
            <View className="flex-row justify-between items-center">
              <Text className="text-gray-500 text-sm">Срок займа</Text>
              <Text className="text-gray-800 font-semibold text-sm">
                {program.minTermMonths} – {program.maxTermMonths} месяцев
              </Text>
            </View>
            <View className="h-px bg-gray-100" />
            <View className="flex-row justify-between items-center">
              <Text className="text-gray-500 text-sm">Направление</Text>
              <Text className="text-gray-800 font-semibold text-sm">
                {program.activityTypes.map((t) => ACTIVITY_LABELS[t] ?? t).join(', ')}
              </Text>
            </View>
          </View>
        </Card>

        {/* Description */}
        {program.description && (
          <Card>
            <Text className="text-gray-800 font-bold mb-3">О программе</Text>
            <Text className="text-gray-600 text-sm leading-5">{program.description}</Text>
          </Card>
        )}

        {/* Requirements */}
        {program.requirements && program.requirements.length > 0 && (
          <Card>
            <Text className="text-gray-800 font-bold mb-3">Требования</Text>
            <View className="gap-2.5">
              {program.requirements.map((req, i) => (
                <View key={i} className="flex-row gap-3 items-start">
                  <View className="w-5 h-5 bg-brand-green-50 rounded-full items-center justify-center mt-0.5">
                    <Ionicons name="checkmark" size={12} color="#1a5c36" />
                  </View>
                  <Text className="text-gray-600 text-sm flex-1 leading-5">{req}</Text>
                </View>
              ))}
            </View>
          </Card>
        )}

        {/* CTA */}
        <View className="gap-3">
          <Button label="Подать заявку" size="lg" onPress={() => router.push('/(tabs)/programs')} />
          <Button
            label="Рассчитать платёж"
            variant="ghost"
            size="lg"
            onPress={() => router.push('/(tabs)/programs')}
          />
        </View>
      </ScrollView>
    </Screen>
  );
}
