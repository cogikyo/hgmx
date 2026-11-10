import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d^="M3 "]');
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });

  const route = el.querySelector('path[d^="M21 19V7"]');
  animate(route, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
