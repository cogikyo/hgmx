import { animate, stagger } from "../motion.js";

export default function (el) {
  const pickets = el.querySelectorAll('path[d^="M4 3"], path[d^="m12 3"], path[d^="m20 3"]');
  animate(pickets, { translateY: [0, -2, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
