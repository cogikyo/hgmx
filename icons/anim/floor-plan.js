import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M9 21"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const walls = el.querySelectorAll('path[d="M9 3v7"], path[d="M21 10h-7"], path[d="M3 15h9"]');
  animate(walls, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
