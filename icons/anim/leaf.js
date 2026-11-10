import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const leaf = el.querySelector("g");
  animate(leaf, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });

  const vein = el.querySelector('path[d^="M2 21"]');
  animate(vein, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
