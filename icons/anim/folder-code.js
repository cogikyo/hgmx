import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M20 20"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const open = el.querySelector('path[d^="M10 10.5"]');
  animate(open, { translateX: [2, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const close = el.querySelector('path[d^="m14 10.5"]');
  animate(close, { translateX: [-2, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
