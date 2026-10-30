import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const waves = el.querySelectorAll('path[d^="M16 9"], path[d^="M19.364"]');
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });

  const slash = el.querySelector('path[d="m2 2 20 20"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
