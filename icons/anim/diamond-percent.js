import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M2.7"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d="m14.5 9.5-5 5"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const dots = el.querySelectorAll('path[d$="h.01"]');
  animate(dots, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09, { startDelay: τmin }) });
}
