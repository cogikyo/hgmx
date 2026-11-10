import { animate, τmax } from "../motion.js";

export default function (el) {
  const right = el.querySelector('path[d^="m16 18"]');
  animate(right, { translateX: [-5, 1, 0], opacity: [0, 1] }, { duration: τmax });

  const left = el.querySelector('path[d^="m8 6"]');
  animate(left, { translateX: [5, -1, 0], opacity: [0, 1] }, { duration: τmax });
}
