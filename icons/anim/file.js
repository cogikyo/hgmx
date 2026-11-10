import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M6 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const fold = el.querySelector('path[d^="M14 2"]');
  animate(fold, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
