import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const base = el.querySelector('path[d="M6 9v12"]');
  animate(base, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const request = el.querySelector('path[d^="M13 6"]');
  animate(request, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const plus = [el.querySelector('path[d="M18 15v6"]'), el.querySelector('path[d="M21 18h-6"]')];
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmin, delay: τmid });
}
