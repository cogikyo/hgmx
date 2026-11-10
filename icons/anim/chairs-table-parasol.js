import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const chairs = el.querySelectorAll('path[d^="M8 21"], path[d^="M3 12"], path[d^="M16 21"], path[d^="M21 12"]');
  animate(chairs, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const parasol = el.querySelector('path[d^="M21 8H3"]');
  animate(parasol, { scaleX: [0.3, 1.15, 1] }, { duration: τmid, delay: τmin });
}
