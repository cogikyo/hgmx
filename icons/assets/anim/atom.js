import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const one = el.querySelector('path[d^="M20.2"]');
  animate(one, { rotate: [0, 180] }, { duration: τmax });

  const other = el.querySelector('path[d^="M15.7"]');
  animate(other, { rotate: [0, -180] }, { duration: τmax });

  const nucleus = el.querySelector("circle");
  animate(nucleus, { scale: [0, 1.75, 1] }, { duration: τmid, delay: τmin });
}
