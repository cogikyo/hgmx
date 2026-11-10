import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pump = el.querySelectorAll('path[d^="M14 21"], path[d^="M2 21"], path[d^="M3 9"]');
  animate(pump, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const hose = el.querySelector('path[d^="M14 13"]');
  animate(hose, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
