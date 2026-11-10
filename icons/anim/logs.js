import { animate, stagger } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d$="h8"]');
  animate(lines, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
