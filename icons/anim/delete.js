import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const key = el.querySelector('path[d^="M10 5"]');
  animate(key, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const x = el.querySelectorAll('path[d="m12 9 6 6"], path[d="m18 9-6 6"]');
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
