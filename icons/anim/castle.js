import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelectorAll('path[d^="M22 9v10"], path[d="M22 11H2"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const merlons = el.querySelectorAll('path[d="M10 5V3"], path[d="M14 5V3"]');
  animate(merlons, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
