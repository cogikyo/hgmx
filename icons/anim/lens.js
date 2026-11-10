import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const reflection = el.querySelector('path[d^="M6 12"]');
  animate(reflection, { rotate: [0, 45, -10, 0] }, { duration: τmax });

  const glint = el.querySelector('path[d^="M15.725"]');
  animate(glint, { scale: [1, 1.75, 1] }, { duration: τmin, delay: τmid });
}
