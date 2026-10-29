import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const slice = el.querySelector('path[d^="M21.964"]');
  animate(slice, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const toppings = el.querySelectorAll('path[d^="m12 14"], path[d^="m13.75"]');
  animate(toppings, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
