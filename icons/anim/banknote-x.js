import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const note = el.querySelector('path[d^="M13 18"]');
  animate(note, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const x = el.querySelectorAll('path[d="m17 17 5 5"], path[d="m22 17-5 5"]');
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
