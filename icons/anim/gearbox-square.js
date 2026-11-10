import { animate, τmax } from "../motion.js";

export default function (el) {
  const gate = el.querySelector('path[d^="M17 7"]');
  animate(gate, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
