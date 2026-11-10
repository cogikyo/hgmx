import { animate, τmax } from "../motion.js";

export default function (el) {
  const arrow = el.querySelectorAll('path[d^="m13 21"], path[d^="M20 18"]');
  animate(arrow, { translateX: [0, -4, 1, 0] }, { duration: τmax });
}
