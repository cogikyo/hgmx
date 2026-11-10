import { animate, τmax } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M10.3"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
  const feather = el.querySelector('path[d^="M7 12"]');
  animate(feather, { rotate: [0, -12, 8, 0] }, { duration: τmax });
}
