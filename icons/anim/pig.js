import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const pig = el.querySelector("g");
  animate(pig, { translateY: [0, -2, 0.5, 0] }, { duration: τmax });

  const tail = el.querySelector('path[d^="M2.3 7"]');
  animate(tail, { rotate: [0, -30, 20, 0] }, { duration: τmid, delay: τmin });
}
