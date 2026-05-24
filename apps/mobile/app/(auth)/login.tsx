import { ActivityIndicator, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { useEffect } from 'react';
import { useAuth } from '~/lib/auth/context';
import { Button } from '~/components/ui/Button';

export default function LoginScreen() {
  const { session, login, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (session) router.replace('/(tabs)');
  }, [session]);

  if (isLoading) {
    return (
      <View className="flex-1 bg-brand-green-500 items-center justify-center">
        <ActivityIndicator size="large" color="white" />
      </View>
    );
  }

  return (
    <View className="flex-1 bg-brand-green-500">
      {/* Top decorative area */}
      <View className="flex-1 items-center justify-center px-8">
        <View className="w-24 h-24 bg-white rounded-3xl items-center justify-center mb-6 shadow-lg">
          <Text className="text-3xl font-black text-brand-green-500">ПКТ</Text>
        </View>
        <Text className="text-white text-2xl font-bold text-center leading-8">
          Первое кредитное{'\n'}товарищество
        </Text>
        <Text className="text-brand-green-100 text-sm mt-2">ЗКО, Казахстан</Text>
      </View>

      {/* Auth card */}
      <View className="bg-white rounded-t-3xl px-6 pt-8 pb-10">
        <Text className="text-xl font-bold text-gray-800 mb-1">Вход в личный кабинет</Text>
        <Text className="text-gray-500 text-sm mb-6">
          Используйте учётные данные ТОО «ПКТ» для входа
        </Text>
        <Button label="Войти через браузер" onPress={login} size="lg" />
        <Text className="text-gray-400 text-xs text-center mt-4">
          Безопасная авторизация OAuth 2.0 / PKCE
        </Text>
      </View>
    </View>
  );
}
