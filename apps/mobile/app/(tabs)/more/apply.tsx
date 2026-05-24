import { ActivityIndicator, Pressable, Text, View } from 'react-native';
import { useState } from 'react';
import { useRouter } from 'expo-router';
import { WebView } from 'react-native-webview';
import { Ionicons } from '@expo/vector-icons';

const WEB_URL = process.env.EXPO_PUBLIC_WEB_URL ?? 'http://localhost:3000';

export default function ApplyScreen() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(true);
  const [hasError, setHasError] = useState(false);

  return (
    <View style={{ flex: 1, backgroundColor: 'white' }}>
      {/* Header */}
      <View className="bg-brand-green-500 px-4 pt-12 pb-4 flex-row items-center gap-3">
        <Pressable
          onPress={() => router.back()}
          className="w-8 h-8 bg-white/20 rounded-full items-center justify-center"
        >
          <Ionicons name="chevron-back" size={18} color="white" />
        </Pressable>
        <View>
          <Text className="text-white text-base font-bold">Подать заявку</Text>
          <Text className="text-brand-green-100 text-xs">Заявка на получение займа</Text>
        </View>
      </View>

      {isLoading && (
        <View className="absolute inset-0 top-20 items-center justify-center bg-white">
          <ActivityIndicator size="large" color="#1a5c36" />
          <Text className="text-gray-400 text-sm mt-3">Загрузка формы...</Text>
        </View>
      )}

      {hasError ? (
        <View className="flex-1 items-center justify-center px-8 gap-4">
          <Ionicons name="wifi-outline" size={48} color="#d1d5db" />
          <Text className="text-gray-700 text-base font-bold text-center">Нет подключения</Text>
          <Text className="text-gray-400 text-sm text-center">
            Форма заявки требует подключения к интернету.
          </Text>
          <Pressable
            className="bg-brand-green-500 rounded-xl px-6 py-3 active:opacity-80"
            onPress={() => setHasError(false)}
          >
            <Text className="text-white font-semibold">Повторить</Text>
          </Pressable>
        </View>
      ) : (
        <WebView
          source={{ uri: `${WEB_URL}/apply` }}
          style={{ flex: 1 }}
          onLoadStart={() => setIsLoading(true)}
          onLoadEnd={() => setIsLoading(false)}
          onError={() => {
            setHasError(true);
            setIsLoading(false);
          }}
          startInLoadingState={false}
          javaScriptEnabled
          domStorageEnabled
          // Forward auth token for pre-filled form
          injectedJavaScriptBeforeContentLoaded={`
            window.__PKT_MOBILE__ = true;
          `}
        />
      )}
    </View>
  );
}
