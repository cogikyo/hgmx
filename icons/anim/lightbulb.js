import { animate, τmax } from "../motion.js";

export default function (el) {
  const bulb = el.querySelector('path[d^="M15 14"]');
  animate(bulb, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
