import { animate, τmax } from "../motion.js";

export default function (el) {
  const lower = el.querySelector('path[d^="M6 2"]');
  animate(lower, { translateX: [0, 1.5, -0.5, 0], translateY: [0, -1.5, 0.5, 0] }, { duration: τmax });

  const upper = el.querySelector('path[d^="M18 22"]');
  animate(upper, { translateX: [0, -1.5, 0.5, 0], translateY: [0, 1.5, -0.5, 0] }, { duration: τmax });
}
