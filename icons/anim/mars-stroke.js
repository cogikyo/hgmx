import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const arm = el.querySelector("g");
  animate(arm, { scale: [1, 1.25, 0.95, 1] }, { duration: τmax });
  const stroke = el.querySelector('path[d^="m14 6"]');
  animate(stroke, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
