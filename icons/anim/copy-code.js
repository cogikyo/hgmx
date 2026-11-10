import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("g");
  animate(front, { translateX: [-6, 1, 0], translateY: [-6, 1, 0] }, { duration: τmax });

  const back = el.querySelector('path[d^="M4 16"]');
  animate(back, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const open = el.querySelector('path[d^="m13 13"]');
  animate(open, { translateX: [2, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const close = el.querySelector('path[d^="m17 13"]');
  animate(close, { translateX: [-2, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
