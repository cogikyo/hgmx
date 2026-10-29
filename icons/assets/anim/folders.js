import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector('path[d^="M20 5"]');
  animate(front, { scaleY: [1, 0.85, 1.05, 1], translateY: [0, 1, 0] }, { duration: τmax });

  const back = el.querySelector('path[d^="M3 8.268"]');
  animate(back, { translateX: [4, -1, 0], translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
