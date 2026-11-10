import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const shell = el.querySelector('path[d^="M12 20c"]');
  animate(shell, { scaleX: [1, 1.15, 0.97, 1] }, { duration: τmax });
  const left = el.querySelector('path[d^="M6.08"]');
  animate(left, { rotate: [0, 14, -4, 0] }, { duration: τmax });
  const right = el.querySelector('path[d^="M19.1"]');
  animate(right, { rotate: [0, -14, 4, 0] }, { duration: τmax });
  const front = el.querySelector('path[d^="M5 4.8"]');
  animate(front, { scaleY: [1, 1.3, 1] }, { duration: τmid, delay: τmin });
}
