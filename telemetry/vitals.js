const endpoint = new URL("../", import.meta.url).pathname;
const send = (body) => navigator.sendBeacon(endpoint, JSON.stringify(body));
const traceparent = (entry) =>
  entry?.serverTiming?.find((t) => t.name === "traceparent")?.description ?? "";

const nav = performance.getEntriesByType("navigation")[0];
send({ type: "view", path: location.pathname, referrer: document.referrer, traceparent: traceparent(nav) });

document.addEventListener("htmx:after:history:push", () => {
  const entry = performance.getEntriesByName(location.href).at(-1);
  send({ type: "view", path: location.pathname, traceparent: traceparent(entry) });
});

const vitals = { type: "vitals", path: location.pathname, traceparent: traceparent(nav) };
if (nav) vitals.ttfb = Math.round(nav.responseStart);

const observe = (type, fn, options = {}) => {
  try {
    new PerformanceObserver((list) => list.getEntries().forEach(fn)).observe({ type, buffered: true, ...options });
  } catch {}
};

observe("paint", (e) => {
  if (e.name === "first-contentful-paint") vitals.fcp = Math.round(e.startTime);
});
observe("largest-contentful-paint", (e) => {
  vitals.lcp = Math.round(e.startTime);
});

let shift = 0;
let first = 0;
let last = 0;
observe("layout-shift", (e) => {
  if (e.hadRecentInput) return;
  if (shift && e.startTime - last < 1000 && e.startTime - first < 5000) {
    shift += e.value;
  } else {
    shift = e.value;
    first = e.startTime;
  }
  last = e.startTime;
  vitals.cls = Math.max(vitals.cls ?? 0, Math.round(shift * 10000) / 10000);
});

const interactions = new Map();
const interact = (e) => {
  const id = e.interactionId || `first:${e.startTime}`;
  interactions.set(id, Math.max(interactions.get(id) ?? 0, e.duration));
};
observe("event", (e) => e.interactionId && interact(e), { durationThreshold: 40 });
observe("first-input", interact);

let sent = false;
const flush = () => {
  if (sent) return;
  sent = true;
  const worst = [...interactions.values()].sort((a, b) => b - a);
  if (worst.length) vitals.inp = Math.round(worst[Math.min(worst.length - 1, Math.floor(worst.length / 50))]);
  send(vitals);
};
addEventListener("visibilitychange", () => document.visibilityState === "hidden" && flush(), true);
addEventListener("pagehide", flush, true);
