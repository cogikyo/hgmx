import { animate, τmax } from "../motion.js";

export default function (el) {
  const petals = el.querySelector('path[d^="M12 16.5A"]');
  animate(petals, { scale: [0, 1.15, 1] }, { duration: τmax });
}
