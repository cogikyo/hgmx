import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M3 12.34"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const arrow = el.querySelectorAll('path[d^="M3 18h6"], path[d^="m6 21-3-3"]');
  animate(arrow, { translateX: [3, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
