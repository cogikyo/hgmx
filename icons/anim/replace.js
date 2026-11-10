import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelectorAll('path[d^="m3 7"], path[d^="M6 10V5"]');
  animate(arrow, { translateY: [0, 2, -0.5, 0] }, { duration: τmid });
  const target = el.querySelector("rect");
  animate(target, { scale: [1, 0.85, 1.1, 1] }, { duration: τmid, delay: τmin });
}
