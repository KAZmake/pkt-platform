# Mobile Distribution Setup

## Prerequisites

### EAS (Expo Application Services)

1. Create Expo account at expo.dev
2. Install EAS CLI: `npm install -g eas-cli`
3. Login: `eas login`
4. Link project: `cd apps/mobile && eas project:init`
   - This sets `extra.eas.projectId` in `app.json`
5. Create `EXPO_TOKEN` in GitHub Secrets:
   - expo.dev → Account Settings → Access Tokens → Create

### iOS (TestFlight)

1. Apple Developer Program membership required ($99/yr)
2. In App Store Connect: create App ID `kz.pkt.mobile`
3. Configure code signing: `eas credentials --platform ios`
   - EAS manages certificates and provisioning profiles
4. Add secrets to GitHub:
   - `EXPO_APPLE_APP_SPECIFIC_PASSWORD` — generate at appleid.apple.com
5. Update `eas.json` → `submit.production.ios`:
   - `ascAppId` — App Store Connect App ID (numeric)
   - `appleTeamId` — from developer.apple.com

### Android (Firebase App Distribution)

1. Create Firebase project at console.firebase.google.com
2. Add Android app with package `kz.pkt.mobile`
3. Enable App Distribution in Firebase console
4. Create a testers group named `testers`
5. Get Firebase CI token: `npx firebase-tools login:ci`
6. Add secrets to GitHub:
   - `FIREBASE_TOKEN` — from login:ci
   - `FIREBASE_APP_ID_ANDROID` — Firebase console → App → App ID (1:xxx:android:xxx)
7. For production Play Store: create service account, download JSON,
   set path in `eas.json` → `submit.production.android.serviceAccountKeyPath`

---

## Workflow: staging build + distribution

```bash
# Via GitHub Actions UI (workflow_dispatch):
# Profile: staging
# Platform: all
# Submit: true

# Or locally (requires Apple/Firebase creds):
cd apps/mobile
eas build --profile staging --platform all
eas submit --profile staging --platform ios --latest
```

## Workflow: production release

```bash
# Bump version in app.json → version and ios.buildNumber / android.versionCode
# Then trigger workflow with profile: production, submit: true
```

---

## GitHub Secrets required

| Secret                             | Where to get                                    |
| ---------------------------------- | ----------------------------------------------- |
| `EXPO_TOKEN`                       | expo.dev → Account → Access Tokens              |
| `EXPO_APPLE_APP_SPECIFIC_PASSWORD` | appleid.apple.com → App-Specific Passwords      |
| `FIREBASE_TOKEN`                   | `npx firebase-tools login:ci`                   |
| `FIREBASE_APP_ID_ANDROID`          | Firebase Console → Project Settings → Your apps |

---

## Build profiles summary

| Profile       | iOS output      | Android output | Target                 |
| ------------- | --------------- | -------------- | ---------------------- |
| `development` | Simulator build | APK            | Local dev              |
| `preview`     | IPA (internal)  | APK            | Internal testing       |
| `staging`     | IPA (internal)  | APK            | TestFlight / Firebase  |
| `production`  | IPA (App Store) | AAB            | App Store / Play Store |

---

## Verification checklist

| #   | Check                                                     | Result |
| --- | --------------------------------------------------------- | ------ |
| 1   | `eas project:init` completed, `app.json` has `projectId`  | `[ ]`  |
| 2   | iOS credentials configured via `eas credentials`          | `[ ]`  |
| 3   | GitHub Secrets set (EXPO_TOKEN, APPLE_PASSWORD, FIREBASE) | `[ ]`  |
| 4   | Staging build triggered and completes on EAS dashboard    | `[ ]`  |
| 5   | iOS build appears in TestFlight → Internal Testing        | `[ ]`  |
| 6   | Android APK appears in Firebase App Distribution          | `[ ]`  |
| 7   | Tester installs app and logs in via Keycloak PKCE         | `[ ]`  |
| 8   | Push notifications received on physical device            | `[ ]`  |
