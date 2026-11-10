import { animate, τmax } from "../motion.js";

export default function (el) {
  const stalks = el.querySelectorAll('path[d$="5.2"]');
  animate(stalks, { scale: [1, 0, 1.25, 1] }, { duration: τmax });
}
