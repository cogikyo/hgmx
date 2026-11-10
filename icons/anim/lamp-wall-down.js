import { animate, τmax } from "../motion.js";

export default function (el) {
  const shade = el.querySelector('path[d^="M19.929"]');
  animate(shade, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
