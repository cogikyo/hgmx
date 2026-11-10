import { animate, stagger } from "../motion.js";

export default function (el) {
  const glints = el.querySelectorAll('path[d^="M11 6"], path[d^="m16 7"]');
  animate(glints, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
