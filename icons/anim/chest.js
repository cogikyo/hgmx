import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const chest = el.querySelectorAll('path[d^="M8 19"], path[d^="M2 11"]');
  animate(chest, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const latch = el.querySelector('path[d^="M16 11"]');
  animate(latch, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
