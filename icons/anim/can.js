import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const can = el.querySelectorAll('path[d^="M21 10.5a"], path[d^="M3 10.5"]');
  animate(can, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lid = el.querySelector('path[d^="M21 10.5A"]');
  animate(lid, { rotate: [0, -12, 4, 0] }, { duration: τmid, delay: τmin });
}
