import { animate, τmax } from "../motion.js";

export default function (el) {
  const tail = el.querySelector('path[d^="M7 10.67"]');
  animate(tail, { rotate: [0, -28, 18, -6, 0] }, { duration: τmax });
}
