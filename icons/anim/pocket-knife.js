import { animate, τmax } from "../motion.js";

export default function (el) {
  const blade = el.querySelector('path[d^="M18 11.66"]');
  animate(blade, { rotate: [0, -12, 4, 0] }, { duration: τmax });
}
