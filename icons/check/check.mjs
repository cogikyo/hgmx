import { mkdirSync, mkdtempSync, rmSync, writeFileSync } from "node:fs";
import { createRequire } from "node:module";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { parseArgs } from "node:util";

const { values: args } = parseArgs({
  options: {
    base: { type: "string", default: "http://127.0.0.1:18493" },
    only: { type: "string", default: "" },
    browsers: { type: "string", default: "chromium,firefox" },
    swaps: { type: "string", default: "5" },
    hovers: { type: "string", default: "2" },
    shots: { type: "string", default: "" },
    modules: { type: "string", default: process.env.NODE_MODULES ?? "node_modules" },
  },
});
if (!/^[1-9]\d*$/.test(args.hovers)) {
  console.error(`--hovers must be a positive integer, got ${JSON.stringify(args.hovers)}`);
  process.exit(2);
}
const require = createRequire(join(args.modules, "_"));
const puppeteer = require("puppeteer-core");

const executables = {
  chromium: { browser: "chrome", executablePath: "/usr/bin/chromium" },
  firefox: { browser: "firefox", executablePath: "/usr/bin/firefox-developer-edition" },
};
const drift = new Set(["orbit"]);
const snap = 16;
const moments = [0.1, 0.22, 0.35];
const viewport = { width: 1600, height: 1000 };

const probe = () => {
  const calls = [];
  const pending = [];
  const frozen = [];
  let run = null;
  let at = null;
  let gate = Promise.resolve();
  let release = () => {};
  const wrap = (animate) =>
    function (targets, ...rest) {
      const list = typeof targets === "string" ? [...document.querySelectorAll(targets)] : targets instanceof Element ? [targets] : targets?.[Symbol.iterator] ? [...targets] : [];
      if (!run) {
        run = { id: calls.length };
        queueMicrotask(() => (run = null));
      }
      const icon = list[0]?.closest?.("svg[data-icon]")?.dataset.icon ?? "?";
      calls.push({ icon, run: run.id, n: list.length });
      const controls = animate.call(this, targets, ...rest);
      if (at !== null) {
        controls.pause();
        controls.time = at;
        frozen.push(controls);
      }
      const finished = controls.finished ?? Promise.resolve();
      pending.push(finished);
      const held = gate;
      Object.defineProperty(controls, "finished", { configurable: true, value: finished.then(() => held) });
      return controls;
    };
  let real;
  Object.defineProperty(globalThis, "Motion", {
    configurable: true,
    get: () => real,
    set(obj) {
      let animate;
      Object.defineProperty(obj, "animate", {
        configurable: true,
        enumerable: true,
        get: () => animate,
        set: (fn) => (animate = wrap(fn)),
      });
      real = obj;
    },
  });
  const hold = () => {
    gate = new Promise((resolve) => (release = resolve));
  };
  const freeze = (t) => {
    at = t;
    for (const controls of frozen) controls.time = t;
  };
  const thaw = () => {
    at = null;
    for (const controls of frozen.splice(0)) controls.play();
  };
  globalThis.__icons = { calls, pending, hold, release: () => release(), freeze, thaw };
};

const snapshot = (names) => {
  const out = {};
  for (const name of names) {
    const svg = document.querySelector(`svg[data-icon="${name}"]`);
    if (!svg) {
      out[name] = null;
      continue;
    }
    out[name] = [...svg.querySelectorAll("*")].map((el) => {
      const cs = getComputedStyle(el);
      const dash = cs.strokeDasharray === "none" ? null : cs.strokeDasharray.split(/[ ,]+/).map(parseFloat);
      const length = el.getAttribute("pathLength");
      return {
        tag: el.tagName,
        opacity: parseFloat(cs.opacity),
        drawn: !dash || (parseFloat(cs.strokeDashoffset) === 0 && dash[0] >= 0.999 * (length ? parseFloat(length) : (el.getTotalLength?.() ?? 1))),
        fill: cs.fill,
        stroke: cs.stroke,
        visibility: cs.visibility,
      };
    });
  }
  return out;
};

