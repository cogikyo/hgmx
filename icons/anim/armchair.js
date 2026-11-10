import { animate, τmax } from "../motion.js";

export default function (el) {
  const back = el.querySelector('path[d^="M19 9"]');
  animate(back, { scaleY: [1, 0.8, 1.1, 1] }, { duration: τmax });
}
