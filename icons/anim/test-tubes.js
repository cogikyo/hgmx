import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const tubes = el.querySelectorAll("g");
  animate(tubes, { rotate: [0, -10, 8, -4, 0] }, { duration: τmax });

  const liquids = el.querySelectorAll('path[d^="M9 16"], path[d^="M20 16"]');
  animate(liquids, { rotate: [0, 12, -8, 0] }, { duration: τmid, delay: τmin });
}
