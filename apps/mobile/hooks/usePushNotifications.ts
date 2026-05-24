import { useEffect, useRef, useState } from 'react';
import type * as Notifications from 'expo-notifications';
import {
  addNotificationReceivedListener,
  addNotificationResponseListener,
  registerForPushNotifications,
} from '~/lib/notifications/service';

export function usePushNotifications() {
  const [expoPushToken, setExpoPushToken] = useState<string | null>(null);
  const [lastNotification, setLastNotification] = useState<Notifications.Notification | null>(null);

  const receivedListener = useRef<Notifications.EventSubscription | null>(null);
  const responseListener = useRef<Notifications.EventSubscription | null>(null);

  useEffect(() => {
    registerForPushNotifications().then(setExpoPushToken);

    receivedListener.current = addNotificationReceivedListener((notification) => {
      setLastNotification(notification);
    });

    responseListener.current = addNotificationResponseListener((_response) => {
      // Handle tap on notification — navigate to relevant screen if needed
    });

    return () => {
      receivedListener.current?.remove();
      responseListener.current?.remove();
    };
  }, []);

  return { expoPushToken, lastNotification };
}
