import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("g");
  animate(front, { translateX: [-6, 1, 0], translateY: [-6, 1, 0] }, { duration: τmax });

  const back = el.querySelector('path[d^="M4 16"]');
  animate(back, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const lines = el.querySelectorAll('path[d$="h6"]');
  animate(lines, { opacity: [0, 1], x: [-3, 0] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
