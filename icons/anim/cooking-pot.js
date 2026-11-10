import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const pot = el.querySelectorAll('path[d^="M2 12"], path[d^="M20 12"]');
  animate(pot, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lid = el.querySelector("g");
  animate(lid, { translateY: [0, -2, 0, -1, 0], rotate: [0, -6, 0, -3, 0] }, { duration: τmax });
}
