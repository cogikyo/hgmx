import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const holder = el.querySelectorAll('path[d^="M18 18"], path[d="m13 13-1-2"]');
  animate(holder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flame = el.querySelector('path[d^="M10 2S"]');
  animate(flame, { scaleY: [1, 0.8, 1.25, 1] }, { duration: τmid, delay: τmin });
}
