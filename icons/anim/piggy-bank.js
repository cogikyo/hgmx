import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const pig = el.querySelector("g");
  animate(pig, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });

  const tail = el.querySelector('path[d^="M2 8"]');
  animate(tail, { rotate: [0, -30, 15, 0] }, { duration: τmid, delay: τmin });
}
