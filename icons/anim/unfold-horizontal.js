import { animate, τmax } from "../motion.js";

export default function (el) {
  const left = el.querySelectorAll('path[d^="M8 12"], path[d^="m5 9"]');
  animate(left, { translateX: [4, -1, 0], opacity: [0, 1] }, { duration: τmax });

  const right = el.querySelectorAll('path[d^="M16 12"], path[d^="m19 15"]');
  animate(right, { translateX: [-4, 1, 0], opacity: [0, 1] }, { duration: τmax });
}
