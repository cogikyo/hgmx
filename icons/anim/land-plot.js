import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const plot = el.querySelectorAll('path:not([d^="m12 8"])');
  animate(plot, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flag = el.querySelector('path[d^="m12 8"]');
  animate(flag, { rotate: [0, -10, 6, -3, 0] }, { duration: τmid, delay: τmin });
}
