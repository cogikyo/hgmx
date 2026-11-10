import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const glass = el.querySelector("g");
  animate(glass, { rotate: [0, -10, 6, 0] }, { duration: τmax });

  const wine = el.querySelector('path[d^="M7 10"]');
  animate(wine, { rotate: [0, 12, -8, 0] }, { duration: τmid, delay: τmin });
}
