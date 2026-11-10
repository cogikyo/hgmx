import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glyph = el.querySelector('path[d^="m6 4"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const stem = el.querySelector('path[d^="M12 12"]');
  animate(stem, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
