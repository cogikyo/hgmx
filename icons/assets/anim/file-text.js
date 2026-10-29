import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M6 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const lines = el.querySelectorAll('path[d$="H8"]');
  animate(lines, { opacity: [0, 1], x: [-3, 0] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
