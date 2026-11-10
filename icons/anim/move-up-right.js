import { animate } from "../motion.js";

export default function (el) {
  const head = el.querySelector('path[d^="M13 5"]');
  animate(head, { translateX: [0, 2, -0.7, 0], translateY: [0, -2, 0.7, 0] }, { duration: 0.3 });
}
