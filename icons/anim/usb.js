import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const stem = el.querySelector('path[d^="M4.7"]');
  animate(stem, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const head = el.querySelector('path[d^="m21 3"]');
  animate(head, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
