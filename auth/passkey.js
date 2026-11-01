const supported = () => typeof PublicKeyCredential?.parseRequestOptionsFromJSON === "function";

async function post(url, body, signal) {
  const res = await fetch(url, {
    method: "POST",
    headers: { Accept: "application/json", "Content-Type": "application/json" },
    body,
    signal,
  });
  if (res.status === 204) return {};
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.problem || "Something went wrong. Try again.");
  return data;
}

async function signIn(mediation, signal) {
  const { publicKey } = await post("/auth/passkey/begin", null, signal);
  const credential = await navigator.credentials.get({
    publicKey: PublicKeyCredential.parseRequestOptionsFromJSON(publicKey),
    mediation,
    signal,
  });
  await post("/auth/passkey/finish", JSON.stringify(credential.toJSON()));
}

async function register() {
  const { publicKey } = await post("/auth/account/passkeys/begin");
  const credential = await navigator.credentials.create({
    publicKey: PublicKeyCredential.parseCreationOptionsFromJSON(publicKey),
  });
  await post("/auth/account/passkeys/finish", JSON.stringify(credential.toJSON()));
}

function show(el, err) {
  const problem = el.closest("[data-auth]")?.querySelector("[data-auth-problem]");
  if (!problem) return;
  problem.textContent =
    err.name === "NotAllowedError" ? "The passkey prompt was closed. Try again." : err.message;
}

let autofill;

async function scan() {
  const input = document.querySelector('input[autocomplete~="webauthn"]');
  if (autofill?.input === input) return;
  autofill?.abort();
  autofill = undefined;
  if (!input || !supported()) return;
  const controller = new AbortController();
  controller.input = input;
  autofill = controller;
  try {
    if (!(await PublicKeyCredential.isConditionalMediationAvailable?.()) || controller.signal.aborted) return;
    await signIn("conditional", controller.signal);
    const button = input.closest("[data-auth]")?.querySelector("[data-auth-passkey]");
    location.assign(button?.dataset.next || "/");
  } catch (err) {
    if (!controller.signal.aborted && err.name !== "NotAllowedError") show(input, err);
  } finally {
    if (autofill === controller) autofill = undefined;
  }
}

document.addEventListener("click", async (event) => {
  const button = event.target.closest("[data-auth-passkey]");
  if (!button || button.disabled) return;
  if (!supported()) {
    show(button, new Error("This browser can't use passkeys here. Update it or use another one."));
    return;
  }
  autofill?.abort();
  autofill = undefined;
  button.disabled = true;
  try {
    if (button.dataset.authPasskey === "register") await register();
    else await signIn();
    location.assign(button.dataset.next || "/");
  } catch (err) {
    show(button, err);
    button.disabled = false;
    scan();
  }
});

// Capture runs before htmx snapshots the form values.
document.addEventListener(
  "submit",
  (event) => {
    const token = event.target.querySelector("[data-auth-link-token]");
    if (token) token.value = location.hash.slice(1);
  },
  true,
);

document.addEventListener("htmx:after:process", scan);
scan();
