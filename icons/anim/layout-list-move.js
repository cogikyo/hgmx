import { animate, stagger } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d^="M14"]');
  animate(lines, { translateY: [0, 3, -1, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
