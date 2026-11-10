import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const drop = el.querySelectorAll('path[d^="M18.715"], path[d^="M8.795"]');
  animate(drop, { translateY: [-2.5, 0.5, 0] }, { duration: τmid });
  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
