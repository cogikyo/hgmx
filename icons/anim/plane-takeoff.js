import { animate, τmax } from "../motion.js";

export default function (el) {
  const plane = el.querySelector('path[d^="M6.36"]');
  animate(plane, { translateX: [0, 2, -0.5, 0], translateY: [0, -3, 1, 0], rotate: [0, -8, 3, 0] }, { duration: τmax });
}
