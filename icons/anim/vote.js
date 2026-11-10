import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = el.querySelector('path[d^="M5 7"]');
  animate(box, { scaleY: [1, 0.9, 1.05, 1] }, { duration: τmax });
  const tick = el.querySelector('path[d^="m9 12"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
