import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M4 21"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const mark = el.querySelector("g");
  animate(mark, { rotate: [0, -18, 14, -8, 0] }, { duration: τmid, delay: τmin });

  const dot = el.querySelector('path[d$="H20"]');
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
