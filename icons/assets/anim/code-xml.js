import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const right = el.querySelector('path[d^="m18 16"]');
  animate(right, { translateX: [-5, 1, 0], opacity: [0, 1] }, { duration: τmax });

  const left = el.querySelector('path[d^="m6 8"]');
  animate(left, { translateX: [5, -1, 0], opacity: [0, 1] }, { duration: τmax });

  const slash = el.querySelector('path[d^="m14.5"]');
  animate(slash, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
