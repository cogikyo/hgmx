import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = el.querySelector('path[d^="M20 5.5"]');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flag = el.querySelector('path[d^="M14 9"]');
  animate(flag, { rotate: [0, -15, 0] }, { duration: τmid, delay: τmin });
}
