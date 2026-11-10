import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bottle = el.querySelector('path[d^="M9 2"]');
  animate(bottle, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cap = el.querySelector('path[d^="M8 2"]');
  animate(cap, { translateY: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
}
