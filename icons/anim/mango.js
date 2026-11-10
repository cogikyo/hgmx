import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const mango = el.querySelector('path[d^="M7 22"]');
  animate(mango, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stem = el.querySelector('path[d^="m17 6"]');
  animate(stem, { rotate: [0, -20, 10, 0] }, { duration: τmid, delay: τmin });
}
