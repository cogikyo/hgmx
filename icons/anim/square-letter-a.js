import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const glyph = el.querySelector('path[d^="m8 16"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const bar = el.querySelector('path[d^="M15 14"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
