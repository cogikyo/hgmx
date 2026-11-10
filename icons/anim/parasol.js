import { animate, τmax } from "../motion.js";

export default function (el) {
  const canopy = el.querySelector('path[d^="M20.425"]');
  animate(canopy, { scaleX: [0.3, 1.15, 1] }, { duration: τmax });
}
