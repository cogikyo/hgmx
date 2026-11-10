import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const tail = el.querySelector('path[d^="M22 16"]');
  animate(tail, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const wraps = [el.querySelector('path[d^="m8 18"]'), el.querySelector('path[d^="m6 11"]')];
  animate(wraps, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.12, { startDelay: τmin }) });
}
