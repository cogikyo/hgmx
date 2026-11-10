import { animate, τmax } from "../motion.js";

export default function (el) {
  const left = el.querySelector('path[d^="M7.5 14"]');
  animate(left, { rotate: [0, -15, 3, 0] }, { duration: τmax });
  const right = el.querySelector('path[d^="M16.5 14"]');
  animate(right, { rotate: [0, 15, -3, 0] }, { duration: τmax });
}
