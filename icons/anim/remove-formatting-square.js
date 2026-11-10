import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const subject = el.querySelectorAll('rect, path[d^="M7 9"], path[d="M13 7 8 17"], path[d="M7 17h3"]');
  animate(subject, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const x = el.querySelectorAll('path[d="m17 14-3 3"], path[d="m14 14 3 3"]');
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
