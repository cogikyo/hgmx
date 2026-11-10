import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const line = [el.querySelector('path[d="M6 3v12"]'), el.querySelector('path[d^="M15 6"]')];
  animate(line, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const tip = el.querySelector('path[d^="M18 9"]');
  animate(tip, { translateX: [-8, 1, 0], translateY: [8, -1, 0], opacity: [0, 1] }, { duration: τmax });

  const plus = [el.querySelector('path[d="M18 15v6"]'), el.querySelector('path[d="M21 18h-6"]')];
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
