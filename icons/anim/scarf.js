import { animate, τmax } from "../motion.js";

export default function (el) {
  const tail = el.querySelectorAll('path[d^="M13 20"], path[d^="M17 12.1"], path[d^="M17 18"], path[d^="M17 20"], path[d^="M21 8.2"]');
  animate(tail, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
}
