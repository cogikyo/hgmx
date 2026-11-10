import { animate, stagger } from "../motion.js";

export default function (el) {
  const air = el.querySelectorAll('path[d^="M6.6"], path[d^="M18 17.5"]');
  animate(air, { scaleY: [1, 1.3, 0.9, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
