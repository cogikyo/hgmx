import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bottle = el.querySelector('path[d^="M6 7"]');
  animate(bottle, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cap = el.querySelector("rect");
  animate(cap, { translateY: [0, -2, 0.5, 0], rotate: [0, -8, 3, 0] }, { duration: τmid, delay: τmin });
}
