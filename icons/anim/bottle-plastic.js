import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bottle = el.querySelector('path[d^="M17 14.5"]');
  animate(bottle, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cap = el.querySelector('path[d^="M10 6.1"]');
  animate(cap, { translateY: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
}
