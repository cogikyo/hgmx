import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const stems = el.querySelectorAll('path[d$="V4"]');
  animate(stems, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(τmin / 2) });
  const bar = el.querySelector('path[d="M6 12H18"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
