import { animate } from "../motion.js";

export default function (el) {
  const head = el.querySelector('path[d^="M6 8"]');
  animate(head, { translateX: [0, -3, 1, 0] }, { duration: 0.3 });
}
