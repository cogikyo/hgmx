import { animate, τmax } from "../motion.js";

export default function (el) {
  const chevron = el.querySelector('path[d^="m7 18"]');
  animate(chevron, { translateX: [0, 2.5, -1, 0] }, { duration: τmax });

  const bar = el.querySelector('path[d="M17 6v12"]');
  animate(bar, { translateX: [0, 1, -0.4, 0] }, { duration: τmax });
}
