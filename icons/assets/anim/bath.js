import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const tub = el.querySelector('path[d^="M9 5"]');
  animate(tub, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const water = el.querySelector('path[d^="M2 12"]');
  animate(water, { translateY: [4, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
