import { animate, stagger } from "../motion.js";

export default function (el) {
  const glints = el.querySelectorAll('path[d^="M10 6.6"], path[d^="M15 7.5"]');
  animate(glints, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
