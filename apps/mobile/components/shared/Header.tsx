import { Pressable, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';

interface HeaderProps {
  title: string;
  subtitle?: string;
  showBack?: boolean;
}

export function Header({ title, subtitle, showBack = true }: HeaderProps) {
  const router = useRouter();

  return (
    <View className="bg-white border-b border-gray-100 px-4 py-3 flex-row items-center gap-3">
      {showBack && (
        <Pressable
          onPress={() => router.back()}
          className="w-8 h-8 items-center justify-center rounded-full active:bg-gray-100"
        >
          <Ionicons name="chevron-back" size={20} color="#374151" />
        </Pressable>
      )}
      <View className="flex-1">
        <Text className="text-gray-800 text-base font-bold" numberOfLines={1}>
          {title}
        </Text>
        {subtitle && (
          <Text className="text-gray-400 text-xs mt-0.5" numberOfLines={1}>
            {subtitle}
          </Text>
        )}
      </View>
    </View>
  );
}
