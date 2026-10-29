import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const tower = el.querySelector('path[d^="M6 21"]');
  animate(tower, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const wings = el.querySelector('path[d^="M6 10"]');
  animate(wings, { scaleY: [0, 1.15, 1] }, { duration: τmid, delay: τmin });

  const door = el.querySelector('path[d^="M14 21"]');
  animate(door, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
