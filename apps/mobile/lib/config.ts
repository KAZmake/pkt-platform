export const config = {
  apiUrl: process.env.EXPO_PUBLIC_API_URL ?? 'http://localhost:8080',
  syncUrl: process.env.EXPO_PUBLIC_SYNC_URL ?? 'http://localhost:8082',
  assistantUrl: process.env.EXPO_PUBLIC_ASSISTANT_URL ?? 'http://localhost:8083',
  keycloakUrl: process.env.EXPO_PUBLIC_KEYCLOAK_URL ?? 'http://localhost:8180',
  keycloakRealm: process.env.EXPO_PUBLIC_KEYCLOAK_REALM ?? 'pkt',
  keycloakClientId: process.env.EXPO_PUBLIC_KEYCLOAK_CLIENT_ID ?? 'pkt-mobile',
} as const;
