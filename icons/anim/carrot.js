import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const carrot = el.querySelector('path[d^="M15 16"]');
  animate(carrot, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const leaves = el.querySelector('path[d^="M15 9q"]');
  animate(leaves, { rotate: [0, -20, 8, 0] }, { duration: τmid, delay: τmin });
}
