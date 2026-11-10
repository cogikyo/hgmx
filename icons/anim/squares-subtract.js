import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const kept = el.querySelector('path[d^="M16 4"]');
  animate(kept, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const removed = el.querySelectorAll('path:not([d^="M16 4"])');
  animate(removed, { translateX: [0, 2, -0.5, 0], translateY: [0, 2, -0.5, 0] }, { duration: τmid, delay: τmin });
}
