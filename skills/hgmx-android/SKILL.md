---
name: hgmx-android
description: Build, change, or review a site-owned hgmx Compose app, including screens, session calls, Keystore storage, identity, generated clients, and Android builds.
---

# Compose apps in an hgmx site

Read the selected hgmx module's `CONVENTIONS.md`: Native App Ownership, App Generation and Checks, and App Launch Requirements.
Use **hgmx** to locate the module and **hgmx-connect** when changing a service or contract.
Inspect the site's adopted files first; `apps/android/` is site code after adoption.

## Add a screen

1. Generate the contract through `go tool hgmx gen` from the site root.
   Java lite messages and Kotlin Connect clients live in `app/src/main/generated/`.
2. Construct the generated service client with `session.api.client` and call through `Session.call`.
   Use Java lite builders directly; do not introduce parallel DTOs or a Kotlin protobuf DSL.
3. Expose loading, snapshots, drafts and problems from a `ViewModel` through `StateFlow`.
   Launch screen work in `viewModelScope` so its lifetime follows the screen.
4. Collect with `collectAsStateWithLifecycle` in the composable and add navigation in `Root.kt`.
   Keep credential and storage effects outside composables.

The kit's `starters/fixture/overlay/apps/android/app/src/main/kotlin/com/example/fixture/notes/Home.kt` shows the Notes client/view-model/view path.
Render server capabilities and apply returned mutation snapshots; do not decide permissions from local state.
Use manual root construction rather than adding repository, use-case or DI scaffolding.

## Session, identity and builds

Under `app/src/main/kotlin/<site package>/`, `App.kt` constructs the app-lifetime `Session` and `Api.kt` owns the configured transport.
`TokenStore.kt` encrypts the bearer with an Android Keystore AES-GCM key and stores ciphertext plus nonce in `noBackupFilesDir/session`.
Features do not read storage or attach tokens; preserve session generation checks across sign-out and replacement.

Edit identity in `app/build.gradle.kts`: namespace/application ID, display name, API origin, public app ID, optional Google web client ID and versions.
Keep Kotlin package paths and declarations consistent if renaming after adoption.
`gradle.properties` records the starter source version; `gradle/libs.versions.toml` owns dependency versions.
Public client IDs belong in the app, never provider secrets.
Google sign-in needs matching server audience/presenter configuration and real provider-token verification before production.

Use JDK 21, the pinned SDK packages and the committed Gradle wrapper.
From `apps/android/`, run `./gradlew :app:lintDebug :app:assembleDebug :app:assembleDebugAndroidTest`; use `:app:connectedDebugAndroidTest` with a device or emulator for session tests.
Commit dependency locks and verification metadata, and verify dependency changes with a fresh Gradle cache so a warm cache cannot hide missing checksums.
Builds do not prove the installed signer, passkeys, associations or store readiness.
Preserve the documented foreground-refresh gap in handoffs until implemented and verified.
When server sign-out is unconfirmed, the starter shows a recovery notice pointing to `<origin>/auth/account` → **Devices**; preserve that notice rather than promising remote revocation.

If `SITE.md` sits beside this file, read it; its rules override this skill.
