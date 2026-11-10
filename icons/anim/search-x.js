import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const lens = el.querySelector("g");
  animate(lens, { scale: [1.35, 1] }, { duration: τmax });
  animate(lens, { translateX: [0, -4, 2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });

  const bars = el.querySelectorAll('path[d^="m13.5"], path[d^="m8.5"]');
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
