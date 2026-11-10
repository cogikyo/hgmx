import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [0, 2, -1, 0], translateY: [0, -2, 1, 0] }, { duration: τmax });

  const head = el.querySelector('path[d^="M15 3"]');
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
