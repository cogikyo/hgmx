import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glass = el.querySelector('path[d^="m18 6"]');
  animate(glass, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const straw = el.querySelector('path[d^="M9 6"]');
  animate(straw, { rotate: [0, -15, 10, 0] }, { duration: τmid, delay: τmin });
}
