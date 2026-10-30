import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { rotate: [0, -360], scale: [1, 0.85, 1] }, { duration: τmax });

  const head = el.querySelector('path[d^="M3 3v5"]');
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
