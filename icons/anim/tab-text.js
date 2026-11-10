import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const tab = el.querySelector('path[d^="M4 20"]');
  animate(tab, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lines = el.querySelectorAll('path[d^="M8 "]');
  animate(lines, { opacity: [0, 1], x: [-3, 0] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
