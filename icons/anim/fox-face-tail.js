import { animate, τmax } from "../motion.js";

export default function (el) {
  const tail = el.querySelector('path[d^="M6.3 20.5"]');
  animate(tail, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
