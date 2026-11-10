import { mkdirSync, mkdtempSync, rmSync, writeFileSync } from "node:fs";
import { createRequire } from "node:module";
import { tmpdir } from "node:os";
import { join, resolve } from "node:path";
import { parseArgs } from "node:util";

const { values: args } = parseArgs({
  options: {
    base: { type: "string", default: "http://127.0.0.1:18493" },
    only: { type: "string", default: "" },
    out: { type: "string", default: "" },
    browser: { type: "string", default: "chromium" },
    size: { type: "string", default: "96" },
    modules: { type: "string", default: process.env.NODE_MODULES ?? "node_modules" },
  },
});
const names = args.only.split(",").map((s) => s.trim()).filter(Boolean);
const size = Number(args.size);
if (!names.length || !args.out || !(size >= 16)) {
  console.error("usage: node icons/check/film.mjs --only a,b --out <dir> [--base URL] [--modules <node_modules>] [--browser chromium|firefox] [--size 96]");
  process.exit(2);
}
const require = createRequire(join(resolve(args.modules), "_"));
const puppeteer = require("puppeteer-core");

const executables = {
  chromium: { browser: "chrome", executablePath: "/usr/bin/chromium", args: ["--no-sandbox"] },
  firefox: { browser: "firefox", executablePath: "/usr/bin/firefox-developer-edition" },
};
const times = [0, 0.05, 0.1, 0.15, 0.21, 0.27, 0.35, 0.42, 0.6];
const labels = [...times.map((t) => `${t}s`), "rest"];
const pad = Math.round(size * 0.375);
const cell = size + 2 * pad;

