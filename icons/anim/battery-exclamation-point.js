import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelectorAll('path[d^="M14 6"], path[d^="M6 18"]');
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const mark = el.querySelector("g");
  animate(mark, { rotate: [0, -18, 14, -8, 0] }, { duration: τmid, delay: τmin });

  const dot = el.querySelector('path[d="M10 17h.01"]');
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
