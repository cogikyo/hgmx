import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelector('path[d^="M4 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const glyph = el.querySelector('path[d^="m8 13"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const bar = el.querySelector('path[d^="M9.1 11"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
