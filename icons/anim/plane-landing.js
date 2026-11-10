import { animate, τmax } from "../motion.js";

export default function (el) {
  const plane = el.querySelector('path[d^="M3.77"]');
  animate(plane, { translateX: [0, 2, -0.5, 0], translateY: [0, 2, -0.5, 0], rotate: [0, 6, -2, 0] }, { duration: τmax });
}
