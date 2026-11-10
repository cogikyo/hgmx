import { animate, τmax } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d^="M3 3v3"]');
  animate(top, { translateY: [0, -2, 0.5, 0] }, { duration: τmax });

  const bottom = el.querySelector('path[d^="M3 21v-3"]');
  animate(bottom, { translateY: [0, 2, -0.5, 0] }, { duration: τmax });
}
