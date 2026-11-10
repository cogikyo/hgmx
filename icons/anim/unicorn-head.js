import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const head = el.querySelector("g");
  animate(head, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });

  const horn = el.querySelectorAll('path[d^="M15.5 10"], path[d^="m15.6"]');
  animate(horn, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
