import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const note = el.querySelector("g");
  animate(note, { rotate: [0, -12, 4, 0], translateY: [0, -4, 1, 0] }, { duration: τmax });

  const fold = el.querySelector('path[d^="M15 3v5"]');
  animate(fold, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });

  const plus = el.querySelectorAll('path[d="M18 15v6"], path[d="M21 18h-6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