const probe = () => {
  const made = [];
  let at = 0;
  const wrap = (animate) =>
    function (targets, ...rest) {
      const first = targets instanceof Element ? targets : targets?.[0];
      const controls = animate.call(this, targets, ...rest);
      controls.pause();
      controls.time = at;
      made.push({ icon: first?.closest?.("svg[data-icon]")?.dataset.icon ?? "?", controls });
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
  const frame = () => new Promise((r) => requestAnimationFrame(() => requestAnimationFrame(r)));
  globalThis.__film = {
    count: (icon) => made.filter((m) => m.icon === icon).length,
    async seek(t) {
      at = t;
      for (const { controls } of made) controls.time = t;
      await frame();
    },
    async release() {
      for (const { controls } of made) controls.play();
      await Promise.all(made.map(({ controls }) => controls.finished));
      await frame();
      await new Promise((r) => setTimeout(r, 50));
      await frame();
    },
  };
};

const layout = (size, cell) => {
  const ul = document.querySelector("#icons ul");
  Object.assign(ul.style, { display: "flex", flexWrap: "wrap", gap: "0", width: `${cell * 6}px` });
  for (const li of ul.children) {
    Object.assign(li.style, { fontSize: `${size}px`, width: `${cell}px`, height: `${cell}px`, padding: "0", justifyContent: "center" });
    li.querySelector("span").style.display = "none";
  }
  window.scrollTo(0, 0);
  return [...ul.querySelectorAll("svg")].map((svg) => {
    const r = svg.getBoundingClientRect();
    return { icon: svg.dataset.icon ?? null, x: r.x + r.width / 2, y: r.y + r.height / 2 };
  });
};

const compose = async (frames, spots, names, labels, size, cell) => {
  const bitmaps = await Promise.all(frames.map(async (s) => createImageBitmap(await (await fetch(`data:image/png;base64,${s}`)).blob())));
  const left = 200;
  const head = 28;
  const encode = async (canvas) => {
    const blob = await canvas.convertToBlob({ type: "image/png" });
    const bytes = new Uint8Array(await blob.arrayBuffer());
    let bin = "";
    for (let i = 0; i < bytes.length; i += 0x8000) bin += String.fromCharCode(...bytes.subarray(i, i + 0x8000));
    return btoa(bin);
  };
  const strip = (ctx, row, y) => {
    ctx.fillStyle = "#222";
    ctx.font = "600 15px system-ui, sans-serif";
    ctx.textBaseline = "middle";
    ctx.fillText(names[row], 12, y + cell / 2);
    bitmaps.forEach((bmp, f) => {
      const x = left + f * cell;
      const { x: cx, y: cy } = spots[row];
      ctx.drawImage(bmp, cx - cell / 2, cy - cell / 2, cell, cell, x, y, cell, cell);
      ctx.strokeStyle = f === labels.length - 1 ? "#9ab" : "#ddd";
      ctx.strokeRect(x + 0.5, y + 0.5, cell - 1, cell - 1);
      const inner = (cell - size) / 2;
      ctx.strokeStyle = "rgba(0, 120, 255, 0.18)";
      ctx.strokeRect(x + inner + 0.5, y + inner + 0.5, size - 1, size - 1);
    });
  };
  const header = (ctx) => {
    ctx.fillStyle = "#555";
    ctx.font = "13px ui-monospace, monospace";
    ctx.textBaseline = "middle";
    ctx.textAlign = "center";
    labels.forEach((l, f) => ctx.fillText(l, left + f * cell + cell / 2, head / 2));
    ctx.textAlign = "start";
  };
  const sheet = (rows) => {
    const canvas = new OffscreenCanvas(left + labels.length * cell, head + rows.length * cell);
    const ctx = canvas.getContext("2d");
    ctx.fillStyle = "#fff";
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    header(ctx);
    rows.forEach((row, i) => strip(ctx, row, head + i * cell));
    return canvas;
  };
  const strips = [];
  for (let i = 0; i < names.length; i++) strips.push(await encode(sheet([i])));
  return { strips, sheet: await encode(sheet(names.map((_, i) => i))) };
};

const dir = mkdtempSync(join(tmpdir(), "icons-film-"));
const browser = await puppeteer.launch({ ...executables[args.browser], headless: true, userDataDir: dir });
try {
  const errors = [];
  const page = await browser.newPage();
  await page.setViewport({ width: cell * 6 + 64, height: 1000 });
  page.on("console", (m) => m.type() === "error" && errors.push(`console: ${m.text()}`));
  page.on("pageerror", (e) => errors.push(`pageerror: ${e.message}`));
  await page.evaluateOnNewDocument(probe);
  await page.goto(`${args.base}/a?only=${names.join(",")}`, { waitUntil: "load" });
  await page.waitForFunction(() => globalThis.Motion?.animate && globalThis._hyperscript);

  const spots = await page.evaluate(layout, size, cell);
  const height = await page.evaluate(() => document.documentElement.scrollHeight);
  await page.setViewport({ width: cell * 6 + 64, height: Math.max(1000, height) });
  const missing = names.filter((n, i) => spots[i]?.icon !== n);
  if (missing.length) throw new Error(`not animated (no module or no data-icon): ${missing.join(", ")}`);

  for (const [i, name] of names.entries()) {
    await page.mouse.move(spots[i].x, spots[i].y);
    await page.waitForFunction((n) => globalThis.__film.count(n) > 0, { timeout: 5000 }, name);
    await page.mouse.move(1, 1);
  }

  const frames = [];
  for (const t of times) {
    await page.evaluate((t) => globalThis.__film.seek(t), t);
    frames.push(await page.screenshot({ encoding: "base64" }));
  }
  await page.evaluate(() => globalThis.__film.release());
  frames.push(await page.screenshot({ encoding: "base64" }));
  if (errors.length) throw new Error(errors.join("\n"));

  const helper = await browser.newPage();
  const { strips, sheet } = await helper.evaluate(compose, frames, spots, names, labels, size, cell);
  mkdirSync(args.out, { recursive: true });
  names.forEach((name, i) => writeFileSync(join(args.out, `${name}.png`), Buffer.from(strips[i], "base64")));
  writeFileSync(join(args.out, "sheet.png"), Buffer.from(sheet, "base64"));
  console.log(`${names.length} strips and sheet.png in ${resolve(args.out)} (${args.browser}, ${size}px, frames at ${labels.join(" ")})`);
} finally {
  await browser.close();
  rmSync(dir, { recursive: true, force: true });
}
