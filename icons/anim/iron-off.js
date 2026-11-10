import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const iron = el.querySelectorAll('path:not([d="M3 22h10"]):not([d="m2 2 20 20"])');
  animate(iron, { translateX: [0, -1.5, 1, 0] }, { duration: τmid });

  const slash = el.querySelector('path[d="m2 2 20 20"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
