import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const waves = el.querySelectorAll('path:not([d^="m2 2"])');
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
