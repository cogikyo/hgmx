import { animate } from "../motion.js";

export default function (el) {
  const head = el.querySelector('path[d^="M8 6"]');
  animate(head, { translateY: [0, -3, 1, 0] }, { duration: 0.3 });
}
