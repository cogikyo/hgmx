import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("g");
  animate(front, { translateX: [-6, 1, 0], translateY: [-6, 1, 0] }, { duration: τmax });

  const back = el.querySelector('path[d^="M4 16"]');
  animate(back, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const bars = el.querySelectorAll('path[d^="m12.5"]');
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmin, delay: stagger(τmin, { startDelay: τmin }) });
}
