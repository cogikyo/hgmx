import { animate, τmax } from "../motion.js";

export default function (el) {
  const arc = el.querySelector('path[d^="M3 11"]');
  animate(arc, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
