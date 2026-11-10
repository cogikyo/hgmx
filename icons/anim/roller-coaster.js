import { animate, τmax } from "../motion.js";

export default function (el) {
  const track = el.querySelector('path[d^="M2 19V9"]');
  animate(track, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
}
