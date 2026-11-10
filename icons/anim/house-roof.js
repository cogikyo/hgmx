import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelector('path[d^="M20 14"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const roof = el.querySelector('path[d^="M2.6 10.4"]');
  animate(roof, { scaleY: [1, 0.85, 1] }, { duration: τmid, delay: τmin });
}
