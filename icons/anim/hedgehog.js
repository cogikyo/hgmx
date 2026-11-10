import { animate, τmax } from "../motion.js";

export default function (el) {
  const spines = el.querySelector('path[d^="M4 11"]');
  animate(spines, { scale: [1, 1.12, 0.97, 1] }, { duration: τmax });
}
