import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const board = el.querySelector('path[d^="M16 4"]');
  animate(board, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lines = el.querySelectorAll('path[d$="h4"]');
  animate(lines, { opacity: [0, 1], x: [-3, 0] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });

  const bullets = el.querySelectorAll('path[d$="h.01"]');
  animate(bullets, { scale: [0, 1.75, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
