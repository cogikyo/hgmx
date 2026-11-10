import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const cord = el.querySelector('path[d^="M19 14"]');
  animate(cord, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const prongs = el.querySelectorAll('path[d="M17 21v-2"], path[d="M21 21v-2"], path[d="M3 5V3"], path[d="M7 5V3"]');
  animate(prongs, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
