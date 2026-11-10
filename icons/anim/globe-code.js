import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const globe = el.querySelectorAll('path[d^="M15.861"], path[d="M2 12h8.5"]');
  animate(globe, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const open = el.querySelector('path[d^="M15.5 10 "]');
  animate(open, { translateX: [2, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const close = el.querySelector('path[d^="M19.5 5 "]');
  animate(close, { translateX: [-2, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
