import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glyph = el.querySelector('path[d^="M19 18"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const leg = el.querySelector('path[d^="M11.38"]');
  animate(leg, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
