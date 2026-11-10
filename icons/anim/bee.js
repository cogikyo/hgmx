import { animate, τmax } from "../motion.js";

export default function (el) {
  const bee = el.querySelector("g");
  animate(bee, { translateY: [0, -1.5, 0.5, 0] }, { duration: τmax });

  const wings = el.querySelector('path[d^="M5 7"]');
  animate(wings, { scaleY: [1, 1.2, 1, 1.2, 1] }, { duration: τmax });
}
