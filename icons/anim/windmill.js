import { animate, τmax } from "../motion.js";

export default function (el) {
  const blades = el.querySelectorAll('path[d^="m10 14"], path[d^="m8 8"], path[d^="m16 12"]');
  animate(blades, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
}
