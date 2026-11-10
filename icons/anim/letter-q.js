import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glyph = el.querySelector('circle');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1], rotate: [270, 90] }, { duration: τmid });
  const tail = el.querySelector('path[d^="M20 20"]');
  animate(tail, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
