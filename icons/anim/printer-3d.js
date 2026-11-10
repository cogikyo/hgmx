import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const head = el.querySelector("g");
  animate(head, { translateX: [0, -1.5, 1.5, 0] }, { duration: τmax });

  const print = el.querySelector('path[d^="M15 20"]');
  animate(print, { scaleY: [0, 1] }, { duration: τmid, delay: τmin });
}
