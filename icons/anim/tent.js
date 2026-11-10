import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const poles = el.querySelectorAll('path[d="M3.5 21 14 3"], path[d="M20.5 21 10 3"]');
  animate(poles, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const door = el.querySelector('path[d^="M15.5 21"]');
  animate(door, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
