import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const mug = el.querySelectorAll('path[d^="M5 8"], path[d^="M17 11"]');
  animate(mug, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const foam = el.querySelector('path[d^="M14 7.5"]');
  animate(foam, { scale: [1, 0.85, 1.15, 1] }, { duration: τmid, delay: τmin });
}
