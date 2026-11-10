import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelectorAll('path[d^="M2 10"], path[d^="m2 4"]');
  animate(arrow, { translateX: [0, 1.5, -0.5, 0], translateY: [0, 1.5, -0.5, 0] }, { duration: τmax });

  const inset = el.querySelector("rect");
  animate(inset, { scale: [1, 1.15, 1] }, { duration: τmin, delay: τmid });
}
