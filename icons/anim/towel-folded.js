import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const towel = el.querySelector("g");
  animate(towel, { scaleY: [1, 0.9, 1.05, 1] }, { duration: τmid });

  const fold = el.querySelector('path[d^="M17 17H7"]');
  animate(fold, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
