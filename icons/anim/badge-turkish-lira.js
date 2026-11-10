import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const badge = el.querySelector('path[d^="M3.85"]');
  animate(badge, { scale: [1, 0.85, 1.1, 1] }, { duration: τmax });

  const glyph = el.querySelector('path[d^="M11 7"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const bar = el.querySelector('path[d="m15 8-6 3"]');
  animate(bar, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
