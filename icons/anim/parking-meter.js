import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const dial = el.querySelector('path[d^="M9 9"]');
  animate(dial, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
