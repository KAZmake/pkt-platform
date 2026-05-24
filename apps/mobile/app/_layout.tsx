import '~/global.css';
import { useEffect } from 'react';
import { View } from 'react-native';
import { Stack } from 'expo-router';
import { QueryClient } from '@tanstack/react-query';
import { PersistQueryClientProvider } from '@tanstack/react-query-persist-client';
import { AuthProvider } from '~/lib/auth/context';
import { asyncStoragePersister } from '~/lib/storage/persister';
import { OfflineBanner } from '~/components/shared/OfflineBanner';
import {
  addNotificationReceivedListener,
  addNotificationResponseListener,
} from '~/lib/notifications/service';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      retry: 2,
      gcTime: 30 * 60 * 1000,
    },
  },
});

export default function RootLayout() {
  useEffect(() => {
    const received = addNotificationReceivedListener((_notification) => {
      // Notification received while app is in foreground — handled by service handler
    });
    const response = addNotificationResponseListener((_response) => {
      // User tapped notification — navigate to relevant screen
      // TODO: use router.push based on response.notification.request.content.data
    });
    return () => {
      received.remove();
      response.remove();
    };
  }, []);

  return (
    <PersistQueryClientProvider
      client={queryClient}
      persistOptions={{ persister: asyncStoragePersister }}
    >
      <AuthProvider>
        <View style={{ flex: 1 }}>
          <OfflineBanner />
          <Stack screenOptions={{ headerShown: false }}>
            <Stack.Screen name="(tabs)" />
            <Stack.Screen name="(auth)" />
            <Stack.Screen name="+not-found" />
          </Stack>
        </View>
      </AuthProvider>
    </PersistQueryClientProvider>
  );
}
