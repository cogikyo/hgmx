import { animate, τmax } from "../motion.js";

export default function (el) {
  const wings = el.querySelector('path[d^="M20 5"]');
  animate(wings, { scaleX: [1, 0.3, 1, 0.5, 1] }, { duration: τmax });
}
