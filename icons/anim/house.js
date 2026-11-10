import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelector('path[d^="M3 10"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const door = el.querySelector('path[d^="M15 21"]');
  animate(door, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
