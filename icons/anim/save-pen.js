import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const disk = el.querySelectorAll('path[d^="M9 21"], path[d^="M13.33"]');
  animate(disk, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pen = el.querySelector('path[d^="M14.363"]');
  animate(pen, { rotate: [-15, 15, -10, 0], translateX: [0, -3, 0] }, { duration: τmax });
}
