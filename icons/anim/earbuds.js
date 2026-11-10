import { animate, stagger } from "../motion.js";

export default function (el) {
  const buds = [el.querySelector('path[d^="M6 3"]'), el.querySelector('path[d^="M18 3"]')];
  animate(buds, { translateY: [0, -2, 1, 0] }, { duration: 0.3, delay: stagger(0.12) });
}
