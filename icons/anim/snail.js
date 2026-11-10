import { animate, τmax } from "../motion.js";

export default function (el) {
  const spiral = el.querySelector('path[d^="M2 13a6"]');
  animate(spiral, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
  const stalks = el.querySelectorAll('path[d$="5.2"]');
  animate(stalks, { scale: [1, 0, 1.25, 1] }, { duration: τmax });
}
