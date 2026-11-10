import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glyph = el.querySelector('path[d^="M18 4"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const bar = el.querySelector('path[d="M6 12h8"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
