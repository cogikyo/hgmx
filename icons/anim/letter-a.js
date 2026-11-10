import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glyph = el.querySelector('path[d^="m4 20"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const bar = el.querySelector('path[d^="M17.117"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