const centers = (names) =>
  names.map((name) => {
    const host = document.querySelector(`svg[data-icon="${name}"]`)?.parentElement;
    if (!host) return null;
    const r = host.getBoundingClientRect();
    return [r.x + r.width / 2, r.y + r.height / 2];
  });

const settle = async () => {
  const { pending } = globalThis.__icons;
  const timeout = new Promise((resolve) => setTimeout(() => resolve("timeout"), 4000));
  while (pending.length) {
    const batch = pending.splice(0);
    if ((await Promise.race([Promise.all(batch), timeout])) === "timeout") return false;
  }
  await new Promise((resolve) => requestAnimationFrame(() => requestAnimationFrame(resolve)));
  return true;
};

function compare(before, after, name) {
  const errs = [];
  if (!after) return ["missing after"];
  if (before.length !== after.length) return [`element count ${before.length} -> ${after.length}`];
  before.forEach((b, i) => {
    const a = after[i];
    const at = `${b.tag}[${i}]`;
    if (Math.abs(b.opacity - a.opacity) > 0.01) errs.push(`${at} opacity ${b.opacity} -> ${a.opacity}`);
    if (!a.drawn) errs.push(`${at} stroke not fully drawn`);
    for (const k of ["fill", "stroke", "visibility"]) if (b[k] !== a[k]) errs.push(`${at} ${k} ${b[k]} -> ${a[k]}`);
  });
  return errs;
}

async function shots(page, names) {
  const out = {};
  for (const name of names) {
    const clip = await page.evaluate((name) => {
      const r = document.querySelector(`svg[data-icon="${name}"]`).getBoundingClientRect();
      return { x: Math.floor(r.x) - 4, y: Math.floor(r.y) - 4, width: Math.ceil(r.width) + 8, height: Math.ceil(r.height) + 8 };
    }, name);
    out[name] = await page.screenshot({ clip, encoding: "base64" });
  }
  return out;
}

const pixels = async (a, b) => {
  const load = async (s) => createImageBitmap(await (await fetch(`data:image/png;base64,${s}`)).blob());
  const [x, y] = await Promise.all([load(a), load(b)]);
  if (x.width !== y.width || x.height !== y.height) return { changed: Infinity, total: x.width * x.height };
  const read = (img) => {
    const c = new OffscreenCanvas(img.width, img.height).getContext("2d");
    c.drawImage(img, 0, 0);
    return c.getImageData(0, 0, img.width, img.height).data;
  };
  const p = read(x);
  const q = read(y);
  let changed = 0;
  let jumped = 0;
  for (let i = 0; i < p.length; i += 4) {
    const delta = Math.max(Math.abs(p[i] - q[i]), Math.abs(p[i + 1] - q[i + 1]), Math.abs(p[i + 2] - q[i + 2]));
    if (delta) changed++;
    if (delta > 128) jumped++;
  }
  return { changed, jumped, total: p.length / 4 };
};

async function hoverAll(page, names, sample = false) {
  const points = await page.evaluate(centers, names);
  const mark = await page.evaluate(() => globalThis.__icons.calls.length);
  if (sample) await page.evaluate((t) => globalThis.__icons.freeze(t), moments[0]);
  for (const [x, y] of points) {
    await page.mouse.move(x, y);
    await new Promise((r) => setTimeout(r, 15));
  }
  await page.mouse.move(1, 1);
  const frames = [];
  for (const t of sample ? moments : []) {
    await page.evaluate(async (t) => {
      globalThis.__icons.freeze(t);
      await new Promise((resolve) => requestAnimationFrame(() => requestAnimationFrame(resolve)));
    }, t);
    frames.push(await shots(page, names));
  }
  if (sample) await page.evaluate(() => globalThis.__icons.thaw());
  const settled = await page.evaluate(settle);
  const calls = await page.evaluate((m) => globalThis.__icons.calls.slice(m), mark);
  return { settled, calls, frames };
}

