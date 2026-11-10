import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const shrimp = el.querySelector("g");
  animate(shrimp, { rotate: [0, -8, 6, -3, 0] }, { duration: τmax });

  const legs = el.querySelector('path[d^="M8 16c"]');
  animate(legs, { rotate: [0, -30, 10, 0] }, { duration: τmid, delay: τmin });
}
