import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("g");
  animate(front, { translateX: [-6, 1, 0], translateY: [-6, 1, 0] }, { duration: τmax });

  const back = el.querySelector('path[d^="M4 16"]');
  animate(back, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const arrow = el.querySelectorAll('path[d="M15 12v6"], path[d^="m12 15"]');
  animate(arrow, { translateY: [-2, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
