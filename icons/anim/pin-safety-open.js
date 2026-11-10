import { animate, τmax } from "../motion.js";

export default function (el) {
  const arm = el.querySelector('path[d^="M9 2s"]');
  animate(arm, { rotate: [0, -12, 6, -2, 0] }, { duration: τmax });
}
