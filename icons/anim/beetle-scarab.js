import { animate, τmax } from "../motion.js";

export default function (el) {
  const wings = el.querySelector('path[d^="M5 4.8"]');
  animate(wings, { scaleY: [1, 1.35, 0.9, 1] }, { duration: τmax });
}
