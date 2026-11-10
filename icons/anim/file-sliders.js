import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M6 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const upper = el.querySelector('path[d="M10 11v2"]');
  animate(upper, { translateX: [0, 3, -1, 0] }, { duration: τmid, delay: τmin });

  const lower = el.querySelector('path[d="M14 16v2"]');
  animate(lower, { translateX: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
