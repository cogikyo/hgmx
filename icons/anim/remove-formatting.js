import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const type = el.querySelectorAll('path[d^="M4 7"], path[d="M5 20h6"], path[d="M13 4 8 20"]');
  animate(type, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const x = el.querySelectorAll('path[d="m15 15 5 5"], path[d="m20 15-5 5"]');
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
