import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cake = el.querySelectorAll('path[d^="M20 21"], path[d^="M4 16"]');
  animate(cake, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flames = el.querySelectorAll('path[d$="h.01"]');
  animate(flames, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
