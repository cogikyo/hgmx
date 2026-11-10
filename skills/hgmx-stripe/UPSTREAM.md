# Stripe sources

Retrieved **2026-11-02** from Stripe's [skill index](https://docs.stripe.com/.well-known/skills/index.json) and its [stripe/ai mirror](https://github.com/stripe/ai/tree/main/skills).
The mirror revision is **`b3ab9149225f837894b7de08c1b14650d7cc64eb`**.

Only these three upstream content files were used:

- [`stripe-best-practices/SKILL.md`](https://github.com/stripe/ai/blob/b3ab9149225f837894b7de08c1b14650d7cc64eb/skills/stripe-best-practices/SKILL.md) is adapted into this directory's `SKILL.md`.
- [`references/payments.md`](https://github.com/stripe/ai/blob/b3ab9149225f837894b7de08c1b14650d7cc64eb/skills/stripe-best-practices/references/payments.md) is copied verbatim under `references/`.
- [`references/security.md`](https://github.com/stripe/ai/blob/b3ab9149225f837894b7de08c1b14650d7cc64eb/skills/stripe-best-practices/references/security.md) is copied verbatim under `references/`.

The MIT license permits redistribution with the copyright and permission notice retained below.
The references keep upstream wording and layout, including guidance for integrations outside hosted Checkout.
Read the local checklist for the hgmx-specific choices and fetch current guidance for other Stripe products.
The upstream skill's SDK table lists Go 87.0.0; the inspected site's v87.1.0 pin is a separate **2026-11** observation, not a claim about the latest SDK.

## Site lessons and API checks

The November 2026 lessons were supplied from sandbox work on **alrodriguez.studio**.
The example files `shop/settle.go` and `store/refunds.go` were inspected read-only at **`c762bbe`** to check fresh-session settlement, per-refund storage, payment locking, and observation windows.
The site's stripe-go pin and module-cache `api_version.go` confirm v87.1.0 and `2026-10-28.endive`.
No Stripe account was connected and no payment tests were run for this documentation update.
The reported restricted-key permissions and Dashboard-only limits were not independently exercised; the inline `price_data` permission additions remain unverified.

These current public sources were also read on **2026-11-02**:

- [Webhook endpoint creation](https://docs.stripe.com/api/webhook_endpoints/create.md) accepts an explicit `api_version`.
- [The endpoint object](https://docs.stripe.com/api/webhook_endpoints/object.md) returns `secret` only at creation.
- [API idempotency](https://docs.stripe.com/api/idempotent_requests.md) explains cached v1 responses and key pruning.
- [Refunds](https://docs.stripe.com/refunds.md) documents lifecycle events and failures up to 30 days after posting.
- [MCP](https://docs.stripe.com/mcp.md) documents the write tool's scope and the November 30 authentication deadline.
- [CLI tools operations](https://docs.stripe.com/cli/tools/operations.md) documents `stripe tools update_brand_settings`.
- [The tools plugin guide](https://docs.stripe.com/stripe-cli/tools-plugin.md) confirms branding operations are live-only.

The CLI pages show different invocation forms: direct operation subcommands versus `stripe tools execute <operation>`.
Check installed help before using the public-preview branding operation; this update did not execute it.

## Re-sync

Resolve a new mirror revision and re-check its license before updating these files.
Replace the two references byte-for-byte, then review the upstream skill's relevant changes against the local checklist.
Keep site lessons and qualifications unless new evidence changes them, and refresh the date and revision here.
Do not add upstream's other product references without a site task that needs them.

## Stripe license

From [`stripe/ai/LICENSE`](https://github.com/stripe/ai/blob/b3ab9149225f837894b7de08c1b14650d7cc64eb/LICENSE):

```text
The MIT License (MIT)

Copyright (c) 2024–2025 Stripe

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
