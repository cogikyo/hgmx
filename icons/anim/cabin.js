import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelector('path[d^="M14 14h6"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const logs = el.querySelectorAll('path[d="M4 18h6"], path[d="M14 18h6"]');
  animate(logs, { scaleX: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
