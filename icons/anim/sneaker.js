import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const shoe = el.querySelector("g");
  animate(shoe, { rotate: [0, -10, 4, 0] }, { duration: τmax });

  const laces = el.querySelectorAll('path[d^="M14.1 "], path[d^="M17.4 "]');
  animate(laces, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
