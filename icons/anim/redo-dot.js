import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { rotate: [0, 30, -10, 0] }, { duration: τmax });

  const head = el.querySelector('path[d^="M21 7"]');
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });

  const dot = el.querySelector("circle");
  animate(dot, { scale: [1, 1.75, 1] }, { duration: τmid, delay: τmin });
}
