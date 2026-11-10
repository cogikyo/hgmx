import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glass = el.querySelector('path[d^="M5.116"]');
  animate(glass, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const water = el.querySelector('path[d^="M6 12"]');
  animate(water, { rotate: [0, -10, 6, 0] }, { duration: τmid, delay: τmin });
}
