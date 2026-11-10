import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = el.querySelector('path[d^="M20 8c"]');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const kernels = el.querySelector('path[d^="M18 8a2"]');
  animate(kernels, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
