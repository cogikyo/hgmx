import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("rect");
  animate(front, { translateX: [-6, 1, 0], translateY: [-6, 1, 0] }, { duration: τmax });

  const behind = [el.querySelector('path[d^="M10 16"]'), el.querySelector('path[d^="M4 10"]')];
  animate(behind, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
