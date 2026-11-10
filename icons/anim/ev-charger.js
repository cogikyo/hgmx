import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pump = el.querySelectorAll('path[d^="M14 21"], path[d^="M2 21"], path[d^="M3 7"]');
  animate(pump, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bolt = el.querySelector('path[d^="m9 11"]');
  animate(bolt, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
