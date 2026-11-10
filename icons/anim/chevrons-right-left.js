import { animate, τmax } from "../motion.js";

export default function (el) {
  const left = el.querySelector('path[d^="m4 17"]');
  animate(left, { translateX: [0, 1.5, -0.5, 0] }, { duration: τmax });

  const right = el.querySelector('path[d^="m20 17"]');
  animate(right, { translateX: [0, -1.5, 0.5, 0] }, { duration: τmax });
}
