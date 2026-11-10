import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const tower = el.querySelector('path[d^="M6 21"]');
  animate(tower, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const wings = el.querySelector("g");
  animate(wings, { scaleY: [0, 1.15, 1] }, { duration: τmid, delay: τmin });

  const plus = el.querySelectorAll('path[d="M16 19h6"], path[d="M19 16v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
