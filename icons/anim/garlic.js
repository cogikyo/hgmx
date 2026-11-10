import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bulb = el.querySelector('path[d^="M15 4"]');
  animate(bulb, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cloves = el.querySelector('path[d^="M9 5"]');
  animate(cloves, { scaleX: [0, 1.1, 1] }, { duration: τmid, delay: τmin });
}
