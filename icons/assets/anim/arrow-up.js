import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateY: [0, -5, 2, 0] }, { duration: τmax });

  const head = el.querySelector('path[d^="m5 12"]');
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
