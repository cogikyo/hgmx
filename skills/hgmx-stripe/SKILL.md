---
name: hgmx-stripe
description: Plan, build, test, or review Stripe payments in an hgmx site, including hosted Checkout, webhook fulfillment, refund reconciliation, restricted keys, CSP, and launch setup.
---

# Stripe payments in an hgmx site

Use this checklist for a site's own-account hosted Checkout integration.
Read [payments](references/payments.md) and [security](references/security.md) before changing payment code.
They are verbatim Stripe references; [UPSTREAM.md](UPSTREAM.md) records their revision, license, and the dated site lessons below.
For subscriptions, Connect, or tax, fetch the relevant current Stripe guidance rather than treating this checklist as complete coverage.

## Checkout

- [ ] **Create hosted Checkout Sessions server-side.** Keep prices and order state under server control.
- [ ] Instantiate a client with `stripe.NewClient`; do not set the global `stripe.Key`.
- [ ] Omit **`payment_method_types`** from Checkout requests.
      Choose methods in the Dashboard or through the Payment Method Configurations API.
- [ ] Use stable idempotency keys for retries of the same create operation.
- [ ] On API versions from `2026-03-25.dahlia`, set `integration_identifier` when creating sessions.
      Stripe's upstream guidance calls for a custom label with an eight-random-letter suffix.
- [ ] Before enabling `automatic_tax`, confirm an active tax registration for the jurisdiction.
      Enabling the flag alone does not make Stripe collect tax.
- [ ] Add **`https://checkout.stripe.com`** to `server.Site.CSP`'s `form-action` override.
      Preserve `'self'` for local forms: overrides replace the base directive's sources.
      Chromium checks the destination of the form POST → 303 redirect too.
      Do not add Stripe.js or embedded-UI CSP sources unless the site actually loads those assets.

## Webhooks and fulfillment

- [ ] Verify the signature against the **raw request body** before processing an event.
      Keep the endpoint's signing secret in the deployment's secret store.
- [ ] Pin the endpoint's **`api_version`** to the version pinned by the site's stripe-go SDK.
      As of **2026-11-02**, stripe-go **v87.1.0** pins **`2026-10-28.endive`**.
      Re-check the SDK's `stripe.APIVersion` when upgrading; do not bypass version validation.
- [ ] Create endpoints through the Webhook Endpoints API with an explicit `api_version` if automating setup.
      Capture the signing secret securely at creation; the API returns it only then.
- [ ] Drive fulfillment from **`checkout.session.completed`** and **`checkout.session.async_payment_succeeded`**.
      Retrieve the session fresh rather than trusting the event body's payment state.
      Require a complete session whose `payment_status` is not `unpaid` (`paid` or `no_payment_required`).
      The success page is not a fulfillment trigger.
- [ ] Make fulfillment **idempotent under concurrency**, including duplicate event deliveries.
      Protect the order transition and its side effects with database constraints and transactions.
- [ ] Handle `checkout.session.async_payment_failed` and `checkout.session.expired` without fulfilling.
      Reject events whose `livemode` disagrees with the configured environment.

## Refund lifecycle

- [ ] Store **each refund by its Stripe ID**, including payment, amount, status, and failure reason.
      A payment-level refunded total alone cannot track failed or partial refunds.
- [ ] Handle **`refund.created`**, **`refund.updated`**, and **`refund.failed`**.
      `charge.refunded` is not needed for this refund-ledger design.
- [ ] Reconcile by **listing all refunds for the payment**, including Dashboard-created refunds.
      Do not replay an idempotent create to observe state: API v1 returns the cached first response while the key is retained.
      Reusing a pruned key can create a new refund.
- [ ] **Serialize observe + save per payment**, including webhook and polling workers.
      Under that lock, save fresh observations, compute the uncovered remainder, and create only that amount.
      Count `pending`, `requires_action`, and `succeeded` refunds as coverage; surface failed or canceled refunds for attention.
- [ ] Keep polling for late Dashboard refunds and late failures, even after observing success.
      Record when Stripe was observed, and stop only after an observation covers the chosen window.
      As of **2026-11-02**, Stripe documents failures up to 30 days after posting.
      Choose and document the site's window; merely reaching its end in wall-clock time is not a final observation.

## Keys and agent access

- [ ] Prefer **restricted server keys** with the least permissions needed.
      `rk_live_` works in place of `sk_live_`; do not reject it merely for its prefix.
      The **2026-11 sandbox-tested permission set** is Checkout Sessions write, Payment Intents read, and Charges and Refunds write.
      Products, Prices, and Shipping Rates write may be needed for inline `price_data`; **this is unverified**.
      Test the site's actual requests with a restricted sandbox key before assigning live permissions.
- [ ] Keep test keys in local development and isolated test environments.
      **Production must refuse `sk_test_` and `rk_test_`.** Never commit or log any key or expose it to a browser.
- [ ] Use separate sandboxes for new local and CI integrations where supported.
      Keep live credentials out of these environments and use the deployment's secret store in production.
- [ ] **Do not connect `mcp.stripe.com` to a live account by default.**
      Its `stripe_api_write` tool supports arbitrary Stripe API POST, PATCH, PUT, and DELETE methods within granted permissions.
      If the user requests live MCP access, agree on its scope and require tool-call approval before connecting.
- [ ] Treat the **December 2026 cutoff as MCP authentication only**.
      Stripe's notice, checked **2026-11-02**, requires OAuth or Agent-tagged API keys after November 30.
      It does not retire normal secret or restricted keys for server API requests.

## Launch checks and Dashboard work

These setup limits are **2026-11 observations**; re-check current Stripe support before automating them.

- [ ] Assign **restricted-key creation, account activation, receipt switches, and public details** to Dashboard work.
      Do not claim these are scripted or complete because Checkout requests succeed.
- [ ] Set own-account branding in the Dashboard by default.
      A public-preview **`stripe tools update_brand_settings`** operation exists but is **live-only**.
      Check the installed tools plugin's help and current [CLI docs](https://docs.stripe.com/cli/tools/operations.md) before use.
      It is not a sandbox branding workflow; get approval before changing live branding.
- [ ] Test delayed payment success and failure, duplicate/concurrent delivery, partial refunds, and Dashboard refund failures.
      Verify the browser's actual redirect under the site's CSP.
      Report sandbox results as sandbox results, not proof of live account activation or permissions.

One implementation example is **alrodriguez.studio at `c762bbe`**, in repository-relative `shop/settle.go` and `store/refunds.go`.
If that repository is available, inspect its fresh-session settlement and serialized refund ledger; access to it is not required to use this checklist.
The site lessons were reported as sandbox-tested in November 2026, not re-tested by this skill's documentation update.

If `SITE.md` sits beside this file, read it; its rules override this skill.
