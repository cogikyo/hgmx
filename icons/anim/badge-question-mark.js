import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const badge = el.querySelector('path[d^="M3.85"]');
  animate(badge, { scale: [1, 0.85, 1.1, 1] }, { duration: τmax });

  const mark = el.querySelectorAll('path[d^="M9.09"], line');
  animate(mark, { rotateY: [0, 360], scale: [1, 1.35, 1] }, { duration: τmid, delay: τmin });
}
