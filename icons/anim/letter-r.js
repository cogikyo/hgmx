import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glyph = el.querySelector('path[d^="M6 12"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const leg = el.querySelector('path[d^="m10 12"]');
  animate(leg, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
