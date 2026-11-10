import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const base = el.querySelector('path[d^="M3 20"]');
  animate(base, { scaleY: [1, 0.7, 1] }, { duration: τmin });

  const head = el.querySelector("g");
  animate(head, { rotate: [0, 12, -4, 0] }, { duration: τmid, delay: τmin });
}
