import { animate, τmax } from "../motion.js";

export default function (el) {
  const pestle = el.querySelector('path[d^="M22 4a2"]');
  animate(pestle, { rotate: [0, -10, 8, -4, 0] }, { duration: τmax });
}
