import { animate, τmax } from "../motion.js";

export default function (el) {
  const left = el.querySelector('path[d^="M3 3h3"]');
  animate(left, { translateX: [0, -2, 0.5, 0] }, { duration: τmax });

  const right = el.querySelector('path[d^="M21 3h-3"]');
  animate(right, { translateX: [0, 2, -0.5, 0] }, { duration: τmax });
}
