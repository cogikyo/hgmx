import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M10.638"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const heart = el.querySelector('path[d^="M14.62"]');
  animate(heart, { scale: [1, 1.25, 0.9, 1.15, 1] }, { duration: τmid, delay: τmin });
}
