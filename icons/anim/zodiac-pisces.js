import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const fish = el.querySelectorAll('path[d^="M19 21"], path[d^="M5 3"]');
  animate(fish, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const cord = el.querySelector('path[d="M20 12H4"]');
  animate(cord, { scaleX: [0, 1.15, 1] }, { duration: τmid, delay: τmin });
}
