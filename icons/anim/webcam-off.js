import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelectorAll('path[d^="M16.566"], path[d^="M8.478"]');
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
