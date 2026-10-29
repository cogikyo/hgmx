import "./theme.js";

const view = { yaw: -0.7, pitch: 0.5 };
const spread = 2;
const turn = 0.12;
const scenes = new WeakMap();
const scrollers = new WeakSet();
const slack = 16;
let focus = null;

function linear(v) {
  return v <= 0.04045 ? v / 12.92 : ((v + 0.055) / 1.055) ** 2.4;
}

function oklab(r, g, b) {
  const l = Math.cbrt(0.4122214708 * r + 0.5363325363 * g + 0.0514459929 * b);
  const m = Math.cbrt(0.2119034982 * r + 0.6806995451 * g + 0.1073969566 * b);
  const s = Math.cbrt(0.0883024619 * r + 0.2817188376 * g + 0.6299787005 * b);
  return [
    0.2104542553 * l + 0.793617785 * m - 0.0040720468 * s,
    1.9779984951 * l - 2.428592205 * m + 0.4505937099 * s,
    0.0259040371 * l + 0.7827717662 * m - 0.808675766 * s,
  ];
}

function parse(hex) {
  const n = parseInt(hex.slice(1), 16);
  return [(n >> 16) & 255, (n >> 8) & 255, n & 255].map((v) => v / 255);
}

const corners = [
  [0, 0, 0], [1, 0, 0], [1, 1, 0], [0, 1, 0],
  [0, 0, 1], [1, 0, 1], [1, 1, 1], [0, 1, 1],
];
const edges = [
  [0, 1], [1, 2], [2, 3], [3, 0], [4, 5], [5, 6], [6, 7], [7, 4], [0, 4], [1, 5], [2, 6], [3, 7],
];

const gamut = edges.map(([a, b]) => {
  const out = [];
  for (let i = 0; i <= 24; i++) {
    const t = i / 24;
    const rgb = corners[a].map((v, k) => v + (corners[b][k] - v) * t);
    out.push({ lab: oklab(...rgb.map(linear)), rgb });
  }
  return out;
});

function setup(canvas) {
  const data = JSON.parse(canvas.dataset.oklab);
  const points = data.points.map((p) => ({ ...p, lab: oklab(...parse(p.hex).map(linear)) }));
  const scene = { canvas, points, base: data.base, ink: data.ink };
  scenes.set(canvas, scene);
  new ResizeObserver(() => draw(scene)).observe(canvas);
  let drag = null;
  canvas.addEventListener("pointerdown", (e) => {
    drag = { x: e.clientX, y: e.clientY };
    canvas.setPointerCapture(e.pointerId);
  });
  canvas.addEventListener("pointerup", () => (drag = null));
  canvas.addEventListener("pointercancel", () => (drag = null));
  canvas.addEventListener("pointermove", (e) => {
    if (!drag) return;
    rotate((e.clientX - drag.x) * 0.01, (e.clientY - drag.y) * 0.01);
    drag = { x: e.clientX, y: e.clientY };
    draw(scene);
  });
  canvas.addEventListener("keydown", (e) => {
    const d = { ArrowLeft: [-turn, 0], ArrowRight: [turn, 0], ArrowUp: [0, -turn], ArrowDown: [0, turn] }[e.key];
    if (!d) return;
    e.preventDefault();
    rotate(...d);
    draw(scene);
  });
}

function rotate(yaw, pitch) {
  view.yaw += yaw;
  view.pitch = Math.max(-1.5, Math.min(1.5, view.pitch + pitch));
}

function project(scene, [L, a, b]) {
  const { w, h, scale } = scene;
  const cy = Math.cos(view.yaw), sy = Math.sin(view.yaw);
  const cp = Math.cos(view.pitch), sp = Math.sin(view.pitch);
  const x = (a * cy - b * sy) * spread;
  const z = (a * sy + b * cy) * spread;
  const y = L - 0.5;
  return { x: w / 2 + x * scale, y: h / 2 - (y * cp - z * sp) * scale, z: y * sp + z * cp };
}

function line(ctx, scene, from, to) {
  const p = project(scene, from), q = project(scene, to);
  ctx.beginPath();
  ctx.moveTo(p.x, p.y);
  ctx.lineTo(q.x, q.y);
  ctx.stroke();
  return q;
}

function draw(scene) {
  const { canvas } = scene;
  const dpr = window.devicePixelRatio || 1;
  const w = canvas.clientWidth, h = canvas.clientHeight;
  if (!w || !h) return;
  canvas.width = Math.round(w * dpr);
  canvas.height = Math.round(h * dpr);
  Object.assign(scene, { w, h, scale: Math.min(w, h) * 0.78 });
  const ctx = canvas.getContext("2d");
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  ctx.fillStyle = scene.base;
  ctx.fillRect(0, 0, w, h);

  ctx.lineWidth = 1.5;
  ctx.globalAlpha = 0.55;
  for (const edge of gamut) {
    for (let i = 1; i < edge.length; i++) {
      const [r, g, b] = edge[i].rgb.map((v) => Math.round(v * 255));
      ctx.strokeStyle = `rgb(${r} ${g} ${b})`;
      line(ctx, scene, edge[i - 1].lab, edge[i].lab);
    }
  }

  ctx.globalAlpha = 0.7;
  ctx.strokeStyle = ctx.fillStyle = scene.ink;
  ctx.lineWidth = 1;
  ctx.font = "11px system-ui, sans-serif";
  const axes = [
    ["L", [0, 0, 0], [1, 0, 0]],
    ["a", [0, -0.2, 0], [0, 0.2, 0]],
    ["b", [0, 0, -0.2], [0, 0, 0.2]],
  ];
  for (const [label, from, to] of axes) {
    const end = line(ctx, scene, from, to);
    ctx.fillText(`+${label}`, end.x + 4, end.y + 4);
  }

  const screen = scene.points.map((p) => ({ p, s: project(scene, p.lab) })).sort((a, b) => a.s.z - b.s.z);
  for (const { p, s } of screen) {
    const lit = !focus || p.families.includes(focus);
    ctx.globalAlpha = lit ? 1 : 0.2;
    ctx.beginPath();
    ctx.arc(s.x, s.y, lit && focus ? 6 : 4.5, 0, Math.PI * 2);
    ctx.fillStyle = p.hex;
    ctx.fill();
    ctx.strokeStyle = scene.ink;
    ctx.lineWidth = 0.75;
    ctx.stroke();
  }
  ctx.globalAlpha = 1;
}

