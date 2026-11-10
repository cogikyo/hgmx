import { animate, τmax } from "../motion.js";

export default function (el) {
  const towel = el.querySelector('path[d^="M6.5 3h11"]');
  animate(towel, { rotate: [0, -5, 3, 0] }, { duration: τmax });
}
