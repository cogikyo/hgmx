---
name: hgmx-ios
description: Build, change, or review a site-owned hgmx SwiftUI app, including screens, session calls, Keychain storage, identity, generation, and iPhone launch prerequisites.
---

# SwiftUI apps in an hgmx site

Read the selected hgmx module's `CONVENTIONS.md`: Native App Ownership, App Generation and Checks, and App Launch Requirements.
Use **hgmx** to locate the module and **hgmx-connect** when changing a service or contract.
Inspect the site's adopted files first; `apps/ios/` is a one-time copy, not a managed native library.

## Add a screen

1. Generate the contract through `go tool hgmx gen` from the site root.
   The local `API` package owns generated Swift under `API/Sources/API/Generated/`.
2. Construct the generated service client with `session.api.client`.
   Route requests through `Session.call`; keep auth grants and account snapshots on the session's dedicated paths.
3. Hold loading, snapshots, drafts and problems in an `@MainActor @Observable` view model.
   Use `Work` for busy/problem state and preserve task cancellation.
4. Render that model in a SwiftUI view and add its route in `Starter/Root.swift`.
   The Xcode project's synchronized folders admit new files under `Starter/` without manual project entries.

The kit's `starters/fixture/overlay/apps/ios/Starter/Home.swift` shows the complete Notes client/model/view path.
Render server capabilities instead of inferring permissions; apply mutation snapshots rather than maintaining a second business-state authority.

## Session, identity and builds

`Starter/StarterApp.swift` constructs the app-lifetime `Session`; `Starter/Api.swift` owns the configured client and metadata.
`Starter/TokenStore.swift` stores the bearer per API origin in non-synchronizing `WhenUnlockedThisDeviceOnly` Keychain storage.
Features do not read Keychain or attach tokens themselves; uninstalling is not session revocation.

Edit app identity and versions in `Config.xcconfig`, including bundle ID, display name, API host, RP ID and public app ID.
`Info.plist` consumes those values, and `Starter.entitlements` declares `webcredentials` from the RP ID.
Choose the real signing team in Xcode and match the server's provisioning prefix/bundle registration before device passkey testing.
Adoption does not supply signing or prove associated domains.

Open `Starter.xcodeproj` and use the `Starter` scheme after generation.
**Keep the `Starter.xcodeproj` name:** hgmx's pins check uses its fixed workspace `Package.resolved` path, so renaming the project breaks `hgmx check`.
Keep the workspace `Package.resolved` committed and use its locked versions for builds; changing dependencies also needs the hgmx pin check.
On non-Mac hosts, `HGMX_MAC` selects the SSH Mac for Swift generation only, not native compilation.
Report simulator builds, tests, physical passkeys and store readiness separately.

## Known gaps

The starter has no Google sign-in, and `init apps ios` rejects `--google-client-id`.
Apple sign-in is also absent; do not add a placeholder button or call the app buyer-launch ready.
Follow App Launch Requirements for the content owner's Apple account, provider revocation, deletion and physical-device prerequisites.
Preserve the documented foreground-refresh gap in any handoff until it is implemented and verified.
When server sign-out is unconfirmed, the starter shows a recovery notice pointing to `<origin>/auth/account` → **Devices**; preserve that notice rather than promising remote revocation.

If `SITE.md` sits beside this file, read it; its rules override this skill.
