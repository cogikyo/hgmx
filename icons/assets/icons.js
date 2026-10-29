import { record } from "./motion.js";

const reduce = matchMedia("(prefers-reduced-motion: reduce)");
const drawn = ["pathLength", "stroke-dasharray", "stroke-dashoffset"];
const authored = new WeakMap();
const latest = new WeakMap();

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

document.addEventListener("hovered", async (e) => {
  const svg = e.target.closest?.("svg[data-icon]");
  if (!svg || reduce.matches) return;
  const { default: play } = await import(`./anim/${svg.dataset.icon}.js`);
  for (const el of svg.querySelectorAll("*")) {
    if (!authored.has(el)) authored.set(el, drawn.map((name) => el.getAttribute(name)));
  }
  const run = record(play, svg);
  latest.set(svg, run);
  await Promise.all(run.map((controls) => controls.finished));
  if (latest.get(svg) === run) restore(svg);
});
