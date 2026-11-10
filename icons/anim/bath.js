import { animate, τmax } from "../motion.js";

export default function (el) {
  const tub = el.querySelector('path[d^="M9 5"]');
  animate(tub, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
}
