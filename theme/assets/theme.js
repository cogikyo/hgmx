const endpoint = new URL("../palette.css", import.meta.url);
let committed = null;
let request = null;
let debounce;

function apply(name, root) {
  root.dataset.theme = name;
}

function cancel() {
  clearTimeout(debounce);
  if (!request) return;
  request.link.remove();
  request.done();
  request = null;
}

export function set(name, root = document.documentElement) {
  cancel();
  apply(name, root);
}

export function live(query, root = document.documentElement) {
  const url = new URL(endpoint);
  url.search = new URLSearchParams(query).toString();
  url.searchParams.set("name", "live");
  if (request?.href === url.href) return request.promise;
  cancel();
  if (committed?.href === url.href) {
    apply("live", root);
    return Promise.resolve();
  }
  const link = document.createElement("link");
  link.rel = "stylesheet";
  link.href = url.href;
  link.dataset.themeSheet = "";
  const mine = { href: url.href, link };
  mine.promise = new Promise((resolve, reject) => {
    mine.done = resolve;
    link.addEventListener("load", () => {
      if (request !== mine) return;
      request = null;
      committed?.remove();
      committed = link;
      apply("live", root);
      resolve();
    }, { once: true });
    link.addEventListener("error", () => {
      if (request !== mine) return;
      request = null;
      link.remove();
      reject(new Error(`theme: ${url.href} failed to load`));
    }, { once: true });
  });
  request = mine;
  document.head.append(link);
  return mine.promise;
}

function seed(form) {
  const q = new URLSearchParams();
  for (const [k, v] of new FormData(form)) {
    if (["base", "primary", "accent"].includes(k) && v) q.set(k, v);
  }
  return q;
}

document.addEventListener("click", (e) => {
  const el = e.target.closest("[data-theme-set], [data-theme-live]");
  if (!el || el.matches("form, select, input")) return;
  if ("themeSet" in el.dataset) set(el.dataset.themeSet);
  else live(el.dataset.themeLive);
});

document.addEventListener("change", (e) => {
  const el = e.target.closest("select[data-theme-set], input[data-theme-set]");
  if (el) set(el.value);
});

document.addEventListener("input", (e) => {
  const form = e.target.closest("form[data-theme-live]");
  if (!form) return;
  clearTimeout(debounce);
  debounce = setTimeout(() => live(seed(form)), 120);
});
