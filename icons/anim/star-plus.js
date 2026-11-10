import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const star = el.querySelector('path[d^="M11.013"]');
  animate(star, { scale: [0.5, 1.15, 1], rotate: [0, 360] }, { duration: τmid });
  animate(star, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M15 18h6"], path[d="M18 15v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