function select(swatch) {
  const root = swatch.closest(".tx-view");
  if (!root) return;
  for (const s of root.querySelectorAll(".tx-swatch[aria-pressed]")) s.removeAttribute("aria-pressed");
  swatch.setAttribute("aria-pressed", "true");
  const detail = JSON.parse(swatch.dataset.detail);
  for (const el of root.querySelectorAll("[data-readout]")) {
    const key = el.dataset.readout;
    if (key === "swatch") el.setAttribute("fill", detail.hex);
    else if (key === "hex") el.firstElementChild.textContent = detail.hex;
    else el.textContent = detail[key];
  }
  focus = swatch.dataset.family;
  for (const dot of root.querySelectorAll(".tx-dot")) {
    dot.classList.toggle("tx-dim", !dot.dataset.family.split(" ").includes(focus));
  }
  for (const canvas of root.querySelectorAll("canvas[data-oklab]")) {
    const scene = scenes.get(canvas);
    if (scene) draw(scene);
  }
}

function reveal(el) {
  const scroller = el.closest(".tx-scroll");
  if (!scroller) return el.scrollIntoView({ block: "nearest" });
  const box = scroller.getBoundingClientRect(), r = el.getBoundingClientRect();
  if (r.right > box.right) scroller.scrollLeft += r.right - box.right + slack;
  else if (r.left < box.left) scroller.scrollLeft -= box.left - r.left + slack;
  const page = el.getBoundingClientRect();
  if (page.top < 0 || page.bottom > window.innerHeight) el.scrollIntoView({ block: "nearest", inline: "nearest" });
}

function move(swatch, key) {
  const group = swatch.closest("[data-swatches]");
  const row = Number(swatch.dataset.row), col = Number(swatch.dataset.col);
  const [r, c] = { ArrowUp: [row - 1, col], ArrowDown: [row + 1, col], ArrowLeft: [row, col - 1], ArrowRight: [row, col + 1] }[key];
  let next = group.querySelector(`.tx-swatch[data-row="${r}"][data-col="${c}"]`);
  if (!next && (key === "ArrowUp" || key === "ArrowDown")) {
    const cells = group.querySelectorAll(`.tx-swatch[data-row="${r}"]`);
    next = cells[cells.length - 1];
  }
  if (!next) return;
  swatch.setAttribute("tabindex", "-1");
  next.setAttribute("tabindex", "0");
  next.focus({ preventScroll: true });
  reveal(next);
  select(next);
}

document.addEventListener("click", (e) => {
  const swatch = e.target.closest?.(".tx-swatch");
  if (!swatch) return;
  for (const s of swatch.closest("[data-swatches]").querySelectorAll('.tx-swatch[tabindex="0"]')) s.setAttribute("tabindex", "-1");
  swatch.setAttribute("tabindex", "0");
  select(swatch);
});

document.addEventListener("keydown", (e) => {
  const swatch = e.target.closest?.(".tx-swatch");
  if (!swatch) return;
  if (e.key === "Enter" || e.key === " ") {
    e.preventDefault();
    select(swatch);
  } else if (e.key.startsWith("Arrow")) {
    e.preventDefault();
    move(swatch, e.key);
  }
});

function shade(scroller) {
  const wrap = scroller.parentElement;
  const max = scroller.scrollWidth - scroller.clientWidth;
  wrap.toggleAttribute("data-left", scroller.scrollLeft > slack);
  wrap.toggleAttribute("data-right", scroller.scrollLeft < max - slack);
}

function init(root) {
  if (!(root instanceof Element || root instanceof Document)) return;
  for (const s of root.querySelectorAll(".tx-scroll")) {
    if (scrollers.has(s)) continue;
    scrollers.add(s);
    s.addEventListener("scroll", () => shade(s), { passive: true });
    new ResizeObserver(() => shade(s)).observe(s);
  }
  for (const g of root.querySelectorAll("[data-swatches]")) {
    if (!g.querySelector('.tx-swatch[tabindex="0"]')) g.querySelector(".tx-swatch")?.setAttribute("tabindex", "0");
  }
  const found = [...root.querySelectorAll("canvas[data-oklab]")];
  if (root instanceof Element && root.matches("canvas[data-oklab]")) found.push(root);
  for (const canvas of found) {
    if (!scenes.has(canvas)) setup(canvas);
  }
  if (root.querySelector?.(".tx-view")) focus = null;
}

document.addEventListener("htmx:after:process", (e) => init(e.target));
init(document);
