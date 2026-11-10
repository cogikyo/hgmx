import { animate, τmax } from "../motion.js";

export default function (el) {
  const flap = el.querySelector('path[d^="M12 12V4"]');
  animate(flap, { scaleX: [1, 0.3, 1.1, 1] }, { duration: τmax });
}
