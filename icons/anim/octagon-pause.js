import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const sign = el.querySelector('path[d^="M2.586"]');
  animate(sign, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = el.querySelectorAll('path[d="M10 15V9"], path[d="M14 15V9"]');
  animate(bars, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
