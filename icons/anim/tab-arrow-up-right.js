import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const tab = el.querySelector('path[d^="M4 20"]');
  animate(tab, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelectorAll('path[d^="m15 9"], path[d^="M9 9"]');
  animate(arrow, { translateX: [0, 1.5, -0.5, 0], translateY: [0, -1.5, 0.5, 0] }, { duration: τmid, delay: τmin });
}
