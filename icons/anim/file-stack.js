import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector('path[d^="M21 6"]');
  animate(front, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const behind = el.querySelectorAll('path[d^="M16 16"], path[d^="M11 21"]');
  animate(behind, { translateX: [5, -1, 0], translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
