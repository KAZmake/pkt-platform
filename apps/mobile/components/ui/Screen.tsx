import { SafeAreaView, ScrollView, View, type ViewProps } from 'react-native';

interface ScreenProps extends ViewProps {
  scroll?: boolean;
  children: React.ReactNode;
}

export function Screen({ scroll = false, children, className, ...props }: ScreenProps) {
  if (scroll) {
    return (
      <SafeAreaView className="flex-1 bg-white">
        <ScrollView
          className="flex-1"
          contentContainerClassName="grow"
          showsVerticalScrollIndicator={false}
        >
          {children}
        </ScrollView>
      </SafeAreaView>
    );
  }
  return (
    <SafeAreaView className="flex-1 bg-white">
      <View className={`flex-1 ${className ?? ''}`} {...props}>
        {children}
      </View>
    </SafeAreaView>
  );
}
