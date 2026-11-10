import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const holder = el.querySelector('path[d^="M18 18"]');
  animate(holder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const wick = el.querySelector('path[d="M10 9v2"]');
  animate(wick, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
