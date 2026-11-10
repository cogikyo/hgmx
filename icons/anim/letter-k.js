import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glyph = el.querySelector('path[d="M6 20V4"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const arms = el.querySelectorAll('path[d^="M17 4"], path[d^="m9 12"]');
  animate(arms, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
