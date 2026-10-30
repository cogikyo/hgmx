const endpoint = new URL("../palette.css", import.meta.url);
const root = document.documentElement;
let committed = null;
let pending = null;

function cancel() {
  pending?.remove();
  pending = null;
}

function set(name) {
  cancel();
  root.dataset.theme = name;
}

function live(query) {
  const url = new URL(endpoint);
  url.search = new URLSearchParams(query).toString();
  url.searchParams.set("name", "live");
  if (pending?.href === url.href) return;
  cancel();
  if (committed?.href === url.href) {
    root.dataset.theme = "live";
    return;
  }
  const link = document.createElement("link");
  link.rel = "stylesheet";
  link.href = url.href;
  link.dataset.themeSheet = "";
  link.addEventListener("load", () => {
    if (pending !== link) return;
    pending = null;
    committed?.remove();
    committed = link;
    root.dataset.theme = "live";
  }, { once: true });
  link.addEventListener("error", () => {
    if (pending !== link) return;
    pending = null;
    link.remove();
    console.warn(`theme: ${url.href} failed to load`);
  }, { once: true });
  pending = link;
  document.head.append(link);
}

document.addEventListener("click", (e) => {
  const el = e.target.closest("[data-theme-set], [data-theme-live]");
  if (!el || el.matches("select, input")) return;
  if ("themeSet" in el.dataset) set(el.dataset.themeSet);
  else live(el.dataset.themeLive);
});

document.addEventListener("change", (e) => {
  const el = e.target.closest("select[data-theme-set], input[data-theme-set]");
  if (el) set(el.value);
});
