import { record } from "./motion.js";

const reduce = matchMedia("(prefers-reduced-motion: reduce)");
const drawn = ["stroke-dasharray", "stroke-dashoffset"];
const authored = new WeakMap();
const latest = new WeakMap();
const plays = new Map();

const load = (name) => import(`./anim/${name}.js`).then(({ default: play }) => plays.set(name, play));

const near = (a, b) => Math.abs(a - b) < 1e-6;

function symmetric(el) {
  const m = new DOMMatrix(getComputedStyle(el).transform);
  const flat = [m.m13, m.m14, m.m23, m.m24, m.m31, m.m32, m.m34, m.m43].every((v) => near(v, 0)) && near(m.m33, 1) && near(m.m44, 1);
  const turn = flat && near(m.e, 0) && near(m.f, 0) && near(m.a, m.d) && near(m.b, -m.c) && near(m.a ** 2 + m.b ** 2, 1);
  if (!turn) return false;
  return near(m.a, 0) || near(m.b, 0) || (el.tagName === "circle" && !el.hasAttribute("transform-origin"));
}

function restore(svg) {
  for (const el of svg.querySelectorAll("*")) {
    authored.get(el)?.forEach((value, i) => {
      if (value === null) el.removeAttribute(drawn[i]);
      else el.setAttribute(drawn[i], value);
    });
    if (el.style.transform && symmetric(el)) el.style.removeProperty("transform");
  }
}

const ahead = new IntersectionObserver(
  (entries) => {
    for (const { isIntersecting, target } of entries) {
      if (!isIntersecting) continue;
      ahead.unobserve(target);
      const { icon } = target.dataset;
      if (!plays.has(icon)) load(icon).catch(() => {});
    }
  },
  { rootMargin: "200px" },
);

function preload(root) {
  if (reduce.matches || document.readyState !== "complete") return;
  const svgs = [...root.querySelectorAll("svg[data-icon]")];
  if (root.matches?.("svg[data-icon]")) svgs.push(root);
  for (const svg of svgs) ahead.observe(svg);
}

addEventListener("load", () => preload(document));
preload(document);
document.addEventListener("htmx:after:process", (e) => preload(e.target));

document.addEventListener("hovered", async (e) => {
  const svg = e.target.closest?.("svg[data-icon]");
  if (!svg || reduce.matches) return;
  const { icon } = svg.dataset;
  if (!plays.has(icon)) await load(icon);
  for (const el of svg.querySelectorAll("*")) {
    if (!authored.has(el)) authored.set(el, drawn.map((name) => el.getAttribute(name)));
  }
  const run = record(plays.get(icon), svg);
  latest.set(svg, run);
  await Promise.all(run.map((controls) => controls.finished));
  if (latest.get(svg) === run) restore(svg);
});
