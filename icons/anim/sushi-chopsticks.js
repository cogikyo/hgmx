import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const roll = el.querySelectorAll('path[d^="M6 11"], path[d^="M18 13"], path[d^="M11 14"]');
  animate(roll, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const left = el.querySelector('path[d="M2 18V2"]');
  animate(left, { translateX: [0, 2, -0.5, 0] }, { duration: τmid, delay: τmin });

  const right = el.querySelector('path[d="M22 18V2"]');
  animate(right, { translateX: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
}
