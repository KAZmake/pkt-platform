import { Stack } from 'expo-router';

export default function MoreLayout() {
  return (
    <Stack screenOptions={{ headerShown: false }}>
      <Stack.Screen name="index" />
      <Stack.Screen name="about" />
      <Stack.Screen name="contacts" />
      <Stack.Screen name="faq" />
      <Stack.Screen name="news" />
    </Stack>
  );
}
