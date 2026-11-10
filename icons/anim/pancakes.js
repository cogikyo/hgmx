import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const stack = el.querySelectorAll('path:not([d^="M16 16"])');
  animate(stack, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const syrup = el.querySelector('path[d^="M16 16"]');
  animate(syrup, { scaleY: [0, 1.15, 1] }, { duration: τmid, delay: τmin });
}