function runs(calls) {
  const by = new Map();
  for (const c of calls) {
    if (!by.has(c.icon)) by.set(c.icon, new Set());
    by.get(c.icon).add(c.run);
  }
  return by;
}

async function launch(name, reduce) {
  const dir = mkdtempSync(join(tmpdir(), `icons-${name}-`));
  const opts = { ...executables[name], headless: true, userDataDir: dir };
  if (name === "chromium") opts.args = ["--no-sandbox", ...(reduce ? ["--force-prefers-reduced-motion"] : [])];
  if (name === "firefox") opts.extraPrefsFirefox = { "ui.prefersReducedMotion": reduce ? 1 : 0 };
  const browser = await puppeteer.launch(opts);
  return { browser, dir };
}

async function open(browser, url, errors) {
  const page = await browser.newPage();
  await page.setViewport(viewport);
  page.on("console", (m) => m.type() === "error" && errors.push(`console: ${m.text()}`));
  page.on("pageerror", (e) => errors.push(`pageerror: ${e.message}`));
  page.on("requestfailed", (r) => errors.push(`requestfailed: ${r.url()}`));
  await page.evaluateOnNewDocument(probe);
  await page.goto(url, { waitUntil: "load" });
  await page.waitForFunction(() => globalThis.Motion?.animate && globalThis._hyperscript);
  const height = await page.evaluate(() => document.documentElement.scrollHeight);
  await page.setViewport({ ...viewport, height: Math.max(viewport.height, height) });
  return page;
}

