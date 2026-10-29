import { animate, τmax } from "../motion.js";

export default function (el) {
  const brackets = el.querySelectorAll("path");
  animate(brackets, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const left = el.querySelector('path[d^="M8 21"]');
  animate(left, { translateX: [-4, 1, 0] }, { duration: τmax });

  const right = el.querySelector('path[d^="M16 3"]');
  animate(right, { translateX: [4, -1, 0] }, { duration: τmax });
}
