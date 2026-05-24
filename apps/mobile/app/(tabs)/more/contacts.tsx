import { Linking, Pressable, ScrollView, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Screen } from '~/components/ui/Screen';
import { Card } from '~/components/ui/Card';

type IoniconName = React.ComponentProps<typeof Ionicons>['name'];

const CONTACTS = [
  {
    icon: 'call-outline' as IoniconName,
    label: 'Call-центр',
    value: '8 800 080 1700',
    sub: 'Бесплатно, пн–пт 9:00–18:00',
    action: () => Linking.openURL('tel:88000801700'),
  },
  {
    icon: 'call-outline' as IoniconName,
    label: 'Приёмная',
    value: '+7 (7112) 50-00-00',
    action: () => Linking.openURL('tel:+77112500000'),
  },
  {
    icon: 'mail-outline' as IoniconName,
    label: 'Email',
    value: 'info@pkt.kz',
    action: () => Linking.openURL('mailto:info@pkt.kz'),
  },
  {
    icon: 'globe-outline' as IoniconName,
    label: 'Сайт',
    value: 'www.pkt.kz',
    action: () => Linking.openURL('https://pkt.kz'),
  },
];

const OFFICES = [
  { city: 'Уральск (головной офис)', address: 'ул. Есет батыра, 37', hours: 'Пн–Пт 9:00–18:00' },
  { city: 'Таскала', address: 'ул. Абая, 15', hours: 'Пн–Пт 9:00–17:00' },
  { city: 'Чапаево', address: 'ул. Ленина, 8', hours: 'Пн–Пт 9:00–17:00' },
];

export default function ContactsScreen() {
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
        <Text className="text-white text-xl font-bold">Контакты</Text>
        <Text className="text-brand-green-100 text-sm mt-1">Связаться с нами</Text>
      </View>

      <ScrollView
        showsVerticalScrollIndicator={false}
        contentContainerStyle={{ padding: 16, gap: 16, paddingBottom: 32 }}
      >
        {/* Contacts */}
        <Card className="p-0 overflow-hidden">
          {CONTACTS.map((c, i) => (
            <Pressable
              key={c.label}
              className={`flex-row items-center gap-3 px-4 py-3.5 active:bg-gray-50 ${i < CONTACTS.length - 1 ? 'border-b border-gray-100' : ''}`}
              onPress={c.action}
            >
              <View className="w-9 h-9 bg-brand-green-50 rounded-xl items-center justify-center">
                <Ionicons name={c.icon} size={18} color="#1a5c36" />
              </View>
              <View className="flex-1">
                <Text className="text-gray-400 text-xs">{c.label}</Text>
                <Text className="text-gray-800 text-sm font-semibold">{c.value}</Text>
                {c.sub && <Text className="text-gray-400 text-xs">{c.sub}</Text>}
              </View>
              <Ionicons name="open-outline" size={14} color="#d1d5db" />
            </Pressable>
          ))}
        </Card>

        {/* Social */}
        <Card>
          <Text className="text-gray-800 font-bold mb-3">Мы в соцсетях</Text>
          <View className="flex-row gap-3">
            {[
              { icon: 'logo-instagram' as IoniconName, label: 'Instagram', color: '#e1306c' },
              { icon: 'logo-facebook' as IoniconName, label: 'Facebook', color: '#1877f2' },
              { icon: 'chatbubbles-outline' as IoniconName, label: 'WhatsApp', color: '#25d366' },
            ].map((s) => (
              <Pressable
                key={s.label}
                className="flex-1 items-center gap-1.5 bg-gray-50 rounded-xl py-3 active:bg-gray-100"
              >
                <Ionicons name={s.icon} size={22} color={s.color} />
                <Text className="text-gray-600 text-xs">{s.label}</Text>
              </Pressable>
            ))}
          </View>
        </Card>

        {/* Offices */}
        <Card>
          <Text className="text-gray-800 font-bold mb-4">Офисы</Text>
          <View className="gap-4">
            {OFFICES.map((o) => (
              <View key={o.city} className="flex-row gap-3">
                <View className="w-8 h-8 bg-brand-green-50 rounded-lg items-center justify-center shrink-0 mt-0.5">
                  <Ionicons name="business-outline" size={16} color="#1a5c36" />
                </View>
                <View className="flex-1">
                  <Text className="text-gray-800 text-sm font-bold">{o.city}</Text>
                  <Text className="text-gray-500 text-sm">{o.address}</Text>
                  <Text className="text-gray-400 text-xs mt-0.5">{o.hours}</Text>
                </View>
              </View>
            ))}
          </View>
        </Card>
      </ScrollView>
    </Screen>
  );
}
