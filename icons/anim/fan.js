import { animate, τmax } from "../motion.js";

export default function (el) {
  const blades = el.querySelector('path[d^="M10.827"]');
  animate(blades, { rotate: [0, 360], scale: [1, 0.85, 1] }, { duration: τmax });
}
