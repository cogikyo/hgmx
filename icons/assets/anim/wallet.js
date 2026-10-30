import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M3 5"]');
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const fold = el.querySelector('path[d^="M19 7"]');
  animate(fold, { translateY: [-3, 1, 0] }, { duration: τmid, delay: τmin });
}
