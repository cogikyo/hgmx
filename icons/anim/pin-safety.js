import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const arm = el.querySelector('path[d^="M2.9 16.9"]');
  animate(arm, { rotate: [0, -10, 0] }, { duration: τmid });

  const clasp = el.querySelector('path[d^="M20.8"]');
  animate(clasp, { scale: [1, 1.1, 1] }, { duration: τmin, delay: τmid });
}
