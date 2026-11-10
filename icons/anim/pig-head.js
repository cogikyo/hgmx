import { animate, τmax } from "../motion.js";

export default function (el) {
  const snout = el.querySelectorAll('path[d^="M14 22"], path[d$="18h.01"]');
  animate(snout, { translateX: [0, -1.5, 1.5, 0] }, { duration: τmax });
}
