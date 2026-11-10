import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("g");
  animate(front, { rotate: [0, -18, 14, -8, 0] }, { duration: τmax });

  const back = el.querySelector('path[d^="M2 7"]');
  animate(back, { rotate: [0, -12, 6, 0] }, { duration: τmid, delay: τmin });
}
