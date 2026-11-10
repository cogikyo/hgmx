import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const tick = el.querySelector('path[d^="M16 6"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const x = el.querySelectorAll('path[d^="m14 14"], path[d^="m20 14"]');
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
