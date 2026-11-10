import { animate, τmax } from "../motion.js";

export default function (el) {
  const head = el.querySelector('path[d^="M16.93"]');
  animate(head, { scaleX: [1, 1.1, 0.97, 1] }, { duration: τmax });
}
