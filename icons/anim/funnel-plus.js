import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const funnel = el.querySelector('path[d^="M13.354"]');
  animate(funnel, { scaleX: [1, 0.8, 1.1, 1], scaleY: [1, 1.1, 0.95, 1] }, { duration: τmax });

  const plus = el.querySelectorAll('path[d="M16 6h6"], path[d="M19 3v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
