import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const greens = el.querySelector("g");
  animate(greens, { scale: [1, 0.85, 1.1, 1] }, { duration: τmax });

  const vein = el.querySelector('path[d="M2 22 17 7"]');
  animate(vein, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
