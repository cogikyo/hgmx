import { existsSync, mkdirSync, mkdtempSync, readdirSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { basename, join, resolve } from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import { parseArgs } from "node:util";

const rules = {
  syntax: ["error", "§1 Mechanics: node --check must pass"],
  imports: ["error", "§1 Mechanics: import only the names you use from ../motion.js"],
  export: ["error", "§1 Mechanics: one module per icon, a default function of el"],
  deferred: ["error", "§1 Mechanics: call every animate() synchronously"],
  state: ["error", "§1 Mechanics: keep no state between hovers, except orbit's WeakMap"],
  data: ["error", "§1 Mechanics: never read or write data-* attributes"],
  scope: ["error", "§1 Mechanics: el is the root <svg>; select parts inside it only"],
  svg: ["error", "§1 Mechanics: anim/<name>.js animates svg/<name>.svg"],
  run: ["error", "§1 Mechanics: the default function runs against its SVG"],
  target: ["error", "§1 Selecting parts: every selector matches a part"],
  call: ["error", "§1 Mechanics: animate(target, keyframes, options)"],
  duration: ["error", "§2 Time: only τmax, τmid, τmin and 0.3 as durations"],
  delay: ["warning", "§2 Beats: delays normally come from the 0, τmin, τmid clock"],
  stagger: ["error", "§2 Burst cascades: stagger(s) with s from 0.03 to 0.21"],
  easing: ["error", "§2 Easing: no ease, type, repeat, repeatType or repeatDelay"],
  direction: ["error", "§4 Draw: do not use pathOffset"],
  keyframes: ["error", "§6 Return to Rest: start every keyframe array explicitly; opacity: [1] only"],
  rest: ["error", "§6 Return to Rest: scale 1, translation 0, opacity 1, pathLength 1"],
  turn: ["error", "§6 Return to Rest: end rotation at a turn that maps the part onto itself"],
  paint: ["error", "§6 Return to Rest: do not animate fill, stroke, stroke-width or colour"],
  symmetry: ["warning", "§6 Return to Rest: a 90° end needs four-fold symmetry"],
  beats: ["warning", "§2 Beats: at most three beats"],
  calls: ["warning", "§7 Unique but in the Family: 2–4 animate() calls"],
  cascade: ["warning", "§2 Burst cascades: the last item finishes by 0.6 s"],
  overshoot: ["warning", "§4 Burst: overshoot at most 1.75"],
  property: ["warning", "§4 Verbs: animated property outside the grammar's verbs"],
  option: ["warning", "§2 Time: option outside duration, delay, times and startDelay"],
};

const exempt = {
  orbit: new Set(["state", "turn", "symmetry"]),
  github: new Set(["cascade"]),
};

const deferrals = new Set(["setTimeout", "setInterval", "setImmediate", "requestAnimationFrame", "requestIdleCallback", "queueMicrotask", "await", "async", "Promise"]);
const continuations = new Set(["then", "catch", "finally"]);
const outside = new Set(["document", "window", "globalThis", "self", "Motion"]);
const easings = new Set(["ease", "type", "repeat", "repeatType", "repeatDelay"]);
const paint = /^(fill|stroke|color)([A-Z-]|$)/i;
const options = new Set(["duration", "delay", "times"]);
const staggerOptions = new Set(["startDelay"]);
const resting = {
  scale: 1, scaleX: 1, scaleY: 1, scaleZ: 1,
  x: 0, y: 0, z: 0, translateX: 0, translateY: 0, translateZ: 0, skewX: 0, skewY: 0,
  opacity: 1, pathLength: 1,
};
const flatTurns = new Set(["rotate", "rotateZ"]);
const depthTurns = new Set(["rotateX", "rotateY"]);
const scales = new Set(["scale", "scaleX", "scaleY"]);
const eps = 1e-9;

const near = (a, b) => Math.abs(a - b) < eps;
const multiple = (v, of) => near(v / of, Math.round(v / of));

const here = fileURLToPath(new URL(".", import.meta.url));
const motionURL = pathToFileURL(resolve(here, "../glyph/motion.js")).href;
const motion = await import(motionURL);
const τ = { τmax: motion.τmax, τmid: motion.τmid, τmin: motion.τmin };
const durations = [τ.τmax, τ.τmid, τ.τmin, 0.3];
const motionNames = new Set(Object.keys(motion));

function fmt(v) {
  if (typeof v !== "number") return JSON.stringify(v);
  for (const [name, t] of Object.entries(τ)) if (near(v, t)) return name;
  return String(Number(v.toFixed(4)));
}

function clocked(v) {
  for (let a = 0; a <= 4; a++) for (let b = 0; a + b <= 4; b++) if (near(a * τ.τmin + b * τ.τmid, v)) return true;
  return false;
}

const idStart = /[\p{ID_Start}$_]/u;
const idRest = /[\p{ID_Continue}$\u200c\u200d]*/uy;
const number = /(?:0[xXoObB][\da-fA-F_]+|\d[\d_]*(?:\.\d*)?|\.\d+)(?:[eE][+-]?\d+)?n?/y;
const regexBefore = new Set(["return", "typeof", "case", "in", "of", "new", "delete", "void", "throw", "yield", "await", "instanceof"]);

function lex(src) {
  const toks = [];
  const depth = [];
  let i = 0;

  const template = () => {
    let text = "";
    while (i < src.length) {
      const c = src[i];
      if (c === "\\") {
        text += src.slice(i, i + 2);
        i += 2;
      } else if (c === "`") {
        i++;
        toks.push({ t: "str", v: text });
        return;
      } else if (c === "$" && src[i + 1] === "{") {
        i += 2;
        toks.push({ t: "str", v: text });
        depth.push(0);
        return;
      } else {
        text += c;
        i++;
      }
    }
    throw new SyntaxError("unterminated template literal");
  };

  const regexAllowed = () => {
    const p = toks.at(-1);
    if (!p) return true;
    if (p.t === "p") return !")]}".includes(p.v);
    return p.t === "id" && regexBefore.has(p.v);
  };

  while (i < src.length) {
    const c = src[i];
    if (/\s/.test(c)) {
      i++;
    } else if (c === "/" && src[i + 1] === "/") {
      const end = src.indexOf("\n", i);
      i = end < 0 ? src.length : end;
    } else if (c === "/" && src[i + 1] === "*") {
      const end = src.indexOf("*/", i + 2);
      if (end < 0) throw new SyntaxError("unterminated comment");
      i = end + 2;
    } else if (c === '"' || c === "'") {
      let j = i + 1;
      let text = "";
      while (j < src.length && src[j] !== c) {
        if (src[j] === "\\") {
          text += src[j + 1];
          j += 2;
        } else text += src[j++];
      }
      if (j >= src.length) throw new SyntaxError("unterminated string");
      toks.push({ t: "str", v: text });
      i = j + 1;
    } else if (c === "`") {
      i++;
      template();
    } else if (c === "}" && depth.length && depth.at(-1) === 0) {
      depth.pop();
      i++;
      template();
    } else if (idStart.test(c)) {
      idRest.lastIndex = i + c.length;
      idRest.exec(src);
      toks.push({ t: "id", v: src.slice(i, idRest.lastIndex) });
      i = idRest.lastIndex;
    } else if (/\d/.test(c) || (c === "." && /\d/.test(src[i + 1] ?? ""))) {
      number.lastIndex = i;
      number.exec(src);
      toks.push({ t: "num", v: src.slice(i, number.lastIndex) });
      i = number.lastIndex;
    } else if (c === "/" && regexAllowed()) {
      let j = i + 1;
      let cls = false;
      while (j < src.length && (cls || src[j] !== "/")) {
        if (src[j] === "\\") j++;
        else if (src[j] === "[") cls = true;
        else if (src[j] === "]") cls = false;
        j++;
      }
      j++;
      while (/[a-z]/.test(src[j] ?? "")) j++;
      toks.push({ t: "re", v: src.slice(i, j) });
      i = j;
    } else {
      const op = src.startsWith("...", i) ? "..." : src.startsWith("?.", i) && !/\d/.test(src[i + 2] ?? "") ? "?." : c;
      if (op === "{" && depth.length) depth[depth.length - 1]++;
      if (op === "}" && depth.length) depth[depth.length - 1]--;
      toks.push({ t: "p", v: op });
      i += op.length;
    }
  }
  return toks;
}

function statics(src, report) {
  const toks = lex(src);
  const member = (k) => k > 0 && toks[k - 1].t === "p" && (toks[k - 1].v === "." || toks[k - 1].v === "?.");
  const key = (k) => toks[k + 1]?.v === ":" && (toks[k - 1]?.v === "{" || toks[k - 1]?.v === ",");
  const imported = new Set();
  const skip = new Set();

  for (let k = 0; k < toks.length; k++) {
    const tok = toks[k];
    if (tok.t !== "id" || tok.v !== "import" || member(k)) continue;
    const start = k;
    if (toks[k + 1]?.v === "(" || toks[k + 1]?.v === ".") {
      report("imports", `import${toks[k + 1].v === "(" ? "()" : ".meta"} is not allowed; import statically from "../motion.js"`);
      continue;
    }
    if (toks[k + 1]?.v !== "{") {
      report("imports", `only named imports are allowed: import { … } from "../motion.js"`);
      continue;
    }
    const names = [];
    k += 2;
    while (k < toks.length && toks[k].v !== "}") {
      if (toks[k].t === "id" && toks[k + 1]?.v !== "as" && toks[k - 1]?.v !== "as") names.push(toks[k].v);
      else if (toks[k].v === "as") report("imports", "renamed imports are not allowed");
      k++;
    }
    const from = toks[k + 1];
    const source = toks[k + 2];
    if (from?.v !== "from" || source?.t !== "str") {
      report("imports", "malformed import declaration");
      continue;
    }
    k += 2;
    if (toks[k + 1]?.v === ";") k++;
    if (source.v !== "../motion.js") report("imports", `imports from ${JSON.stringify(source.v)}; only "../motion.js" is allowed`);
    for (const name of names) {
      if (imported.has(name)) report("imports", `${name} is imported twice`);
      imported.add(name);
    }
    for (let j = start; j <= k; j++) skip.add(j);
  }

  const used = new Set();
  for (let k = 0; k < toks.length; k++) {
    const tok = toks[k];
    if (skip.has(k)) continue;
    if (tok.t === "str" && tok.v.includes("data-")) report("data", `string ${JSON.stringify(tok.v)} touches a data-* attribute`);
    if (tok.t !== "id") continue;
    if (member(k)) {
      if (continuations.has(tok.v)) report("deferred", `.${tok.v}() defers work past the hover`);
      if (tok.v === "dataset") report("data", "reads dataset");
      continue;
    }
    if (tok.v === "dataset") report("data", "reads dataset");
    if (deferrals.has(tok.v)) report("deferred", `${tok.v} defers work past the hover; use delay or stagger`);
    if (tok.v === "WeakMap") report("state", "WeakMap keeps state between hovers");
    if (key(k)) continue;
    if (outside.has(tok.v)) report("scope", `${tok.v} reaches outside el`);
    used.add(tok.v);
  }

  for (const name of imported) if (!used.has(name)) report("imports", `${name} is imported but unused`);
  for (const name of motionNames) if (used.has(name) && !imported.has(name)) report("imports", `${name} is used but not imported`);
}

const internals = new WeakMap();

class NodeList {
  constructor(nodes) {
    nodes.forEach((node, i) => (this[i] = node));
    internals.set(this, nodes);
  }
  get length() {
    return internals.get(this).length;
  }
  item(i) {
    return internals.get(this)[i] ?? null;
  }
  forEach(fn, self) {
    internals.get(this).forEach((node, i) => fn.call(self, node, i, this));
  }
  entries() {
    return internals.get(this).entries();
  }
  keys() {
    return internals.get(this).keys();
  }
  values() {
    return internals.get(this).values();
  }
  [Symbol.iterator]() {
    return internals.get(this).values();
  }
}

class Element {
  constructor(tag, attrs, parent, id) {
    internals.set(this, { tag, attrs, parent, children: [], id });
  }
  querySelector(selector) {
    return select(this, selector)[0] ?? null;
  }
  querySelectorAll(selector) {
    return new NodeList(select(this, selector));
  }
}

const node = (el) => internals.get(el);

class UnsupportedSelector extends Error {}

function parseSvg(text) {
  const tags = /<!--[\s\S]*?-->|<(\/?)([A-Za-z][\w:.-]*)((?:\s+[^\s=>/]+\s*=\s*(?:"[^"]*"|'[^']*'))*)\s*(\/?)>/g;
  const attr = /([^\s=>/]+)\s*=\s*(?:"([^"]*)"|'([^']*)')/g;
  const stack = [];
  let root = null;
  let id = 0;
  for (const m of text.matchAll(tags)) {
    if (!m[2]) continue;
    const [, close, tag, rawAttrs, self] = m;
    if (close) {
      if (node(stack.pop() ?? {})?.tag !== tag) throw new Error(`mismatched </${tag}>`);
      continue;
    }
    const attrs = new Map([...rawAttrs.matchAll(attr)].map((a) => [a[1], a[2] ?? a[3]]));
    const parent = stack.at(-1) ?? null;
    const el = new Element(tag, attrs, parent, id++);
    if (parent) node(parent).children.push(el);
    else if (root) throw new Error("more than one root element");
    else root = el;
    if (!self) stack.push(el);
  }
  if (!root || node(root).tag !== "svg") throw new Error("root element is not <svg>");
  if (stack.length) throw new Error(`unclosed <${node(stack.at(-1)).tag}>`);
  return root;
}

const selectors = new Map();

function compile(src) {
  let i = 0;
  const fail = (why = "invalid") => {
    throw new SyntaxError(`${why} selector ${JSON.stringify(src)}`);
  };
  const ws = () => {
    const from = i;
    while (/\s/.test(src[i] ?? "")) i++;
    return i > from;
  };
  const ident = () => {
    const m = /^-?[_a-zA-Z][-\w]*/.exec(src.slice(i));
    if (!m) fail();
    i += m[0].length;
    return m[0];
  };
  const value = () => {
    const q = src[i];
    if (q === '"' || q === "'") {
      const end = src.indexOf(q, i + 1);
      if (end < 0) fail();
      const v = src.slice(i + 1, end);
      i = end + 1;
      return v;
    }
    return ident();
  };
  const integer = () => {
    if (src[i] !== "(") fail();
    const m = /^\(\s*(\d+)\s*\)/.exec(src.slice(i));
    if (!m) throw new UnsupportedSelector(src);
    i += m[0].length;
    return Number(m[1]);
  };
  const siblings = (n) => (node(n).parent ? node(node(n).parent).children : [n]);
  const ofType = (n) => siblings(n).filter((s) => node(s).tag === node(n).tag);
  const ops = {
    "=": (v, w) => v === w,
    "^=": (v, w) => w !== "" && v.startsWith(w),
    "$=": (v, w) => w !== "" && v.endsWith(w),
    "*=": (v, w) => w !== "" && v.includes(w),
    "~=": (v, w) => v.split(/\s+/).includes(w),
    "|=": (v, w) => v === w || v.startsWith(`${w}-`),
  };

  function compound() {
    const tests = [];
    if (src[i] === "*") i++;
    else if (/[a-zA-Z]/.test(src[i] ?? "")) {
      const tag = ident();
      tests.push((n) => node(n).tag === tag);
    }
    for (;;) {
      const c = src[i];
      if (c === ".") {
        i++;
        const cls = ident();
        tests.push((n) => (node(n).attrs.get("class") ?? "").split(/\s+/).includes(cls));
      } else if (c === "#") {
        i++;
        const id = ident();
        tests.push((n) => node(n).attrs.get("id") === id);
      } else if (c === "[") {
        i++;
        ws();
        const name = ident();
        ws();
        if (src[i] === "]") {
          i++;
          tests.push((n) => node(n).attrs.has(name));
          continue;
        }
        const op = /^[~|^$*]?=/.exec(src.slice(i))?.[0];
        if (!op) fail();
        i += op.length;
        ws();
        const want = value();
        ws();
        if (src[i] !== "]") throw new UnsupportedSelector(src);
        i++;
        tests.push((n) => node(n).attrs.has(name) && ops[op](node(n).attrs.get(name), want));
      } else if (c === ":") {
        i++;
        const pseudo = ident();
        if (pseudo === "scope") tests.push((n, scope) => n === scope);
        else if (pseudo === "not" || pseudo === "is") {
          if (src[i] !== "(") fail();
          i++;
          const inner = list(")");
          i++;
          const any = (n, scope) => inner.some((c) => matches(c, c.length - 1, n, scope));
          tests.push(pseudo === "not" ? (n, scope) => !any(n, scope) : any);
        } else if (pseudo === "first-child") tests.push((n) => siblings(n)[0] === n);
        else if (pseudo === "last-child") tests.push((n) => siblings(n).at(-1) === n);
        else if (pseudo === "first-of-type") tests.push((n) => ofType(n)[0] === n);
        else if (pseudo === "last-of-type") tests.push((n) => ofType(n).at(-1) === n);
        else if (pseudo === "nth-child") {
          const k = integer();
          tests.push((n) => siblings(n)[k - 1] === n);
        } else if (pseudo === "nth-of-type") {
          const k = integer();
          tests.push((n) => ofType(n)[k - 1] === n);
        } else throw new UnsupportedSelector(src);
      } else break;
    }
    if (!tests.length && src[i - 1] !== "*") fail();
    return tests;
  }

  function complex() {
    ws();
    const parts = [{ comb: null, tests: compound() }];
    for (;;) {
      const spaced = ws();
      let comb;
      if (src[i] === ">") {
        i++;
        ws();
        comb = ">";
      } else if (src[i] === "+" || src[i] === "~") throw new UnsupportedSelector(src);
      else if (spaced && i < src.length && src[i] !== "," && src[i] !== ")") comb = " ";
      else break;
      parts.push({ comb, tests: compound() });
    }
    return parts;
  }

  function list(end) {
    const out = [complex()];
    while (src[i] === ",") {
      i++;
      out.push(complex());
    }
    if (end ? src[i] !== end : i < src.length) fail();
    return out;
  }

  return list();
}

function matches(parts, k, n, scope) {
  if (!parts[k].tests.every((test) => test(n, scope))) return false;
  if (k === 0) return true;
  const up = node(n).parent;
  if (parts[k].comb === ">") return up !== null && matches(parts, k - 1, up, scope);
  for (let p = up; p; p = node(p).parent) if (matches(parts, k - 1, p, scope)) return true;
  return false;
}

function select(scope, selector) {
  selector = String(selector);
  if (!selectors.has(selector)) selectors.set(selector, compile(selector));
  const list = selectors.get(selector);
  const found = [];
  const walk = (n) => {
    for (const child of node(n).children) {
      if (list.some((parts) => matches(parts, parts.length - 1, child, scope))) found.push(child);
      walk(child);
    }
  };
  walk(scope);
  return found;
}

class Stagger {
  constructor(step, options) {
    this.step = step;
    this.options = options;
  }
}

const plain = (v) => v !== null && typeof v === "object" && Object.getPrototypeOf(v) === Object.prototype;

function targets(t) {
  if (t instanceof Element) return [t];
  if (t instanceof NodeList) return [...t];
  if (Array.isArray(t)) return t;
  return null;
}

function forbidden(obj, where, report) {
  for (const [k, v] of Object.entries(obj)) {
    if (easings.has(k)) report("easing", `${where} sets ${k}`);
    if (plain(v)) forbidden(v, `${where}.${k}`, report);
  }
}

function analyze(name, calls, report) {
  if (calls.length > 4) report("calls", `${calls.length} animate() calls`);
  const beats = new Set();
  const delays = [];

  calls.forEach((args, c) => {
    const at = `call ${c + 1}`;
    const [target, keyframes, opts] = args;
    if (args.length !== 3 || !plain(keyframes) || !plain(opts)) {
      report("call", `${at}: expected animate(target, keyframes, options)`);
      return;
    }

    const parts = targets(target);
    if (typeof target === "string") report("scope", `${at}: target is the selector string ${JSON.stringify(target)}; select inside el`);
    else if (!parts) report("target", `${at}: target is ${target === null ? "null; its selector matches nothing" : typeof target}`);
    else if (!parts.length) report("target", `${at}: target list is empty; its selector matches nothing`);
    else {
      const bad = parts.flatMap((p, i) => (p instanceof Element ? [] : [i]));
      if (bad.length) report("target", `${at}: target items ${bad.join(", ")} are not parts (a selector matches nothing)`);
    }
    const circles = parts?.length > 0 && parts.every((p) => p instanceof Element && node(p).tag === "circle" && !node(p).attrs.has("transform-origin"));

    forbidden(keyframes, `${at} keyframes`, report);
    forbidden(opts, `${at} options`, report);

    for (const [prop, frames] of Object.entries(keyframes)) {
      if (easings.has(prop)) continue;
      if (paint.test(prop)) {
        report("paint", `${at}: animates ${prop}`);
        continue;
      }
      if (prop === "pathOffset") {
        report("direction", `${at}: animates pathOffset`);
        continue;
      }
      const allowedNoop = prop === "opacity" && Array.isArray(frames) && frames.length === 1 && frames[0] === 1;
      if (!allowedNoop && (!Array.isArray(frames) || frames.length < 2)) {
        report("keyframes", `${at}: ${prop} is ${JSON.stringify(frames)}; use an array of at least 2 keyframes`);
        continue;
      }
      if (!frames.every((f) => typeof f === "number" && Number.isFinite(f))) {
        report("keyframes", `${at}: ${prop} has non-numeric keyframes ${JSON.stringify(frames)}`);
        continue;
      }
      const end = frames.at(-1);
      if (scales.has(prop)) {
        const peak = Math.max(...frames);
        if (peak > 1.75 + eps) report("overshoot", `${at}: ${prop} overshoots to ${fmt(peak)}`);
      }
      if (prop in resting) {
        if (!near(end, resting[prop])) report("rest", `${at}: ${prop} ends at ${fmt(end)}, not ${resting[prop]}`);
      } else if (flatTurns.has(prop)) {
        if (circles) continue;
        if (!multiple(end, 90)) report("turn", `${at}: ${prop} ends at ${fmt(end)}°, not a multiple of 90°`);
        else if (!multiple(end, 180)) report("symmetry", `${at}: ${prop} ends at ${fmt(end)}°; the part needs four-fold symmetry`);
      } else if (depthTurns.has(prop)) {
        if (!multiple(end, 180)) report("turn", `${at}: ${prop} ends at ${fmt(end)}°, not a multiple of 180°; the part stays foreshortened`);
      } else report("property", `${at}: ${prop} has no rest rule`);
    }

    for (const k of Object.keys(opts)) if (!options.has(k) && !easings.has(k)) report("option", `${at}: option ${k}`);

    const { duration, delay = 0 } = opts;
    if (duration === undefined) report("duration", `${at}: no duration`);
    else if (!durations.some((d) => near(d, duration))) report("duration", `${at}: duration ${fmt(duration)}`);

    let start = 0;
    if (delay instanceof Stagger) {
      const { step, options: so = {} } = delay;
      if (typeof step !== "number" || step < 0.03 - eps || step > 0.21 + eps) report("stagger", `${at}: stagger step ${fmt(step)}`);
      if (!plain(so)) report("stagger", `${at}: stagger options ${JSON.stringify(so)}`);
      else {
        forbidden(so, `${at} stagger`, report);
        for (const k of Object.keys(so)) if (!staggerOptions.has(k) && !easings.has(k)) report("option", `${at}: stagger option ${k}`);
        start = so.startDelay ?? 0;
        if (typeof start !== "number" || !clocked(start)) report("delay", `${at}: stagger startDelay ${fmt(start)}`);
      }
      if (parts?.length && typeof step === "number" && typeof duration === "number") {
        const last = start + step * (parts.length - 1) + duration;
        if (last > 0.6 + eps) report("cascade", `${at}: item ${parts.length} of the cascade ends at ${fmt(last)} s`);
      }
    } else if (typeof delay !== "number") report("call", `${at}: delay is ${typeof delay}`);
    else {
      start = delay;
      delays.push({ at: c + 1, delay });
    }
    if (typeof start === "number") beats.add(Math.round(start * 1e6));
  });

  const off = delays.filter((d) => !clocked(d.delay));
  const values = [...new Set(delays.map((d) => Math.round(d.delay * 1e6)))].sort((a, b) => a - b).map((v) => v / 1e6);
  const step = values[1] - values[0];
  const run = off.length && values.length >= 3 && values.every((v, i) => near(v, values[0] + i * step)) && step > 0.03 - eps && step < 0.21 + eps;
  if (run) {
    const from = values[0] ? `, { startDelay: ${fmt(values[0])} }` : "";
    report("delay", `calls ${delays.map((d) => d.at).join(", ")}: delays ${values.map(fmt).join(", ")} step by ${fmt(step)}; use stagger(${fmt(step)}${from}) on shared keyframes`);
  } else for (const d of off) report("delay", `call ${d.at}: delay ${fmt(d.delay)} is off the 0, τmin, τmid clock`);

  if (beats.size > 3) report("beats", `${beats.size} beats at ${[...beats].sort((a, b) => a - b).map((b) => fmt(b / 1e6)).join(", ")}`);
}

function serialize(calls) {
  return JSON.stringify(calls, (_, v) => {
    if (v instanceof Element) return `#${node(v).id}`;
    if (v instanceof NodeList) return [...v];
    if (v instanceof Stagger) return { stagger: v.step, ...v.options };
    if (typeof v === "function") return `fn:${v.name}`;
    return v;
  });
}

const { values: args, positionals } = parseArgs({
  allowPositionals: true,
  options: {
    only: { type: "string", default: "" },
    svg: { type: "string", default: "" },
    rules: { type: "boolean", default: false },
  },
});
if (args.rules) {
  for (const [rule, [level, ref]] of Object.entries(rules)) console.log(`${rule}: ${level}: ${ref}`);
  process.exit(0);
}
if (positionals.length > 1) {
  console.error("usage: node icons/check/lint.mjs [--only a,b] [--svg dir] [dir]");
  process.exit(2);
}
const dir = resolve(positionals[0] ?? resolve(here, "../anim"));
const svgDir = resolve(args.svg || join(dir, "../svg"));
const present = readdirSync(dir).filter((f) => f.endsWith(".js")).map((f) => basename(f, ".js")).sort();
const wanted = args.only ? args.only.split(",").map((s) => s.trim()).filter(Boolean) : present;

const findings = [];
const began = performance.now();
const scratch = mkdtempSync(join(tmpdir(), "hgmx-icons-lint-"));
const hook = Symbol.for("hgmx.icons.lint");
const recorder = { calls: null };
globalThis[hook] = recorder;

try {
  mkdirSync(join(scratch, "anim"));
  writeFileSync(
    join(scratch, "motion.js"),
    [
      `export * from ${JSON.stringify(motionURL)};`,
      `const lint = globalThis[Symbol.for("hgmx.icons.lint")];`,
      `export function animate(...args) { if (!lint.calls) throw new Error("animate() called outside the hover"); lint.calls.push(args); return Object.freeze({}); }`,
      `export function stagger(...args) { return lint.stagger(...args); }`,
      "",
    ].join("\n"),
  );
  recorder.stagger = (step, options) => new Stagger(step, options);

  const modules = wanted.map((name) => {
    const report = (rule, detail) => {
      if (!exempt[name]?.has(rule)) findings.push({ name, rule, detail });
    };
    const file = join(dir, `${name}.js`);
    if (!existsSync(file)) {
      report("export", `no module ${file}`);
      return null;
    }
    const src = readFileSync(file, "utf8");
    try {
      statics(src, report);
    } catch (e) {
      report("syntax", e.message);
    }
    const copy = join(scratch, "anim", `${name}.js`);
    writeFileSync(copy, src);
    return { name, report, url: pathToFileURL(copy).href };
  });

  const loaded = await Promise.all(modules.map((m) => m && import(m.url).then((ns) => ({ ns }), (error) => ({ error }))));

  modules.forEach((m, i) => {
    if (!m) return;
    const { name, report } = m;
    const { ns, error } = loaded[i];
    if (error) {
      if (error.code === "ERR_MODULE_NOT_FOUND") return;
      const msg = String(error.message).split("\n")[0];
      report(/does not provide an export named/.test(msg) ? "imports" : "syntax", msg);
      return;
    }
    const extra = Object.keys(ns).filter((k) => k !== "default");
    if (extra.length) report("export", `exports ${extra.join(", ")} besides default`);
    if (typeof ns.default !== "function") {
      report("export", "no default function export");
      return;
    }

    const svgFile = join(svgDir, `${name}.svg`);
    if (!existsSync(svgFile)) {
      report("svg", `no ${svgFile}`);
      return;
    }
    let root;
    try {
      root = parseSvg(readFileSync(svgFile, "utf8"));
    } catch (e) {
      report("svg", `${name}.svg: ${e.message}`);
      return;
    }

    const runs = [];
    for (let r = 0; r < 2; r++) {
      recorder.calls = [];
      try {
        const out = ns.default(root);
        if (r === 0 && out !== undefined) report("deferred", `default function returns ${out?.then ? "a promise" : typeof out}`);
        runs.push(recorder.calls);
      } catch (e) {
        report("run", e instanceof UnsupportedSelector ? `the lint cannot evaluate selector ${JSON.stringify(e.message)}` : `${e.name}: ${e.message}`);
        return;
      } finally {
        recorder.calls = null;
      }
    }
    if (!runs[0].length) report("call", "no animate() calls");
    if (serialize(runs[0]) !== serialize(runs[1])) report("state", "a second hover records different animations");
    analyze(name, runs[0], report);
  });
} finally {
  delete globalThis[hook];
  rmSync(scratch, { recursive: true, force: true });
}

let errors = 0;
let warnings = 0;
findings.sort((a, b) => (a.name < b.name ? -1 : a.name > b.name ? 1 : 0));
for (const { name, rule, detail } of findings) {
  const warning = rules[rule][0] === "warning";
  if (warning) warnings++;
  else errors++;
  console.log(`${name}: ${rule}: ${warning ? "warning: " : ""}${detail}`);
}
const seconds = ((performance.now() - began) / 1000).toFixed(2);
console.error(`${wanted.length} modules, ${errors} errors, ${warnings} warnings, ${seconds} s`);
process.exit(errors ? 1 : 0);
