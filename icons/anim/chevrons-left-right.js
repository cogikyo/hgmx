import { animate, τmax } from "../motion.js";

export default function (el) {
  const left = el.querySelector('path[d^="m9 7"]');
  animate(left, { translateX: [0, -3, 1, 0] }, { duration: τmax });

  const right = el.querySelector('path[d^="m15 7"]');
  animate(right, { translateX: [0, 3, -1, 0] }, { duration: τmax });
}
