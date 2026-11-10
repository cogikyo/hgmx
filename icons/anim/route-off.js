import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const path = el.querySelectorAll('path:not([d^="m2 2 "])');
  animate(path, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m2 2 "]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
