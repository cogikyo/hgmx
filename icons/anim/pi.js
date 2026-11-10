import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bar = el.querySelector('path[d^="M4 7"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const legs = el.querySelector("g");
  animate(legs, { scaleY: [0, 1] }, { duration: τmin, delay: τmid });
}
