import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelectorAll('path[d^="M15 2"], path[d^="M16.706"]');
  animate(front, { translateX: [-4, 1, 0], translateY: [4, -1, 0] }, { duration: τmax });

  const back = el.querySelector('path[d^="M5 7"]');
  animate(back, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid, delay: τmin });
}
