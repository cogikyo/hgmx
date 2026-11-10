import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const base = el.querySelector('path[d="M6 9v12"]');
  animate(base, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const x = el.querySelectorAll('path[d^="m15.5"]');
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });

  const head = el.querySelector('circle[cx="18"]');
  animate(head, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
