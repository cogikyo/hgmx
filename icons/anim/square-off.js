import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelectorAll('path:not([d^="M22 22"])');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d^="M22 22"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
