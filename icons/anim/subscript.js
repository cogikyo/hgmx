import { animate, τmax } from "../motion.js";

export default function (el) {
  const two = el.querySelector('path[d^="M20 19"]');
  animate(two, { translateY: [0, 3, -1, 0] }, { duration: τmax });
}
