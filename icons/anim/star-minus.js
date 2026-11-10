import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const star = el.querySelector('path[d^="M17.688"]');
  animate(star, { scale: [0.5, 1.15, 1], rotate: [0, 360] }, { duration: τmid });
  animate(star, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const minus = el.querySelector('path[d="M15 18h6"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
