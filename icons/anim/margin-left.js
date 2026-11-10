import { animate, τmax } from "../motion.js";

export default function (el) {
  const edge = el.querySelector('path[d="M3 3v18"]');
  animate(edge, { translateX: [0, -2, 0.5, 0] }, { duration: τmax });
}
