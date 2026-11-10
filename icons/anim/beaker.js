import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const beaker = el.querySelector("g");
  animate(beaker, { rotate: [0, -10, 8, -4, 0] }, { duration: τmax });

  const liquid = el.querySelector('path[d^="M6 14"]');
  animate(liquid, { rotate: [0, 12, -8, 0] }, { duration: τmid, delay: τmin });
}
