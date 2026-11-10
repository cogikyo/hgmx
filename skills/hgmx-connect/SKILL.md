---
name: hgmx-connect
description: Build, change, or review an hgmx site's native Connect-RPC API, protobuf contracts, access policies, generation, and installed-client compatibility.
---

# Native APIs in an hgmx site

Read the selected hgmx module's `CONVENTIONS.md`, especially Serving an App API, Native App Ownership, App Generation and Checks, and App Release Baselines.
Use the base **hgmx** skill to locate that module; inspect the site's wiring and `go doc` for the auth and generated APIs before editing.
Use **hgmx-ios** or **hgmx-android** for the client side.

## Contract and handler

- **Keep site contracts in `proto/<site>/api/v1/`**, under a distinct protobuf package.
  Leave `proto/hgmx/` and `buf.yaml` managed by hgmx.
- Put Connect adapters in `api/`, importing generated types, domain code and store, never routes or views.
  Share business rules at their domain/store owner instead of calling an HTML handler.
- Build `api` with `accounts.API()`, register services on `api.Server()`, and classify each procedure through `api.Allow`.
  Use generated procedure constants and mount through `api.Mount`, not direct Connect mounting.
- Keep browser cookies and HTML-only middleware out of native RPC.
  Let auth enforce app metadata, bearer scope, admission and minimum builds.
- Use explicit sign-in, linking and mutation operations with the authority each needs.
  Return authoritative snapshots and capabilities; keep presentation and drafts in native code.
- Keep changes additive within a supported API major.
  Added validation must preserve released-client behavior, even when the schema remains compatible.

## Generate and verify

Run `go tool hgmx gen`, then inspect the changed generated file sets and run `go tool hgmx check`.
Sites use hgmx's Go descriptors for its contracts and generate adjacent Go only for their own protos.
Do not edit generated files, add a second generator chain or use remote Buf plugins.
Build affected native apps separately: a Go check does not prove Swift or Kotlin compilation.

Inspect `apps/compatibility.json` before changing a shipped contract.
Preserve recorded release commits and supported majors; missing history is a prerequisite failure, not permission to skip a baseline.
The current retained-client-check gap blocks recorded releases; follow App Release Baselines rather than erasing release evidence.
For kit changes, use the repository's adoption and seeded-failure fixture described in README.

If `SITE.md` sits beside this file, read it; its rules override this skill.
