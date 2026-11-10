import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const tube = el.querySelector("g");
  animate(tube, { rotate: [0, -10, 8, -4, 0] }, { duration: τmax });

  const liquid = el.querySelector('path[d^="M14.5 16"]');
  animate(liquid, { rotate: [0, 12, -8, 0] }, { duration: τmid, delay: τmin });
}
