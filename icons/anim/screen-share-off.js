import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector('path[d^="M13 3H4"]');
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const x = el.querySelectorAll('path[d^="m22 3"], path[d^="m17 3"]');
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
