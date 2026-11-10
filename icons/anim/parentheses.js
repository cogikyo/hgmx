import { animate, τmax } from "../motion.js";

export default function (el) {
  const left = el.querySelector('path[d^="M8 21"]');
  animate(left, { translateX: [0, 2, -0.5, 0] }, { duration: τmax });

  const right = el.querySelector('path[d^="M16 3"]');
  animate(right, { translateX: [0, -2, 0.5, 0] }, { duration: τmax });
}
