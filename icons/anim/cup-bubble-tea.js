import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cup = el.querySelectorAll('path[d^="m6 10"], path[d^="M5 10"]');
  animate(cup, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pearls = el.querySelectorAll('path[d$="h.01"]');
  animate(pearls, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
