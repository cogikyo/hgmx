import { animate, τmax } from "../motion.js";

export default function (el) {
  const upper = el.querySelector('path[d^="M10 13"]');
  animate(upper, { translateX: [3, -1, 0], translateY: [-3, 1, 0] }, { duration: τmax });

  const lower = el.querySelector('path[d^="M14 11"]');
  animate(lower, { translateX: [-3, 1, 0], translateY: [3, -1, 0] }, { duration: τmax });
}
