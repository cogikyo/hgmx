import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const envelope = el.querySelector('path[d^="M21.2"]');
  animate(envelope, { scaleY: [0.85, 1.1, 1] }, { duration: τmax });

  const fold = el.querySelector('path[d^="m22 10"]');
  animate(fold, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