async function check(name, url, wanted) {
  const failures = [];
  const fail = (icon, msg) => failures.push(`${name} ${icon}: ${msg}`);
  const errors = [];

  let { browser, dir } = await launch(name, false);
  try {
    const page = await open(browser, url, errors);
    const helper = await browser.newPage();
    await page.bringToFront();
    let worst = 0;
    const jumps = new Map();
    const names = await page.evaluate(() => [...document.querySelectorAll("svg[data-icon]")].map((s) => s.dataset.icon));
    for (const n of wanted ?? []) if (!names.includes(n)) fail(n, "not animated (no data-icon)");

    const hovers = Number(args.hovers);
    let first = null;
    let repeat = 0;
    for (let pass = 0; pass <= hovers; pass++) {
      const swapped = pass === hovers;
      if (swapped) {
        for (let i = 0; i < Number(args.swaps); i++) {
          const to = i % 2 === 0 ? "/b" : "/a";
          await page.click(to === "/b" ? "#b" : "#a");
          await page.waitForFunction((p) => document.querySelector("#page")?.textContent === p, {}, to);
        }
        await page.waitForFunction(() => globalThis._hyperscript && document.querySelector("svg[data-icon]"));
        await new Promise((r) => setTimeout(r, 300));
      }
      const label = swapped ? `after ${args.swaps} swaps` : `hover ${pass + 1}`;
      const before = await page.evaluate(snapshot, names);
      const look = await shots(page, names);
      await page.evaluate(() => globalThis.__icons.hold());
      const { settled, calls, frames } = await hoverAll(page, names, !swapped);
      if (!settled) fail("*", `${label}: animations did not finish within 4s`);
      first ??= frames;
      for (const n of pass && !swapped ? names : []) {
        if (drift.has(n)) continue;
        for (const [i, t] of moments.entries()) {
          const { changed, total } = await helper.evaluate(pixels, first[i][n], frames[i][n]);
          repeat = Math.max(repeat, changed);
          if (!changed) continue;
          fail(n, `${label}: differs from hover 1 at ${t}s: ${changed} of ${total} pixels`);
          if (args.shots) {
            mkdirSync(args.shots, { recursive: true });
            const stem = join(args.shots, `${name}-${n}-pass${pass + 1}-${t}s`);
            writeFileSync(`${stem}-first.png`, Buffer.from(first[i][n], "base64"));
            writeFileSync(`${stem}-repeat.png`, Buffer.from(frames[i][n], "base64"));
          }
          break;
        }
      }
      const quiet = await page.evaluate(() => globalThis.__icons.calls.length);
      const end = await shots(page, names);
      await page.evaluate(async () => {
        globalThis.__icons.release();
        await new Promise((resolve) => requestAnimationFrame(() => requestAnimationFrame(resolve)));
      });
      const after = await page.evaluate(snapshot, names);
      const relook = await shots(page, names);
      const stray = await page.evaluate((m) => globalThis.__icons.calls.slice(m).map((c) => c.icon), quiet);
      for (const icon of new Set(stray)) fail(icon, `${label}: animated while the checker took rest screenshots`);
      for (const n of names) {
        if (drift.has(n)) continue;
        const { jumped, changed: blurred } = await helper.evaluate(pixels, look[n], end[n]);
        const [j, b] = jumps.get(n) ?? [0, 0];
        jumps.set(n, [Math.max(j, jumped), Math.max(b, blurred)]);
        if (jumped > snap) fail(n, `${label}: snaps back after the animation ends: ${jumped} pixels jump`);
        const { changed, total } = await helper.evaluate(pixels, look[n], relook[n]);
        worst = Math.max(worst, changed);
        if (changed) {
          fail(n, `${label}: not at rest: ${changed} of ${total} pixels differ`);
          if (args.shots) {
            mkdirSync(args.shots, { recursive: true });
            const stem = join(args.shots, `${name}-${n}-pass${pass + 1}`);
            writeFileSync(`${stem}-rest.png`, Buffer.from(look[n], "base64"));
            writeFileSync(`${stem}-after.png`, Buffer.from(relook[n], "base64"));
          }
        }
      }
      const by = runs(calls);
      for (const n of names) {
        const count = by.get(n)?.size ?? 0;
        if (count !== 1) fail(n, `${label}: ${count} runs for one hover`);
        for (const e of compare(before[n], after[n], n)) fail(n, `${label}: not at rest: ${e}`);
      }
      for (const icon of by.keys()) if (!names.includes(icon)) fail(icon, `${label}: animate() outside a known icon`);
    }
    console.log(`${name}: ${names.length} icons hovered ${hovers} times, then once after ${args.swaps} boosted swaps; worst rest diff ${worst} pixels`);
    console.log(`${name}: worst repeat-hover diff at ${moments.join(", ")}s: ${repeat} pixels`);
    const top = [...jumps].sort((a, b) => b[1][0] - a[1][0] || b[1][1] - a[1][1]).slice(0, 5).map(([n, [j, b]]) => `${n} ${j} (${b} differ)`);
    console.log(`${name}: largest end-frame jumps (bound ${snap}): ${top.join(", ")}`);
  } finally {
    await browser.close();
    rmSync(dir, { recursive: true, force: true });
  }

  ({ browser, dir } = await launch(name, true));
  try {
    const page = await open(browser, url, errors);
    const reduced = await page.evaluate(() => matchMedia("(prefers-reduced-motion: reduce)").matches);
    if (!reduced) fail("*", "reduced motion emulation did not apply");
    const names = await page.evaluate(() => [...document.querySelectorAll("svg[data-icon]")].map((s) => s.dataset.icon));
    const { calls } = await hoverAll(page, names);
    await new Promise((r) => setTimeout(r, 300));
    const late = await page.evaluate(() => globalThis.__icons.calls.length);
    if (calls.length || late) fail("*", `reduced motion: ${late} animate() calls`);
    console.log(`${name}: reduced motion, ${names.length} icons hovered, ${late} animate() calls`);
  } finally {
    await browser.close();
    rmSync(dir, { recursive: true, force: true });
  }

  for (const e of errors) fail("*", e);
  return failures;
}

const wanted = args.only ? args.only.split(",").filter(Boolean) : null;
const url = `${args.base}/a${wanted ? `?only=${wanted.join(",")}` : ""}`;
let failures = [];
for (const name of args.browsers.split(",")) failures = failures.concat(await check(name, url, wanted));
for (const f of failures) console.log(`FAIL ${f}`);
console.log(failures.length ? `${failures.length} failures` : "ok");
process.exit(failures.length ? 1 : 0);
