import { animate, τmax } from "../motion.js";

export default function (el) {
  const swirl = el.querySelector('path[d^="M11 11"]');
  animate(swirl, { rotate: [0, 360], scale: [1, 0.85, 1] }, { duration: τmax });
}
