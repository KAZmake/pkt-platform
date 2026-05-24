import { Link, Stack } from 'expo-router';
import { Text, View } from 'react-native';

export default function NotFoundScreen() {
  return (
    <>
      <Stack.Screen options={{ title: 'Не найдено' }} />
      <View className="flex-1 items-center justify-center p-8 bg-white">
        <Text className="text-5xl font-black text-brand-green-100 mb-4">404</Text>
        <Text className="text-xl font-bold text-gray-800 mb-2">Страница не найдена</Text>
        <Text className="text-gray-500 text-sm text-center mb-8">
          Запрошенная страница не существует.
        </Text>
        <Link href="/">
          <Text className="text-brand-green-500 font-semibold text-base">← На главную</Text>
        </Link>
      </View>
    </>
  );
}
