import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const base = el.querySelector('path[d="M5 9v12"]');
  animate(base, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const request = el.querySelector('path[d="M12 6h5a2 2 0 0 1 2 2v3"]');
  animate(request, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const arrow = el.querySelector('path[d^="m15 9"]');
  animate(arrow, { translateX: [3, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const plus = [el.querySelector('path[d="M19 15v6"]'), el.querySelector('path[d="M22 18h-6"]')];
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmin, delay: τmid });
}
