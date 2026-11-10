import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const lemon = el.querySelector('path[d^="M17.6"]');
  animate(lemon, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const glint = el.querySelector('path[d^="M7 12"]');
  animate(glint, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
