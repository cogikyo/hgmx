import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glyph = el.querySelector('path[d^="m8 16"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const bar = el.querySelector('path[d^="M15 14"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
  const left = el.querySelector('path[d="M3 3v18"]');
  animate(left, { translateX: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
  const right = el.querySelector('path[d="M21 3v18"]');
  animate(right, { translateX: [0, 2, -0.5, 0] }, { duration: τmid, delay: τmin });
}
