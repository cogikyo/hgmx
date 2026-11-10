import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelectorAll('path[d^="m17 21"], path[d^="M10 18"]');
  animate(arrow, { translateX: [0, 4, -1, 0] }, { duration: τmax });

  const digit = el.querySelector('rect[x="15"]');
  animate(digit, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
