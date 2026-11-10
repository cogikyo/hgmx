import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelectorAll('path[d^="M21 15.3"], path[d^="M6.4 6.4"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d="m2 2 20 20"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
